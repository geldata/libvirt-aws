package server

import (
	"testing"

	"github.com/geldata/libvirt-aws/db"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func testDBOpts() *db.DBOpts {
	return &db.DBOpts{
		DBFile: ":memory:",
		Config: &gorm.Config{},
	}
}

func TestNewAWSEmulatorServer(t *testing.T) {
	tests := []struct {
		name string
		opts *AWSEmulatorServerOpts
		want *AWSEmulatorServer
	}{
		{
			name: "default options",
			opts: &AWSEmulatorServerOpts{
				BindTo: "",
				Port:   5100,
				Debug:  false,
				Region: "us-east-2",
				DBOpts: testDBOpts(),
			},
			want: &AWSEmulatorServer{
				Addr: ":5100",
				Opts: &AWSEmulatorServerOpts{
					BindTo: "",
					Port:   5100,
					Debug:  false,
					Region: "us-east-2",
				},
			},
		},
		{
			name: "custom options",
			opts: &AWSEmulatorServerOpts{
				BindTo: "169.254.169.254",
				Port:   9090,
				Debug:  true,
				Region: "us-west-2",
				DBOpts: testDBOpts(),
			},
			want: &AWSEmulatorServer{
				Addr: "169.254.169.254:9090",
				Opts: &AWSEmulatorServerOpts{
					BindTo: "169.254.169.254",
					Port:   9090,
					Debug:  true,
					Region: "us-west-2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAWSEmulatorServer(tt.opts)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Opts.BindTo, got.Opts.BindTo)
			assert.Equal(t, tt.want.Opts.Port, got.Opts.Port)
			assert.Equal(t, tt.want.Opts.Debug, got.Opts.Debug)
			assert.Equal(t, tt.want.Opts.Region, got.Opts.Region)
			assert.Equal(t, tt.want.Addr, got.Addr)
		})
	}
}

func TestNewAWSEmulatorServerFromFlags(t *testing.T) {
	tests := []struct {
		name string
		fs   *pflag.FlagSet
		want *AWSEmulatorServer
	}{
		{
			name: "custom options",
			fs: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("database", ":memory:", "custom database file")
				fs.Bool("debug", false, "enable gorm logger")
				fs.String("bind-to", "33.33.33.33", "bind to address")
				fs.Int("port", 8080, "port to listen on")
				fs.String("region", "us-west-2", "AWS region")
				return fs
			}(),
			want: &AWSEmulatorServer{
				Addr: "33.33.33.33:8080",
				Opts: &AWSEmulatorServerOpts{
					BindTo: "33.33.33.33",
					Port:   8080,
					Region: "us-west-2",
					DBOpts: &db.DBOpts{
						DBFile: ":memory:",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAWSEmulatorServerFromFlags(tt.fs)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Opts.BindTo, got.Opts.BindTo)
			assert.Equal(t, tt.want.Opts.Port, got.Opts.Port)
			assert.Equal(t, tt.want.Opts.Region, got.Opts.Region)
			assert.Equal(t, tt.want.Addr, got.Addr)
			assert.Equal(t, tt.want.Opts.DBOpts.DBFile, got.Opts.DBOpts.DBFile)
			assert.Equal(t, tt.want.Opts.DBOpts.EnableGormLogger, got.Opts.DBOpts.EnableGormLogger)
		})
	}
}
