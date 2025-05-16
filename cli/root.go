package cli

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "libvirt-aws",
	Short: "AWS API emulation powered by libvirt",
	Long:  "AWS API emulation powered by libvirt",
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to execute command")
	}
}
