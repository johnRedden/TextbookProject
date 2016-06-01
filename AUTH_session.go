// Session
// This package handles all basic Authentication/Session Management.
package main

/*
AUTH_session.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotLoggedIn    = errors.New("Session: Cannot Maintain Session, No Logged In User") // ErrNotLoggedIn is thrown when a session cannot find user information.
	ErrTimedOut       = errors.New("Session: User has timed out.")                        // ErrTimedOut is thrown when a session is missing a validation cookie or memcache
	ErrInvalidSession = errors.New("Session: Invalid session information.")
	// Common duration time for session storage
	StorageDuration = time.Hour * time.Duration(24)

	// Common key for a session validation cookie
	SessionCookie = "Session"
)

type Session struct {
	UserKey string
	Valid   time.Time
	ID      int64 `datastore:",noindex"`
}

// Method: Key
// Implements Retrivable interface
func (p *Session) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Sessions", "", id.(int64), nil)
}

func NewSession(res http.ResponseWriter, req *http.Request, email string) (Session, error) {
	ctx := appengine.NewContext(req)
	s := Session{}
	s.UserKey = email
	s.Valid = time.Now().Add(StorageDuration)

	// Optional: Limit Sessions to one?

	sK, putErr := PlaceInDatastore(ctx, int64(0), &s)
	if putErr != nil {
		return Session{}, nil, putErr
	}
	s.ID = sK.IntID()
	ToCookie(res, SessionCookie, fmt.Sprint(sK.IntID()), StorageDuration)
	return s, sK, nil
}

func DeleteSession(res http.ResponseWriter, req *http.Request, id int64) error {
	ctx := appengine.NewContext(req)
	DeleteCookie(res, SessionCookie)
	return DeleteFromDatastore(ctx, id, &Session{})
}

func GetSession(res http.ResponseWriter, req *http.Request) (Session, error) {
	val, cErr := FromCookie(req, key) // Get session info from cookie
	if cErr != nil {
		return Session{}, ErrNotLoggedIn
	}

	id, convErr := strconv.ParseInt(val, 10, 64) // Change cookie val into key
	if convErr != nil {
		return Session{}, ErrInvalidSession
	}

	s := Session{}
	getErr := GetFromDatastore(ctx, id, &s) // Get actual session from datastore
	if getErr != nil {
		return Session{}, ErrTimedOut
	}
	s.ID = id

	if s.Valid.After(time.Now()) { // is that session still valid?
		s.Valid = time.Now().Add(StorageDuration)
		putErr := PlaceInDatastore(ctx, s.ID, &s)
		if putErr != nil {
			return s, putErr
		}
	}

	return s, nil
}

// Internal Function
// Description:
// Given an http call, this will return a boolean true if there is a valid session.
//
// Returns:
//      exists(boolean) - true if session exists
func HasSession(res http.ResponseWriter, req *http.Request) bool {
	_, err := GetSession(res, req)
	return err == nil
}

// Internal Function
// Description:
// Given an http call, this will retrieve the current user's PermissionUser.
//
// Returns:
//      user(PermissionUser) - Prepared PermissionUser
//      failure?(error) - Any errors are stored here if exists.
func GetUserFromSession(res http.ResponseWriter, req *http.Request) (PermissionUser, error) {
	ctx := appengine.NewContext(req)

	s, sessErr := GetSession(res, req)
	if sessErr != nil {
		return PermissionUser{}, sessErr
	}

	u := PermissionUser{}
	getErr := GetFromDatastore(ctx, s.UserKey, &u)
	return u, getErr
}
