package config

import (
	"github.com/spf13/viper"
)

func getParameterNameMapping() map[string]string {
	return map[string]string{
		"DEBUG":     keyDebug,
		"DEBUG_DB":  keyDebugDB,
		"DEBUG_CTX": keyDebugCTX,
	}
}

func loadEnv(cfg *viper.Viper) error {
	// Load each specified simple parameter
	for envKey, cfgKey := range getParameterNameMapping() {
		if err := cfg.BindEnv(cfgKey, envKey); err != nil {
			return err
		}
	}
	return nil
}
