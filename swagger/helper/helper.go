package helper

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheParams struct {
	ExpirationTime int
	AttemptsLimit  int
}

var CacheEnv *CacheParams = &CacheParams{}

func GetEnvVariables(expiration, attempts int) {
	CacheEnv.AttemptsLimit = attempts
	CacheEnv.ExpirationTime = expiration
}

var cacheInit = cache.New(0, 1*time.Hour)

type LoginCounter struct {
	Number int
}

func SetCache(key string, loginCounter interface{}) {
	cacheInit.Set(key, loginCounter, time.Duration(CacheEnv.ExpirationTime)*time.Hour)
}

var loginCounter LoginCounter

func GetCache(key string) (LoginCounter, bool) {
	data, found := cacheInit.Get(key)
	if found {
		loginCounter = data.(LoginCounter)
	}
	return loginCounter, found
}

func (l LoginCounter) LimitReached() bool {
	return l.Number == CacheEnv.AttemptsLimit
}
