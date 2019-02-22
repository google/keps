package cache

import (
	"github.com/shibukawa/configdir"
)

func Dir() string {
	// TODO remove this once properly integrated into settings.Runtime

	configDirs := configdir.New("kubernetes", "keps")
	cache := configDirs.QueryCacheFolder()

	return cache.Path
}
