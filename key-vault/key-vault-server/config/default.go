package config

import (
	"github.com/spf13/viper"
)

func setDefaults(cfg *viper.Viper) {
	cfg.SetDefault(keyPort, "3030")
	cfg.SetDefault(keyRootUrl, "http://localhost")
	cfg.SetDefault(keyHydraAdminUrl, "http://localhosst:4445")
}
