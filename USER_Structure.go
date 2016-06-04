package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	// "errors"
	// "fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	// "strconv"
	"strings"
)

//// --------------------------
// User,
// Our local user.
// Has a true name.
////

// Type: User
// Our users with Name and Permission
type User struct {
	Name       string
	Email      string
	Permission int   `datastore:"-"`
	ID         int64 `datastore:"-"`
}

// Method: Key
// Implements Retrivable interface
func (u *User) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Users", "", id.(int64), nil)
}

func (u *User) StoreKey(k *datastore.Key) {
	u.ID = k.IntID()
}

// Internal Function
// Description:
//
// Returns:
//      user(User) - Prepared User
func MakeUser(name, email string) User {
	return User{
		Name:  name,
		Email: strings.ToLower(email),
	}
}
