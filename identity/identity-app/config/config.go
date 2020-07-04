package config

import (
	"github.com/spf13/viper"
)

// Config reader keys
const (
	keyDebug    = "debug.general"
	keyDebugDB  = "debug.db"
	keyDebugCTX = "debug.ctx"
)

// Master application configuration
type Config struct {
	// General debug level
	Debug bool
	// Debug DB operation
	DebugDB bool
	// Debug operation context
	DebugCTX bool
}

func Load() (*Config, error) {
	// Initialize reader
	cfgReader := getReader()

	// Load each source of configuration
	if err := loadEnv(cfgReader); err != nil {
		return nil, err
	}

	// Build final configuration
	return &Config{
		Debug:    cfgReader.GetBool(keyDebug),
		DebugDB:  cfgReader.GetBool(keyDebugDB),
		DebugCTX: cfgReader.GetBool(keyDebugCTX),
	}, nil
}

func getReader() *viper.Viper {
	return viper.New()
}
