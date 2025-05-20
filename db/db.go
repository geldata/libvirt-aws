package db

import (
	"fmt"

	"github.com/geldata/libvirt-aws/models"
	"github.com/glebarez/sqlite"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBOpts struct {
	DBFile           string
	EnableGormLogger bool
	Config           *gorm.Config
}

func DefaultDBOpts() *DBOpts {
	return &DBOpts{
		DBFile:           "pool.db",
		EnableGormLogger: false,
		Config: &gorm.Config{
			Logger: logger.Discard,
		},
	}
}

func NewDBOpts(fs *pflag.FlagSet) *DBOpts {
	var err error
	opts := &DBOpts{}
	opts.DBFile, err = fs.GetString("database")
	if err != nil {
		opts.DBFile = "pool.db"
	}
	opts.EnableGormLogger, err = fs.GetBool("debug")
	if err != nil {
		opts.EnableGormLogger = false
	}
	opts.Config = &gorm.Config{
		Logger: logger.Discard,
	}
	return opts
}

func NewDB(opts *DBOpts) (*gorm.DB, error) {
	if opts.EnableGormLogger {
		opts.Config.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(sqlite.Open(opts.DBFile), opts.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	err = db.AutoMigrate(models.All()...)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}
