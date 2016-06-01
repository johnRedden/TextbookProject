package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// This is an interface to help a struct
// directly interface with datastore.
//
// Any struct that implements Retrivable will
// have access to the functionality of placement
// in datastore, retrival from datastore, and deletion
type Retrievable interface {
	// Key will create a datastore key based on
	// un-exported information from an interface
	// This is a method tied to the struct pointer
	Key(context.Context, interface{}) *datastore.Key
}

func PlaceInDatastore(ctx context.Context, key interface{}, source Retrievable) (*datastore.Key, error) {
	uk := source.Key(ctx, key)
	return datastore.Put(ctx, uk, source)
}
func GetFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	return datastore.Get(ctx, uk, source)
}
func DeleteFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	return datastore.Delete(ctx, uk)
}
