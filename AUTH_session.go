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
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"
	"time"
)

var (
	// ErrNotLoggedIn is thrown when a session cannot find user OAuth information.
	ErrNotLoggedIn = errors.New("Session: Cannot Maintain Session, No Logged In User")

	// ErrTimedOut is thrown when a session is missing a validation cookie or memcache
	ErrTimedOut = errors.New("Session: User has timed out.")

	// Common duration time for session storage
	StorageDuration = time.Hour

	// Common key for a session validation cookie
	CookieKey = "Session"
)

// Internal Function, Outbound Service
// Description:
// This will create a google login url to have a user login using their gmail. It will then redirect back to an internal url of our choosing.
//
// Returns:
//		url(string) - Google login url with redirect.
func GetLoginURL(ctx context.Context, redirect string) string {
	login, _ := user.LoginURLFederated(ctx, redirect, "")
	return login
}

// Internal Function
// Description:
// If a session exists, this function will refresh all timers back to StorageDuration.
//
// Returns:
//		failure?(error) - Any errors are stored here if exists.
func MaintainSession(res http.ResponseWriter, req *http.Request) error {
	ctx := appengine.NewContext(req)

	if u := user.Current(ctx); u != nil {

		// Verify that a valid cookie exists on the client machine.
		if cookieValue, cErr := FromCookie(req, CookieKey); cErr != nil || cookieValue != u.ID {
			// No cookie? They must've timed out.
			return ErrTimedOut
		}

		// Verify that a copy of the local memory exits.
		if _, memErr := FromMemcache(ctx, u.Email); memErr != nil { // We can count on our local memory still being valid.
			// No Memory? Probably timed out.
			return ErrTimedOut
		}

		// If yes, reset those timers.
		UpdateCookie(res, req, CookieKey, StorageDuration)
		UpdateMemcache(ctx, u.Email, StorageDuration)
		return nil
	} else {
		// No user? Not logged in.
		return ErrNotLoggedIn
	}
}

// Internal Function
// Description:
// If a user has a valid OAuth token, this function will create a new session.
//
// Returns:
//		failure?(error) - Any errors are stored here if exists.
func CreateSession(res http.ResponseWriter, req *http.Request, dataToMemchache string) error {
	ctx := appengine.NewContext(req)
	if u := user.Current(ctx); u != nil {
		ToCookie(res, CookieKey, u.ID, StorageDuration)
		return ToMemcache(ctx, u.Email, dataToMemchache, StorageDuration)
	} else {
		return ErrNotLoggedIn
	}
}
