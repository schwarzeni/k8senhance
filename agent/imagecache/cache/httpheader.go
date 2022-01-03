package cache

import (
	"net/http"
	"sync"
)

var (
	// httpHeaderCache cache response header for certain url
	httpHeaderCache     = map[string]http.Header{}
	httpHeaderCacheLock sync.RWMutex
)

func SetHTTPHeaderCache(url string, header http.Header) {
	httpHeaderCacheLock.Lock()
	defer httpHeaderCacheLock.Unlock()
	httpHeaderCache[url] = header
}

func HTTPHeaderCache(url string) (http.Header, bool) {
	httpHeaderCacheLock.RLock()
	defer httpHeaderCacheLock.RUnlock()
	res, ok := httpHeaderCache[url]
	return res, ok
}
