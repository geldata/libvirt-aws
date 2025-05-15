package cli

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the libvirt-aws server",
	Long:  "Starts the libvirt-aws server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("bind-to", "b", "", "Address to listen on")
	startCmd.Flags().IntP("port", "p", 5100, "TCP port to listen on")
	startCmd.Flags().StringP("database", "d", "pool.db", "Path to db file")
	startCmd.Flags().StringP("libvirt-uri", "l", "qemu:///system", "Libvirtd URI")
	startCmd.Flags().StringP("libvirt-image-pool", "i", "default", "Name or UUID of libvirt image pool to use for EBS emulation.")
	startCmd.Flags().StringP("libvirt-network", "n", "default", "Name or UUID of libvirt network to use for EIP emulation.")
	startCmd.Flags().StringP("region", "r", "us-east-2", "AWS region to pretend to be in")
}

func start() {
	log.Info().Msg("Starting libvirt-aws server...")
}
