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
	"google.golang.org/appengine/user"
	"net/http"
	"time"
)

var (
	ErrNotLoggedIn  = errors.New("Session: Cannot Maintain Session, No Logged In User")
	ErrTimedOut     = errors.New("Session: User has timed out.")
	StorageDuration = time.Hour
	CookieKey       = "Session"
)

// GetLoginURL
// Internal, Outbound Service
// This will create a google login url to have a user login using their gmail. It will then redirect back to an internal url of our choosing.
func GetLoginURL(ctx context.Context, redirect string) string {
	login, _ := user.LoginURLFederated(ctx, redirect, "")
	return login
}

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

// Create Session:
// Requires a logged in user.
func CreateSession(res http.ResponseWriter, req *http.Request, dataToMemchache string) error {
	ctx := appengine.NewContext(req)
	if u := user.Current(ctx); u != nil {
		ToCookie(res, CookieKey, u.ID, StorageDuration)
		return ToMemcache(ctx, u.Email, dataToMemchache, StorageDuration)
	} else {
		return ErrNotLoggedIn
	}
}
