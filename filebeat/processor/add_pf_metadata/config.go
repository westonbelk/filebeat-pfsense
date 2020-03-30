package add_pf_metadata

import (
	"time"
)

type Config struct {
	CacheTTL       time.Duration   `config:"cache.ttl"`
	Name           string          `config:"name"`
}

func defaultConfig() Config {
	return Config {
		CacheTTL: 15* time.Second,
	}
}