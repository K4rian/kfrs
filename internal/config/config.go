package config

import (
	"fmt"
	"net"
	"os"
	"slices"
	"sync"
)

type Config struct {
	Host          string
	Port          int
	Directory     string
	MaxRequests   int
	BanTime       int
	LogToFile     bool
	LogLevel      string
	LogFile       string
	LogFileFormat string
	LogMaxSize    int
	LogMaxBackups int
	LogMaxAge     int
}

func (c *Config) Validate() error {
	// Host validation (v4, v6, or domain)
	if net.ParseIP(conf.Host) == nil {
		_, err := net.LookupHost(conf.Host)
		if err != nil {
			return fmt.Errorf("invalid host: %s", conf.Host)
		}
	}

	// Port validation
	if conf.Port <= 1024 || conf.Port >= 65535 {
		return fmt.Errorf("port must be between 1025 and 65534: %d", conf.Port)
	}

	// The directory must exist
	if _, err := os.Stat(conf.Directory); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", conf.Directory)
	}

	// Max Requests validation (> 0 and < 100)
	if conf.MaxRequests <= 0 || conf.MaxRequests >= 100 {
		return fmt.Errorf("max requests should be between 1 and 99: %d", conf.MaxRequests)
	}

	// Ban Time validation (should be > 0)
	if conf.BanTime <= 0 {
		return fmt.Errorf("ban time should be greater than 0: %d", conf.BanTime)
	}

	// Log level
	validLevels := []string{"info", "debug", "warn", "error"}
	if !slices.Contains(validLevels, conf.LogLevel) {
		return fmt.Errorf("invalid log level: %s", conf.LogLevel)
	}

	// Log file format (text or json)
	validFormat := []string{"text", "json"}
	if !slices.Contains(validFormat, conf.LogFileFormat) {
		return fmt.Errorf("invalid log file format: %s", conf.LogFileFormat)
	}

	// Log max size (should be > 0)
	if conf.LogMaxSize < 1 {
		return fmt.Errorf("log max size should be greater than 0: %d", conf.LogMaxSize)
	}

	// Log max backups (should be > 0)
	if conf.LogMaxBackups < 1 {
		return fmt.Errorf("log max backups should be greater than 0: %d", conf.LogMaxBackups)
	}

	// Log max age (should be > 0)
	if conf.LogMaxAge < 1 {
		return fmt.Errorf("log max backups should be greater than 0: %d", conf.LogMaxAge)
	}
	return nil
}

var (
	conf *Config
	once sync.Once
)

func Get() *Config {
	once.Do(initConfig)
	return conf
}

func initConfig() {
	conf = &Config{
		Host:          DefaultHost,
		Port:          DefaultPort,
		Directory:     DefaultDirectory,
		MaxRequests:   DefaultMaxRequests,
		BanTime:       DefaultBanTime,
		LogToFile:     DefaultLogToFile,
		LogLevel:      DefaultLogLevel,
		LogFile:       DefaultLogFile,
		LogFileFormat: DefaultLogFileFormat,
		LogMaxSize:    DefaultLogMaxSize,
		LogMaxBackups: DefaultLogMaxBackups,
		LogMaxAge:     DefaultLogMaxAge,
	}
}
