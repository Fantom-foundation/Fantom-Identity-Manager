package db

import (
	"errors"
	"github.com/volatiletech/authboss"
	"identity-app/config"
	"sync"
)

var (
	storers = make([]func() StorerBase, 0)
	mtx     sync.Mutex

	undefinedStorerError = errors.New("no storer defined for selected specification")
)

type StorerBase interface {
	Close()
	CanHandle(dsn string) bool
	FromConfig(cfg *config.Config) (StorerBase, error)
	authboss.ServerStorer
}

func RegisterStorer(provider func() StorerBase) {
	mtx.Lock()
	storers = append(storers, provider)
	mtx.Unlock()
}

func LoadStorer(cfg *config.Config) StorerBase {
	storer, err := getCompatibleStorer(cfg.DSN)
	if err != nil {
		return nil
	}
	storer, err = storer.FromConfig(cfg)
	if err != nil {
		return nil
	}
	return storer
}

func getCompatibleStorer(dsn string) (StorerBase, error) {
	for _, provider := range storers {
		storer := provider()
		if storer.CanHandle(dsn) {
			return storer, nil
		}
	}
	return nil, undefinedStorerError
}
