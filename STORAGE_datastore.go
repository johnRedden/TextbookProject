package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func PlaceInDatastore(ctx context.Context, key interface{}, source retrievable.Retrievable) (*datastore.Key, error) {
	return retrievable.PlaceInDatastore(ctx, key, source)
}
func GetFromDatastore(ctx context.Context, key interface{}, source retrievable.Retrievable) error {
	return retrievable.GetFromDatastore(ctx, key, source)
}
func DeleteFromDatastore(ctx context.Context, key interface{}, source retrievable.Retrievable) error {
	return retrievable.DeleteFromDatastore(ctx, key, source)
}
