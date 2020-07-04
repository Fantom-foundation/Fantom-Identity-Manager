package config

import (
	"github.com/spf13/viper"
)

func getParameterNameMapping() map[string]string {
	return map[string]string{
		"DEBUG":             keyDebug,
		"DEBUG_DB":          keyDebugDB,
		"DEBUG_CTX":         keyDebugCTX,
		"COOKIE_STORE_KEY":  keyCookieStoreKey,
		"SESSION_STORE_KEY": keySessionStoreKey,
		"PORT":              keyPort,
		"ROOT_URL":          keyRootUrl,
		"HYDRA_ADMIN_URL":   keyHydraAdminUrl,
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
