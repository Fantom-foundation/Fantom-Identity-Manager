package config

import (
	"github.com/spf13/viper"
)

func getParameterNameMapping() map[string]string {
	return map[string]string{
		"PORT":            keyPort,
		"ROOT_URL":        keyRootUrl,
		"HYDRA_ADMIN_URL": keyHydraAdminUrl,
		"GRAPHQL_PATH":    keyGraphqlPath,
		"DSN":             keyDsn,
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
