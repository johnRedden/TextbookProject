package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
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

type Permission struct {
	Permission int
}

// Method: Key
// Implements Retrivable interface
func (p *Permission) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Permissions", "", id.(int64), nil)
}

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
	u, sessErr := GetUserFromSession(res, req)
	if sessErr != nil {
		return false, sessErr
	}
	ctx := appengine.NewContext(req)
	perm := Permission{}
	getErr := retrievable.GetEntity(ctx, &perm, u.ID)
	if getErr != nil {
		// No Permissions? Default.
		perm = Permission{
			Permission: ReadPermissions,
		}
	}

	if perm.Permission >= minimumRequiredPermission {
		return true, nil
	}
	return false, ErrInvalidPermission
}

func NewPermission(appu *user.User, recommended int) Permission {
	if appu != nil && appu.Admin {
		return Permission{AdminPermissions}
	}
	return Permission{recommended}
}

func SetPermission(ctx context.Context, uID int64, perm *Permission) error {
	_, putErr := retrievable.PlaceEntity(ctx, uID, perm)
	return putErr
}

func GetPermission(ctx context.Context, uID int64) int {
	perm := Permission{}
	getErr := retrievable.GetEntity(ctx, &perm, uID)
	if getErr != nil {
		return ReadPermissions
	}
	return perm.Permission
}
