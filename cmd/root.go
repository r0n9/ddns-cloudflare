package cmd

import (
	"fmt"
	"github.com/r0n9/ddns-cloudflare/cmd/flags"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "ddns-cloudflare",
	Short: "A dynamic DNS records manage program by using Cloudflare API.",
	Long: `Cloudflare provides an API that allows you to manage DNS records programmatically. 
To set up a Cloudflare dynamic DNS, you'll need to run a process on a client inside your network 
that does two main actions: get your network's current public IP address and automatically update the corresponding DNS record.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&flags.Conf, "conf", "config.json", "config file")
}
