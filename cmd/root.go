package cmd

import (
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/K4rian/kfrs/internal/config"
	"github.com/K4rian/kfrs/internal/log"
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

	conf := []struct {
		Key         string
		Default     interface{}
		Description string
	}{
		{"host", config.DefaultHost, "IP/Host to bind to"},
		{"port", config.DefaultPort, "TCP port to listen on"},
		{"serve-dir", config.DefaultServeDir, "Directory to serve"},
		{"max-requests", config.DefaultMaxRequests, "Max requests per IP/minute"},
		{"ban-time", config.DefaultBanTime, "Ban duration (in minutes)"},
		{"log-to-file", config.DefaultLogToFile, "Enable log file output"},
		{"log-level", config.DefaultLogLevel, "Set the log level (debug, info, warn, error)"},
		{"log-file", config.DefaultLogFile, "Specify the log file path"},
		{"log-file-format", config.DefaultLogFileFormat, "Specify the log format (text or json)"},
		{"log-max-size", config.DefaultLogMaxSize, "Set the maximum log file size in MB"},
		{"log-max-backups", config.DefaultLogMaxBackups, "Set the maximum number of backup log files to keep"},
		{"log-max-age", config.DefaultLogMaxAge, "Set the maximum number of days to retain old log files"},
	}
	for _, cfg := range conf {
		switch v := cfg.Default.(type) {
		case bool:
			rootCmd.Flags().Bool(cfg.Key, v, cfg.Description)
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
	conf := config.Get()

	// Set configuration values
	setConfigValues(cmd, conf)

	// Validate configuration values
	if err := conf.Validate(); err != nil {
		stdlog.Fatalf("Configuration validation error: %v", err)
	}

	// Init the logger
	log.Init()

	// Print the final configuration
	printConfigValues(conf)
}

func getConfigValue(cmd *cobra.Command, key string, defaultValue any) any {
	switch defaultValue.(type) {
	case bool:
	case int:
	case string:
	default:
		stdlog.Fatalf("Unsupported type for flag %s: %T", key, defaultValue)
	}

	// CLI arguments take precedence over environment variables
	if cmd.Flags().Changed(key) {
		val := cmd.Flag(key).Value.String()
		switch defaultValue.(type) {
		case bool:
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				stdlog.Fatalf("Invalid bool value for flag %s: %s", key, val)
			}
			return boolVal
		case int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				stdlog.Fatalf("Invalid integer value for flag %s: %s", key, val)
			}
			return intVal
		case string:
			return val
		}
	}

	if viper.IsSet(key) {
		switch defaultValue.(type) {
		case bool:
			return viper.GetBool(key)
		case int:
			return viper.GetInt(key)
		case string:
			return viper.GetString(key)
		}
	}
	return defaultValue
}

func setConfigValues(cmd *cobra.Command, conf *config.Config) {
	conf.Host = getConfigValue(cmd, "host", config.DefaultHost).(string)
	conf.Port = getConfigValue(cmd, "port", config.DefaultPort).(int)
	conf.ServeDir = getConfigValue(cmd, "serve-dir", config.DefaultServeDir).(string)
	conf.MaxRequests = getConfigValue(cmd, "max-requests", config.DefaultMaxRequests).(int)
	conf.BanTime = getConfigValue(cmd, "ban-time", config.DefaultBanTime).(int)
	conf.LogToFile = getConfigValue(cmd, "log-to-file", config.DefaultLogToFile).(bool)
	conf.LogLevel = strings.ToLower(getConfigValue(cmd, "log-level", config.DefaultLogLevel).(string))
	conf.LogFile = getConfigValue(cmd, "log-file", config.DefaultLogFile).(string)
	conf.LogFileFormat = strings.ToLower(getConfigValue(cmd, "log-file-format", config.DefaultLogFileFormat).(string))
	conf.LogMaxSize = getConfigValue(cmd, "log-max-size", config.DefaultLogMaxSize).(int)
	conf.LogMaxBackups = getConfigValue(cmd, "log-max-backups", config.DefaultLogMaxBackups).(int)
	conf.LogMaxAge = getConfigValue(cmd, "log-max-age", config.DefaultLogMaxAge).(int)
}

func printConfigValues(conf *config.Config) {
	log.Logger.Info("===================================================")
	log.Logger.Info("                   KFRS Settings                   ")
	log.Logger.Info("===================================================")
	log.Logger.Info(fmt.Sprintf(" ● Host            → %s", conf.Host))
	log.Logger.Info(fmt.Sprintf(" ● Port            → %d", conf.Port))
	log.Logger.Info(fmt.Sprintf(" ● Served Dir.     → %s", conf.ServeDir))
	log.Logger.Info(fmt.Sprintf(" ● Max Requests    → %d/minute", conf.MaxRequests))
	log.Logger.Info(fmt.Sprintf(" ● Ban Time        → %d minute(s)", conf.BanTime))
	log.Logger.Info(fmt.Sprintf(" ● Log To File     → %t", conf.LogToFile))
	log.Logger.Info(fmt.Sprintf(" ● Log Level       → %s", conf.LogLevel))
	log.Logger.Info(fmt.Sprintf(" ● Log File        → %s", conf.LogFile))
	log.Logger.Info(fmt.Sprintf(" ● Log File Format → %s", conf.LogFileFormat))
	log.Logger.Info(fmt.Sprintf(" ● Log Max Size    → %d MB", conf.LogMaxSize))
	log.Logger.Info(fmt.Sprintf(" ● Log Max Backups → %d", conf.LogMaxBackups))
	log.Logger.Info(fmt.Sprintf(" ● Log Max Age     → %d days", conf.LogMaxAge))
	log.Logger.Info("====================================================")
}
