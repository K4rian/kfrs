package cmd

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultHost        = "0.0.0.0"
	defaultPort        = 9090
	defaultDirectory   = "./redirect"
	defaultMaxRequests = 20
	defaultBanTime     = 15
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
		Use:   "kfrs",
		Short: "KF HTTP Redirect Server",
		Long:  "A command-line tool to run an HTTP Redirect Server used by the Killing Floor Dedicated Server (kfds).",
		Run: func(cmd *cobra.Command, args []string) {
			Host = getConfigValue(cmd, "host", defaultHost).(string)
			Port = getConfigValue(cmd, "port", defaultPort).(int)
			Directory = getConfigValue(cmd, "directory", defaultDirectory).(string)
			MaxRequests = getConfigValue(cmd, "max-requests", defaultMaxRequests).(int)
			BanTime = getConfigValue(cmd, "ban-time", defaultBanTime).(int)

			printFlags()
		},
	}

	helpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		helpFunc(cmd, args)
		os.Exit(0)
	})

	config := []struct {
		Key         string
		Default     interface{}
		Description string
	}{
		{"host", defaultHost, "IP/Host to bind to"},
		{"port", defaultPort, "TCP port to listen on"},
		{"directory", defaultDirectory, "Directory to serve"},
		{"max_requests", defaultMaxRequests, "Max requests per IP/minute"},
		{"ban_time", defaultBanTime, "Ban duration (in minutes)"},
	}
	for _, cfg := range config {
		switch v := cfg.Default.(type) {
		case string:
			rootCmd.Flags().String(cfg.Key, v, cfg.Description)
		case int:
			rootCmd.Flags().Int(cfg.Key, v, cfg.Description)
		}
		viper.BindPFlag(cfg.Key, rootCmd.Flags().Lookup(cfg.Key))
		viper.BindEnv(cfg.Key)
		viper.SetDefault(cfg.Key, cfg.Default)
	}

	viper.SetEnvPrefix("KFRS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return rootCmd
}

func getConfigValue(cmd *cobra.Command, key string, defaultValue any) any {
	// CLI arguments take precedence over environment variables
	if cmd.Flags().Changed(key) {
		val := cmd.Flag(key).Value.String()

		switch v := defaultValue.(type) {
		case int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				log.Fatalf("Invalid integer value for flag %s: %s", key, val)
			}
			return intVal
		case string:
			return val
		default:
			log.Fatalf("Unsupported type for flag %s: %T", key, v)
		}
	}

	if viper.IsSet(key) {
		switch v := defaultValue.(type) {
		case int:
			return viper.GetInt(key)
		case string:
			return viper.GetString(key)
		default:
			log.Fatalf("Unsupported type for flag %s: %T", key, v)
		}
	}
	return defaultValue
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
