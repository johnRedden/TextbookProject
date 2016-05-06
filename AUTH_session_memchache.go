package main

/*
AUTH_session_memchache.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"time"
)

// Internal Function
// Description:
// This function will add a key:value pair into memcache with a life of time.Duration.
//
// Returns:
//		failure?(error) - If any errors occur they exist here.
func ToMemcache(ctx context.Context, key string, value string, expiration time.Duration) error {
	mI := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: expiration,
	}
	return memcache.Set(ctx, mI)
}

// Internal Function
// Description:
// This function will retrieve a value that may exist in key from memcache.
//
// Returns:
//		value(string) - Value of key
//		failure?(error) - If any errors occur they exist here.
func FromMemcache(ctx context.Context, key string) (string, error) {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// Internal Function
// Description:
// This function will update a value with new time.Duration
//
// Returns:
//		failure?(error) - If any errors occur they exist here.
func UpdateMemcache(ctx context.Context, key string, expiration time.Duration) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	item.Expiration = expiration // Update the memcache expiration.
	return memcache.Set(ctx, item)
}

// Internal Function
// Description:
// This function will delete key from memcache
//
// Returns:
//		failure?(error) - If any errors occur they exist here.
func DeleteMemcache(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}
