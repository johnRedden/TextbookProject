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
	"google.golang.org/appengine"
	"net/http"
	"strconv"
)

var (
	ErrNotLoggedIn    = errors.New("Session: No Logged In User")   // ErrNotLoggedIn is thrown when a session cannot find user information.
	ErrTimedOut       = errors.New("Session: User has timed out.") // ErrTimedOut is thrown when a session is missing a validation cookie or memcache
	ErrInvalidSession = errors.New("Session: Invalid session information.")
)

func SetSesson(res http.ResponseWriter, u *User) {
	http.SetCookie(res, &http.Cookie{
		Name:  "Session",
		Value: fmt.Sprint(u.ID),
		Path:  "/",
	})
}

func MakeSessionCookie(u *User) *http.Cookie {
	return &http.Cookie{
		Name:  "Session",
		Value: fmt.Sprint(u.ID),
		Path:  "/",
	}
}

func DeleteSession(res http.ResponseWriter) {
	http.SetCookie(res, &http.Cookie{
		Name:   "Session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func GetUserFromSession(res http.ResponseWriter, req *http.Request) (*User, error) {
	ctx := appengine.NewContext(req)

	c, cerr := req.Cookie("Session")
	if cerr != nil {
		return &User{}, ErrNotLoggedIn
	}

	uid, convErr := strconv.ParseInt(c.Value, 10, 64)
	if convErr != nil || uid == 0 {
		return &User{}, ErrInvalidSession
	}

	u := &User{}
	rerr := retrievable.GetEntity(ctx, u, uid)
	if rerr != nil {
		return &User{}, rerr
	}

	return u, nil
}
