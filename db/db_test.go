package db

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewDBOpts(t *testing.T) {
	tests := []struct {
		name string
		fs   *pflag.FlagSet
		want *DBOpts
	}{
		{
			name: "default options",
			fs:   pflag.NewFlagSet("test", pflag.ContinueOnError),
			want: &DBOpts{
				DBFile:           "pool.db",
				EnableGormLogger: false,
				Config: &gorm.Config{
					Logger: logger.Discard,
				},
			},
		},
		{
			name: "custom options",
			fs: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("database", "custom.db", "custom database file")
				fs.Bool("debug", true, "enable gorm logger")
				return fs
			}(),
			want: &DBOpts{
				DBFile:           "custom.db",
				EnableGormLogger: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDBOpts(tt.fs)
			assert.Equal(t, tt.want.DBFile, got.DBFile)
			assert.Equal(t, tt.want.EnableGormLogger, got.EnableGormLogger)
		})
	}
}
func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		opts *DBOpts
		want *gorm.DB
	}{
		{
			name: "in-memory options",
			opts: &DBOpts{
				DBFile:           ":memory:",
				EnableGormLogger: false,
				Config: &gorm.Config{
					Logger: logger.Discard,
				},
			},
		},
		{
			name: "pflag options",
			opts: NewDBOpts(pflag.NewFlagSet("test", pflag.ContinueOnError)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(tt.opts)
			assert.NoError(t, err)
			assert.NotNil(t, db)
		})
	}
}
