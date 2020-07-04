package config

import (
	"github.com/spf13/viper"
)

// Config reader keys
const (
	keyDebug             = "debug.general"
	keyDebugDB           = "debug.db"
	keyDebugCTX          = "debug.ctx"
	keyCookieStoreKey    = "storekey.cookie"
	keySessionStoreKey   = "storekey.session"
	keySessionCookieName = "storekey.sessionname"
	keyPort              = "net.port"
	keyRootUrl           = "net.url"
	keyHydraAdminUrl     = "hydra.admin-url"
	keyDsn               = "dsn"
)

// Master application configuration
type Config struct {
	// General debug level
	Debug bool
	// Debug DB operation
	DebugDB bool
	// Debug operation context
	DebugCTX bool
	// Encrypted cookie key
	CookieStoreKey string
	// Encrypted session
	SessionStoreKey string
	// Session specification
	SessionCookieName string
	// User facing port
	Port string
	// User facing address
	RootUrl string
	// Entrypoint manged by hydra server to join to
	HydraAdminUrl string
	// Database connection specification
	DSN string
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
		Debug:             cfgReader.GetBool(keyDebug),
		DebugDB:           cfgReader.GetBool(keyDebugDB),
		DebugCTX:          cfgReader.GetBool(keyDebugCTX),
		CookieStoreKey:    cfgReader.GetString(keyCookieStoreKey),
		SessionStoreKey:   cfgReader.GetString(keySessionStoreKey),
		SessionCookieName: cfgReader.GetString(keySessionStoreKey),
		Port:              cfgReader.GetString(keyPort),
		RootUrl:           cfgReader.GetString(keyRootUrl),
		HydraAdminUrl:     cfgReader.GetString(keyHydraAdminUrl),
		DSN:               cfgReader.GetString(keyDsn),
	}, nil
}

func getReader() *viper.Viper {
	return viper.New() // Default reader
}
