package apis

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var _cache = cache.New(15*time.Minute, 5*time.Minute)
