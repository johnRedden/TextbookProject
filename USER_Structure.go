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
// An appengine/user.User with Name and Permission
type PermissionUser struct {
	Name       string
	Permission int
	ID         string
	Email      string
}

// Method: ToString
// Converts a permission user into a marshalled string
func (u PermissionUser) ToString() string {
	return fmt.Sprintf("%s�%s�%d�%s", u.Name, u.Email, u.Permission, u.ID)
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
	if len(data) < 4 {
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
		ID:         data[3],
	}, nil
}

// Internal Function
// Description:
// Given an appengine/user.User, a name, and a permission level, will create a valid permission user.
//
// Returns:
//      user(PermissionUser) - Prepared PermissionUser
func MakePermissionUser(name string, permission int, u *user.User) PermissionUser {
	return PermissionUser{
		Name:       name,
		Permission: permission,
		Email:      strings.ToLower(u.Email),
		ID:         u.ID,
	}
}

//// --------------------------
// Permission User, Datastore
// This collection of functions handle the insertion,
// retrieval, and deletion of PermissionUsers from
// datastore.
// All PermissionUsers exist on table Users
////

func PutPermissionUserToDatastore(ctx context.Context, keyname string, pu *PermissionUser) error {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	_, putErr := datastore.Put(ctx, userkey, pu)
	return putErr
}
func GetPermissionUserFromDatastore(ctx context.Context, keyname string) (PermissionUser, error) {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	pu := PermissionUser{}
	getErr := datastore.Get(ctx, userkey, &pu)
	return pu, getErr
}
func RemovePermissionUserFromDatastore(ctx context.Context, keyname string) error {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	return datastore.Delete(ctx, userkey)
}
