package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/geldata/libvirt-aws/awsapi"
	"github.com/geldata/libvirt-aws/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

const (
	DefaultPort     = 5100
	TimeFieldFormat = "2006-01-02T15:04:05.000Z"
)

type AWSEmulatorServerOpts struct {
	BindTo string
	Port   int
	Debug  bool
	Region string // AWS region to pretend to be in
	DBOpts *db.DBOpts
}

type AWSEmulatorServer struct {
	Opts   *AWSEmulatorServerOpts
	Addr   string
	stop   chan os.Signal
	awsapi *awsapi.AWSAPI
}

func newAWSEmulatorServerOpts(fs *pflag.FlagSet) *AWSEmulatorServerOpts {
	var err error
	opts := &AWSEmulatorServerOpts{}
	opts.BindTo, err = fs.GetString("bind-to")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get bind-to flag")
	}
	opts.Port, err = fs.GetInt("port")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get port flag")
	}
	opts.Debug, err = fs.GetBool("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get debug flag")
	}
	opts.Region, err = fs.GetString("region")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get region flag")
	}
	opts.DBOpts = db.DefaultDBOpts()
	opts.DBOpts.DBFile, err = fs.GetString("database")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database flag")
	}
	opts.DBOpts.EnableGormLogger, err = fs.GetBool("debug")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get debug flag")
	}
	return opts
}

func DefaultAWSEmulatorServerOpts() *AWSEmulatorServerOpts {
	return &AWSEmulatorServerOpts{
		BindTo: "",
		Port:   DefaultPort,
		Debug:  false,
		Region: "us-east-2",
	}
}

func mergeAWSEmulatorServerOpts(opts *AWSEmulatorServerOpts) *AWSEmulatorServerOpts {
	if opts == nil {
		return DefaultAWSEmulatorServerOpts()
	}
	defaultOpts := DefaultAWSEmulatorServerOpts()
	if opts.Port == 0 {
		opts.Port = defaultOpts.Port
	}
	if opts.Region == "" {
		opts.Region = defaultOpts.Region
	}
	if opts.DBOpts == nil {
		opts.DBOpts = db.DefaultDBOpts()
	}
	return opts
}

func NewAWSEmulatorServer(opts *AWSEmulatorServerOpts) (*AWSEmulatorServer, error) {
	opts = mergeAWSEmulatorServerOpts(opts)
	db, err := db.NewDB(opts.DBOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	api := awsapi.NewAWSAPI(db, opts.Region)
	return &AWSEmulatorServer{
		Opts:   opts,
		Addr:   fmt.Sprintf("%s:%d", opts.BindTo, opts.Port),
		stop:   make(chan os.Signal, 1),
		awsapi: api,
	}, nil
}

func NewAWSEmulatorServerFromFlags(fs *pflag.FlagSet) (*AWSEmulatorServer, error) {
	opts := newAWSEmulatorServerOpts(fs)
	opts.DBOpts = db.NewDBOpts(fs)
	return NewAWSEmulatorServer(opts)
}

func (s *AWSEmulatorServer) configureLogger() {
	zerolog.TimeFieldFormat = TimeFieldFormat
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = zerolog.New(os.Stderr).With().Str("svc", "awsapi").Logger()
	if s.Opts.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug logging enabled")
	}
}

func (s *AWSEmulatorServer) Start() error {
	s.configureLogger()
	mux := http.NewServeMux()
	s.registerRoutes(mux)
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: mux,
	}
	signal.Notify(s.stop, os.Interrupt)
	go func() {
		log.Info().Msgf("starting server at %s", s.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("failed to start server at %s", s.Addr)
		}
	}()
	<-s.stop
	log.Info().Msg("shutting down server")
	srv.Shutdown(context.Background())
	return nil
}

func (s *AWSEmulatorServer) Stop() {
	log.Info().Msg("sending server interrupt signal")
	s.stop <- os.Interrupt
	close(s.stop)
}
