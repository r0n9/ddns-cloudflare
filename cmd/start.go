package cmd

import "github.com/spf13/cobra"

// StartCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Silent start ddns_cloudflare",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {
	initDaemon()

}
