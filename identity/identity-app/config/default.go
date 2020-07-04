package config

import (
	"github.com/spf13/viper"
)

func setDefaults(cfg *viper.Viper) {
	cfg.SetDefault(keyDebug, false)
	cfg.SetDefault(keyDebugDB, false)
	cfg.SetDefault(keyDebugCTX, false)
	cfg.SetDefault(keyCookieStoreKey, "01234567890123456789012345678901")
	cfg.SetDefault(keySessionStoreKey, "01234567890123456789012345678901")
	cfg.SetDefault(keySessionCookieName, "fantomSession")
	cfg.SetDefault(keyPort, "3000")
	cfg.SetDefault(keyRootUrl, "http://localhost")
	cfg.SetDefault(keyHydraAdminUrl, "http://localhosst:4445")
}
