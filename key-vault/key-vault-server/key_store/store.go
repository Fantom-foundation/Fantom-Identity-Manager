package key_store

import (
	"errors"
	"key-vault-server/config"
	"key-vault-server/logging"
	"sync"
)

var (
	stores = make([]func() StoreBase, 0)
	mtx    sync.Mutex

	undefinedStoreError = errors.New("no store defined for selected specification")
)

type StoreBase interface {
	Close()
	CanHandle(dsn string) bool
	FromConfig(cfg *config.Config, log *logging.Logger) (StoreBase, error)

	// TODO: There should be a key path as another parameter
	getKey(userId string) string
}

func RegisterStore(provider func() StoreBase) {
	mtx.Lock()
	stores = append(stores, provider)
	mtx.Unlock()
}

func LoadStorer(cfg *config.Config, logger *logging.Logger) (StoreBase, error) {
	storer, err := getCompatibleStore(cfg.DSN)
	if err != nil {
		return nil, err
	}
	storer, err = storer.FromConfig(cfg, logger)
	if err != nil {
		return nil, err
	}
	return storer, nil
}

func getCompatibleStore(dsn string) (StoreBase, error) {
	for _, provider := range stores {
		storer := provider()
		if storer.CanHandle(dsn) {
			return storer, nil
		}
	}
	return nil, undefinedStoreError
}
