package cli

import (
	"github.com/geldata/libvirt-aws/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the libvirt-aws emulator server",
	Long:  "Starts the libvirt-aws emulator server",
	Run: func(cmd *cobra.Command, args []string) {
		opts := server.NewAWSEmulatorServerOpts(cmd.Flags())
		log.Info().Msg("Starting libvirt-aws server...")
		server.NewAWSEmulatorServer(opts).Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("bind-to", "b", "", "Address to listen on")
	startCmd.Flags().IntP("port", "p", 5100, "TCP port to listen on")
	startCmd.Flags().BoolP("debug", "", false, "Enable debug logging")
	startCmd.Flags().StringP("database", "d", "pool.db", "Path to db file")
	startCmd.Flags().StringP("libvirt-uri", "l", "qemu:///system", "Libvirtd URI")
	startCmd.Flags().StringP("libvirt-image-pool", "i", "default", "Name or UUID of libvirt image pool to use for EBS emulation.")
	startCmd.Flags().StringP("libvirt-network", "n", "default", "Name or UUID of libvirt network to use for EIP emulation.")
	startCmd.Flags().StringP("region", "r", "us-east-2", "AWS region to pretend to be in")
}
