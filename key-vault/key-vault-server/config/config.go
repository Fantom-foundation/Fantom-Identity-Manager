package config

import (
	"github.com/spf13/viper"
)

// Config reader keys
const (
	keyPort          = "net.port"
	keyRootUrl       = "net.url"
	keyHydraAdminUrl = "hydra.admin-url"
	keyGraphqlPath   = "graphql.path"
)

// Master application configuration
type Config struct {
	// User facing port
	Port string
	// User facing address
	RootUrl string
	// Entrypoint manged by hydra server to join to
	HydraAdminUrl string
	// Graphql entrypoint path
	GraphqlEntrypoint string
}

func Load() (*Config, error) {
	// Initialize reader
	cfgReader := getReader()
	setDefaults(cfgReader)

	// Load each source of configuration
	if err := loadEnv(cfgReader); err != nil {
		return nil, err
	}

	// Build final configuration
	return &Config{
		Port:              cfgReader.GetString(keyPort),
		RootUrl:           cfgReader.GetString(keyRootUrl),
		HydraAdminUrl:     cfgReader.GetString(keyHydraAdminUrl),
		GraphqlEntrypoint: cfgReader.GetString(keyGraphqlPath),
	}, nil
}

func getReader() *viper.Viper {
	return viper.New() // Default reader
}
