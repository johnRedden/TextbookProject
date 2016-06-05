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
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
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
	UserKey int64
	Valid   time.Time
	ID      int64 `datastore:",noindex"`
}

// Method: Key
// Implements Retrivable interface
func (p *Session) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Sessions", "", id.(int64), nil)
}
func (p *Session) StoreKey(k *datastore.Key) {
	p.ID = k.IntID()
}

func NewSession(res http.ResponseWriter, req *http.Request, userId int64) (Session, error) {
	ctx := appengine.NewContext(req)
	s := Session{}
	if userId == int64(0) {
		return Session{}, ErrInvalidSession
	}
	s.UserKey = userId
	s.Valid = time.Now().Add(StorageDuration)

	// Optional: Limit Sessions to one?

	sK, putErr := retrievable.PlaceEntity(ctx, int64(0), &s)
	if putErr != nil {
		return Session{}, putErr
	}
	ToCookie(res, SessionCookie, fmt.Sprint(sK.IntID()), StorageDuration)
	return s, nil
}

func DeleteSession(res http.ResponseWriter, req *http.Request, id int64) error {
	ctx := appengine.NewContext(req)
	DeleteCookie(res, SessionCookie)
	return retrievable.DeleteEntity(ctx, (&Session{}).Key(ctx, id))
}

func GetSession(res http.ResponseWriter, req *http.Request) (Session, error) {
	val, cErr := FromCookie(req, SessionCookie) // Get session info from cookie
	if cErr != nil {
		return Session{}, ErrNotLoggedIn
	}

	id, convErr := strconv.ParseInt(val, 10, 64) // Change cookie val into key
	if convErr != nil {
		return Session{}, ErrInvalidSession
	}

	ctx := appengine.NewContext(req)
	s := Session{}
	getErr := retrievable.GetEntity(ctx, &s, id) // Get actual session from datastore
	if getErr != nil {
		return Session{}, ErrTimedOut
	}

	if s.Valid.After(time.Now()) { // is that session still valid?
		s.Valid = time.Now().Add(StorageDuration)
		_, putErr := PlaceInDatastore(ctx, s.ID, &s)
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
// Given an http call, this will retrieve the current user's User.
//
// Returns:
//      user(User) - Prepared PermissionUser
//      failure?(error) - Any errors are stored here if exists.
func GetUserFromSession(res http.ResponseWriter, req *http.Request) (User, error) {
	ctx := appengine.NewContext(req)

	s, sessErr := GetSession(res, req)
	if sessErr != nil {
		return User{}, sessErr
	}

	u := User{}
	getErr := retrievable.GetEntity(ctx, &u, s.UserKey)
	if getErr == nil {
		u.Permission = GetPermission(ctx, u.ID)
	}
	return u, getErr
}
