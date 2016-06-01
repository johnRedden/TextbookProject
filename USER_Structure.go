package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"strconv"
	"strings"
)

var (
	// ErrPermissionUserMarshall, this error is thrown when there is an incorrect number of values to unmarshal into a Permission User
	ErrPermissionUserMarshall = errors.New("MarshallPermissionUser: Cannot Unmarshal String, Too Few Values")
)

//// --------------------------
// Permission User,
// A user with a permission value.
// Has a true name.
////

// Type: PermissionUser
// Our users with Name and Permission
type PermissionUser struct {
	Name       string
	Permission int
	Email      string
}

// Method: Key
// Implements Retrivable interface
func (u *PermissionUser) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Users", id.(string), 0, nil)
}

// Method: ToString
// Converts a permission user into a marshalled string
func (u PermissionUser) ToString() string {
	return fmt.Sprintf("%s�%s�%d", u.Name, u.Email, u.Permission)
}

// Internal Function
// Description:
// Function takes a stringed permission user and unmarshalls the values back into their original struct.
//
// Returns:
//      user(PermissionUser) - Unpacked PermissionUser
//      failure?(error) - Any errors are stored here if exists.
func MarshallPermissionUser(p string) (PermissionUser, error) {
	data := strings.Split(p, "�")
	if len(data) < 3 {
		return PermissionUser{}, ErrPermissionUserMarshall
	}
	permLevel, convErr := strconv.Atoi(data[2])
	if convErr != nil {
		return PermissionUser{}, convErr
	}
	return PermissionUser{
		Name:       data[0],
		Email:      data[1],
		Permission: permLevel,
	}, nil
}

// Internal Function
// Description:
// Given an appengine/user.User, a name, and a permission level, will create a valid permission user.
//
// Returns:
//      user(PermissionUser) - Prepared PermissionUser
func MakePermissionUser(name, email string, permission int) PermissionUser {
	return PermissionUser{
		Name:       name,
		Permission: permission,
		Email:      strings.ToLower(email),
	}
}
