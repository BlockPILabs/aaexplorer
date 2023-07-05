package config

import (
	"fmt"
	aimos "github.com/BlockPILabs/aa-scan/internal/os"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"path/filepath"
)

// DefaultDirPerm is the default permissions used when creating directories.
const DefaultDirPerm = 0700

var (
	DefaultHomeDir   = ".aim"
	defaultConfigDir = "config"
	defaultDataDir   = "data"

	defaultConfigFileName = "config.toml"

	defaultConfigFilePath = filepath.Join(defaultConfigDir, defaultConfigFileName)
)

const (

	// FuzzModeDrop is a mode in which we randomly drop reads/writes, connections or sleep
	FuzzModeDrop = iota
	// FuzzModeDelay is a mode in which we randomly sleep
	FuzzModeDelay

	// LogFormatPlain is a format for colored text
	LogFormatPlain = "plain"
	// LogFormatJSON is a format for json output
	LogFormatJSON = "json"

	// DefaultLogLevel defines a default log level as INFO.
	DefaultLogLevel = "info"

	// Mempool versions. V1 is prioritized mempool (deprecated), v0 is regular mempool.
	// Default is v0.
	MempoolV0 = "v0"
	MempoolV1 = "v1"
)

// BaseConfig defines the base configuration
type BaseConfig struct { //nolint: maligned

	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`

	// Output format: 'plain' (colored text) or 'json'
	LogFormat string `mapstructure:"log_format"`

	// Output level for logging
	LogLevel string `mapstructure:"log_level"`
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg BaseConfig) ValidateBasic() error {
	switch cfg.LogFormat {
	case LogFormatPlain, LogFormatJSON:
	default:
		return errors.New("unknown log_format (must be 'plain' or 'json')")
	}
	return nil
}

// DefaultBaseConfig returns a default base configuration
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		RootDir:   "",
		LogFormat: LogFormatPlain,
		LogLevel:  DefaultLogLevel,
	}
}

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for services
	Api *ApiConfig `mapstructure:"api"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Api:        DefaultApiConfig(),
	}
}

// EnsureRoot creates the root, config, and data directories if they don't exist,
// and panics if it fails.
func EnsureRoot(rootDir string) {
	if err := aimos.EnsureDir(rootDir, DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := aimos.EnsureDir(filepath.Join(rootDir, defaultConfigDir), DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := aimos.EnsureDir(filepath.Join(rootDir, defaultDataDir), DefaultDirPerm); err != nil {
		panic(err.Error())
	}

	configFilePath := filepath.Join(rootDir, defaultConfigFilePath)

	// Write default config file if missing.
	if !aimos.FileExists(configFilePath) {
		writeDefaultConfigFile(configFilePath)
	}
}

func writeDefaultConfigFile(configFilePath string) {
	WriteConfigFile(configFilePath, DefaultConfig())
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config *Config) {

	bts, err := toml.Marshal(config)
	if err != nil {
		return
	}

	aimos.MustWriteFile(configFilePath, bts, 0644)
}

// SetRoot sets the RootDir for all Config structs
func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	cfg.Api.RootDir = root
	return cfg
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg *Config) ValidateBasic() error {
	if err := cfg.BaseConfig.ValidateBasic(); err != nil {
		return err
	}
	if err := cfg.Api.ValidateBasic(); err != nil {
		return fmt.Errorf("error in [rpc] section: %w", err)
	}

	return nil
}

func (cfg *Config) CheckDeprecated() []string {
	var warnings []string
	return warnings
}
