package rest

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func lastModified(r *http.Request) time.Time {
	since := r.Header.Get("If-Modified-Since")
	if since == "" {
		return time.Time{}
	}

	client, err := time.Parse(http.TimeFormat, since)
	if err != nil {
		return time.Time{}
	}
	return client
}

func clearCache(ctx context.Context, key string) {
	if err := memcache.Delete(ctx, key); err != nil {
		log.Errorf(ctx, "Unable to remove cache key '%s':%s\n", key, err.Error())
	}
}

func clientEncache(ctx context.Context, page string, t time.Time) {
	i := memcache.Item{
		Key:   "LastModified/" + page,
		Value: []byte(t.Format(http.TimeFormat)),
	}
	if err := memcache.Set(ctx, &i); err != nil {
		log.Errorf(ctx, "Unable to encache page %s: %s\n", page, err.Error())
	}
}

func clientCached(ctx context.Context, w http.ResponseWriter, r *http.Request, page string) bool {
	i, err := memcache.Get(ctx, "LastModified/"+page)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false
		}
		log.Errorf(ctx, "Unable check the cache for page %s: %s\n", page, err.Error())
		return false
	}

	server, err := time.Parse(http.TimeFormat, string(i.Value))
	if err != nil {
		log.Errorf(ctx, "Unable to parse the cached date %s: %s\n", string(i.Value), err.Error())
		return false
	}
	client := lastModified(r)
	if !server.After(client) {
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	return false
}
