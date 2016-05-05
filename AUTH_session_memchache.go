package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"time"
)

func ToMemcache(ctx context.Context, key string, value string, expiration time.Duration) error {
	mI := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expiration,
	}
	return memcache.Set(ctx, mI)
}

func FromMemcache(ctx context.Context, key string) (string, error) {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func UpdateMemcache(ctx context.Context, key string, expiration time.Duration) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	item.Expiration = expiration // Update the memcache expiration.
	return memcache.Set(ctx, item)
}

func DeleteMemchache(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}
