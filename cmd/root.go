package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("directory", rootCmd.Flags().Lookup("directory"))
	viper.BindPFlag("max_requests", rootCmd.Flags().Lookup("max-requests"))
	viper.BindPFlag("ban_time", rootCmd.Flags().Lookup("ban-time"))

	viper.SetEnvPrefix("KFRS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return rootCmd
}

func printFlags() {
	log.Println("===================================================")
	log.Println("                   KFRS Settings                   ")
	log.Println("===================================================")
	log.Printf(" ● Host         → %s\n", Host)
	log.Printf(" ● Port         → %d\n", Port)
	log.Printf(" ● Directory    → %s\n", Directory)
	log.Printf(" ● Max Requests → %d\n", MaxRequests)
	log.Printf(" ● Ban Time     → %d\n", BanTime)
	log.Println("====================================================")
}
