package cmd

import (
	"log"
	"net"
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
		Run:   runRootCommand,
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

func runRootCommand(cmd *cobra.Command, args []string) {
	// Init configuration values
	initConfig(cmd)

	// Validate config values
	validateConfig()

	// Print the final configuration
	printConfig()
}

func initConfig(cmd *cobra.Command) {
	Host = getConfigValue(cmd, "host", defaultHost).(string)
	Port = getConfigValue(cmd, "port", defaultPort).(int)
	Directory = getConfigValue(cmd, "directory", defaultDirectory).(string)
	MaxRequests = getConfigValue(cmd, "max-requests", defaultMaxRequests).(int)
	BanTime = getConfigValue(cmd, "ban-time", defaultBanTime).(int)
}

func getConfigValue(cmd *cobra.Command, key string, defaultValue any) any {
	switch defaultValue.(type) {
	case int:
	case string:
	default:
		log.Fatalf("Unsupported type for flag %s: %T", key, defaultValue)
	}

	// CLI arguments take precedence over environment variables
	if cmd.Flags().Changed(key) {
		val := cmd.Flag(key).Value.String()
		switch defaultValue.(type) {
		case int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				log.Fatalf("Invalid integer value for flag %s: %s", key, val)
			}
			return intVal
		case string:
			return val
		}
	}

	if viper.IsSet(key) {
		switch defaultValue.(type) {
		case int:
			return viper.GetInt(key)
		case string:
			return viper.GetString(key)
		}
	}
	return defaultValue
}

func validateConfig() {
	// Host validation (v4, v6, or domain)
	if net.ParseIP(Host) == nil {
		_, err := net.LookupHost(Host)
		if err != nil {
			log.Fatalf("Invalid host: %s", Host)
		}
	}

	// Port validation
	if Port <= 1024 || Port >= 65535 {
		log.Fatalf("Port must be between 1025 and 65534: %d", Port)
	}

	// The directory must exist
	if _, err := os.Stat(Directory); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", Directory)
	}

	// Max Requests validation (> 0 and < 100)
	if MaxRequests <= 0 || MaxRequests >= 100 {
		log.Fatalf("Max requests should be between 1 and 99: %d", MaxRequests)
	}

	// Ban Time validation (should be > 0)
	if BanTime <= 0 {
		log.Fatalf("Ban time should be greater than 0: %d", BanTime)
	}
}

func printConfig() {
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
