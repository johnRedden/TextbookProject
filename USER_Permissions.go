package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
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
	u, sessErr := GetUserFromSession(res, req)
	if sessErr != nil {
		return false, sessErr
	}

	if u.Permission >= minimumRequiredPermission {
		return true, nil
	}
	return false, ErrInvalidPermission
}

func SetPermission(u *User, permission int) {
	u.Permission = permission
}
