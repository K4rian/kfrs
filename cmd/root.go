package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	Host        string
	Port        int
	Directory   string
	MaxRequests int
	BanTime     int
)

func BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "./kfrs",
		Short: "KF HTTP Redirect Server",
		Long:  "A command-line tool to run an HTTP Redirect Server used by the Killing Floor Dedicated Server (kfds).",
		Run: func(cmd *cobra.Command, args []string) {
			printFlags()
		},
	}

	helpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpFunc(cmd, args)
		os.Exit(0)
	})

	rootCmd.Flags().StringVar(&Host, "host", "0.0.0.0", "IP/Host to bind to")
	rootCmd.Flags().IntVar(&Port, "port", 9090, "TCP port to listen on")
	rootCmd.Flags().StringVar(&Directory, "directory", "./redirect", "directory to serve")
	rootCmd.Flags().IntVar(&MaxRequests, "max-requests", 20, "max requests per IP/minute")
	rootCmd.Flags().IntVar(&BanTime, "ban-time", 15, "ban duration (in minutes)")

	return rootCmd
}

func printFlags() {
	log.Println("===================================================")
	log.Println("                   KFRS Settings                   ")
	log.Println("===================================================")
	log.Printf(" ● Host        → %s\n", Host)
	log.Printf(" ● Port        → %d\n", Port)
	log.Printf(" ● Directory   → %s\n", Directory)
	log.Printf(" ● MaxRequests → %d\n", MaxRequests)
	log.Printf(" ● BanTime     → %d\n", BanTime)
	log.Println("====================================================")
}
