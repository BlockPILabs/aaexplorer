package config

// TaskConfig defines the configuration options for the Task
type TaskConfig struct {
	Networks         []string `mapstructure:"networks" toml:"networks" json:"networks"`
	BlockScanThreads int      `mapstructure:"blockScanThreads" toml:"blockScanThreads" json:"blockScanThreads"`
}

// DefaultTaskConfig returns a default configuration for the Task
func DefaultTaskConfig() *TaskConfig {
	return &TaskConfig{
		Networks:         []string{},
		BlockScanThreads: 10,
	}
}
