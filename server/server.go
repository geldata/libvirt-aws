package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

const (
	DefaultPort = 5100
)

type AWSEmulatorServerOpts struct {
	BindTo string
	Port   int
	Debug  bool
	Region string // AWS region to pretend to be in
}

type AWSEmulatorServer struct {
	BindTo string
	Port   int
	Addr   string
	Debug  bool
	Region string
}

func NewAWSEmulatorServerOpts(fs *pflag.FlagSet) *AWSEmulatorServerOpts {
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
	return opts
}

func NewAWSEmulatorServer(opts *AWSEmulatorServerOpts) *AWSEmulatorServer {
	opts = mergeAWSEmulatorServerOpts(opts)
	return &AWSEmulatorServer{
		BindTo: opts.BindTo,
		Port:   opts.Port,
		Addr:   fmt.Sprintf("%s:%d", opts.BindTo, opts.Port),
		Debug:  opts.Debug,
		Region: opts.Region,
	}
}

func (s *AWSEmulatorServer) configureLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if s.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = zerolog.New(os.Stderr).With().Str("svc", "awsapi").Logger()
}

func (s *AWSEmulatorServer) Start() error {
	s.configureLogger()
	s.registerRoutes()
	log.Info().Msgf("starting server at %s", s.Addr)
	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start server at %s", s.Addr)
	}
	return nil
}
