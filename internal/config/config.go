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
	ServeDir      string
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
	if net.ParseIP(c.Host) == nil {
		_, err := net.LookupHost(c.Host)
		if err != nil {
			return fmt.Errorf("invalid host: %s", c.Host)
		}
	}

	// Port validation
	if c.Port <= 1024 || c.Port >= 65535 {
		return fmt.Errorf("port must be between 1025 and 65534: %d", c.Port)
	}

	// The directory must exist
	if _, err := os.Stat(c.ServeDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", c.ServeDir)
	}

	// Max Requests validation (> 0 and < 100)
	if c.MaxRequests <= 0 || c.MaxRequests >= 100 {
		return fmt.Errorf("max requests should be between 1 and 99: %d", c.MaxRequests)
	}

	// Ban Time validation (should be > 0)
	if c.BanTime <= 0 {
		return fmt.Errorf("ban time should be greater than 0: %d", c.BanTime)
	}

	// Log level
	validLevels := []string{"info", "debug", "warn", "error"}
	if !slices.Contains(validLevels, c.LogLevel) {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Log file format (text or json)
	validFormat := []string{"text", "json"}
	if !slices.Contains(validFormat, c.LogFileFormat) {
		return fmt.Errorf("invalid log file format: %s", c.LogFileFormat)
	}

	// Log max size (should be > 0)
	if c.LogMaxSize < 1 {
		return fmt.Errorf("log max size should be greater than 0: %d", c.LogMaxSize)
	}

	// Log max backups (should be > 0)
	if c.LogMaxBackups < 1 {
		return fmt.Errorf("log max backups should be greater than 0: %d", c.LogMaxBackups)
	}

	// Log max age (should be > 0)
	if c.LogMaxAge < 1 {
		return fmt.Errorf("log max backups should be greater than 0: %d", c.LogMaxAge)
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
		ServeDir:      DefaultServeDir,
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
