package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
)

var (
	// ErrInvalidPermission, this error is thrown when a user fails a minimum permission level check.
	ErrInvalidPermission = errors.New("Permission Error: User Does Not Have Required Permission Level!")
)

//// --------------------------
// Permission Levels
// This collection details the different levels
// of permissions a user can hold, verification
// that a user meets a minimum permission requirement
// and insertion, retrieval, and deletion of permission
// levels into datastore.
// Permission Levels in datastore are on table Permissions
////

const (
	// Permission Levels.
	// These are const integers. Please refer to them always by name, never number.
	ReadPermissions  = iota
	EditPermissions  = iota
	WritePermissions = iota
	AdminPermissions = iota
)

// Internal Function
// Description:
// Given a response, request, and minimum permission level.
// This function will return a boolean if the current user
// does or does not meet the requirement.
//
// Returns:
//      valid?(bool) - True/False if user meets requirement
//      failure?(error) - Any errors are stored here if exists.
func HasPermission(res http.ResponseWriter, req *http.Request, minimumRequiredPermission int) (bool, error) {
	if sessErr := MaintainSession(res, req); sessErr != nil { // Must have a session
		return false, sessErr
	} else {
		ctx := appengine.NewContext(req)
		if u, permissionErr := GetPermissionUserFromSession(ctx); permissionErr != nil { // Must have a valid permission user.
			return false, permissionErr
		} else {
			if u.Permission < minimumRequiredPermission { // That permission user must be at least the minimum.
				return false, ErrInvalidPermission
			}
		}
	}
	return true, nil
}

func PutPermissionLevelToDatastore(ctx context.Context, keyname string, permLevel int) error {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	toDatastore := &struct{ PL int }{permLevel}
	_, putErr := datastore.Put(ctx, permkey, toDatastore)
	return putErr
}
func GetPermissionLevelFromDatastore(ctx context.Context, keyname string) (int, error) {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	pl := struct{ PL int }{}
	getErr := datastore.Get(ctx, permkey, &pl)
	return pl.PL, getErr
}
func RemovePermissionLevelFromDatastore(ctx context.Context, keyname string) error {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	return datastore.Delete(ctx, permkey)
}
