// Authentication
// This package handles all advanced Authentication/Session Management.
package main

/*
AUTH_authentication.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"
	"strings"
	"time"
)

var (
	RegisterUUIDTime = time.Minute * time.Duration(3) // Three minutes until cookie deletes itself.
)

// Internal Function, Outbound Service
// Description:
// This will create a google login url to have a user login using their gmail. It will then redirect back to an internal url of our choosing.
//
// Returns:
//      url(string) - Google login url with redirect.
func GetOAuthURL(ctx context.Context, redirect string) string {
	login, _ := user.LoginURLFederated(ctx, redirect, "")
	return login
}

func NewUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

//// --------------------------
// Login Process,
// Login, Logout, and registration handlers.
////

// Call: /logout
// Description:
// Deletes the user's local session data.
// Will attempt to redirect the user to page at option:redirect if exists.
//
// Method: GET
// Results: HTTP.Redirect
// Mandatory Options:
// Optional Options: redirect
func AUTH_Logout_GET(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	DeleteSession(res)
	// ctx := appengine.NewContext(req)
	// lgo, _ := user.LogoutURL(ctx, "/"+req.FormValue("redirect"))
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

////----------------------//
// Register
////

// Call: /register
// Description:
// After user gains an OAuth token, will ask the user
// for additional information to create the user's
// PermissionUser.
// Will forward value of option:redirect.
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options: redirect
func AUTH_Register_GET(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	screen := User{}
	// Step 1: UUID
	registerUUID, cErr := FromCookie(req, "UUID")
	if cErr != nil { // No uuid, we'll make one
		registerUUID = NewUUID()
		ToCookie(res, "UUID", registerUUID, RegisterUUIDTime)
	} else {
		getErr := retrievable.GetFromMemcache(ctx, registerUUID, &screen)
		if ErrorPage(res, "UUID Error!", getErr) {
			return
		}
	}

	if authed, _ := FromCookie(req, "AUTHED"); authed != "yes" {
		// if req.FormValue("Oauth") != "yes" { // has oauth occured?
		// First time oauth. Go out to google.
		ToCookie(res, "AUTHED", "yes", RegisterUUIDTime)
		retrievable.PlaceInMemcache(ctx, registerUUID, screen, 0)
		lgn, _ := user.LoginURLFederated(ctx, "/register?Oauth=yes&redirect="+req.FormValue("redirect"), "")
		http.Redirect(res, req, lgn, http.StatusFound)
		return
	}

	DeleteCookie(res, "AUTHED")

	u := user.Current(ctx)
	if u == nil {
		ErrorPage(res, "OAuth Error!", ErrNotLoggedIn)
		return
	}
	screen.Email = strings.ToLower(u.Email)
	if u.Admin {
		screen.Permission = AdminPermissions
	}
	// Output and get Name
	retrievable.PlaceInMemcache(ctx, registerUUID, screen, 0)
	ServeTemplateWithParams(res, "register.html", screen)
	return
}
func AUTH_Register_GET_USINGUUID(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ToCookie(res, "UUID", params.ByName("UUID"), 0)
	http.Redirect(res, req, "/register?redirect="+req.FormValue("redirect"), http.StatusFound)
}

// Call: /register
// Description:
// After user submits the additional information required.
// This call will make the permission user and issue a new session.
// Will redirect to option:redirect if exists.
//
// Method: POST
// Results: HTTP.Redirect
// Mandatory Options:
// Optional Options: redirect
func AUTH_Register_POST(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	ruuid, cErr := FromCookie(req, "UUID")
	if ErrorPage(res, "UUID Error: Cannot find UUID cookie.", cErr) {
		return
	}

	uNew := &User{}
	getErr := retrievable.GetFromMemcache(ctx, ruuid, uNew)
	if ErrorPage(res, "UUID Error: Cannot find UUID memcache value.", getErr) {
		return
	}

	// Validate info
	if req.FormValue("Name") == "" || req.FormValue("Email") == "" {
		ErrorPage(res, "Incoming information invalid. Attempt to register again.", ErrInvalidSession)
		return
	}

	lvl := ReadPermissions
	if req.FormValue("Level") == "writer" {
		lvl = WritePermissions
	}

	// Set updated information
	uNew.Name = req.FormValue("Name")
	uNew.Permission = func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}(uNew.Permission, lvl)

	_, putErr := retrievable.PlaceEntity(ctx, int64(0), uNew)
	if ErrorPage(res, "User Creation Error", putErr) {
		return
	}

	createErr := CreateLogin(ctx, uNew)
	if ErrorPage(res, "Login Creation Error!", createErr) {
		return
	}
	http.SetCookie(res, MakeSessionCookie(uNew))
	DeleteCookie(res, "UUID")
	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
}

/////-----------------------//
// Login
////

// Call: /login
// Description:
// This is a login request from a user. This handles
// all requirements for login and will even redirect
// to a registration page if user is attempting a
// first time login.
//
// If a value is given to option:redirect, this call will
// attempt to redirect the user to said page.
//
// If option:changeuser is set to the string 'yes'
// then this function will force the user to reauthenticate.
//
// Method: GET
// Results: HTTP.Redirect
// Mandatory Options:
// Optional Options: redirect, changeuser
func AUTH_Login_GET(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// User has requested a login procedure.

	ctx := appengine.NewContext(req)

	if authed, _ := FromCookie(req, "AUTHED"); authed != "yes" {
		ToCookie(res, "AUTHED", "yes", RegisterUUIDTime)
		lgn, _ := user.LoginURLFederated(ctx, "/login?Oauth=yes&redirect="+req.FormValue("redirect"), "")
		http.Redirect(res, req, lgn, http.StatusFound)
		return
	}

	ou := user.Current(ctx)
	if ou == nil {
		DeleteCookie(res, "AUTHED") // spoofed AUTHED cookie. Delete this one and give an error page.
		ErrorPage(res, "Invalid Login State: OAuth", ErrNotLoggedIn)
		return
	}
	uid, getErr := GetUIDFromLogin(ctx, strings.ToLower(ou.Email))
	if getErr != nil || uid == 0 {
		http.Redirect(res, req, "/register?login=user_not_found&redirect="+req.FormValue("redirect"), http.StatusSeeOther)
		return
	}

	DeleteCookie(res, "AUTHED") // We've gone past the point of needing the auth cookie. go fourth.

	u := &User{}
	getuErr := retrievable.GetEntity(ctx, u, uid)
	if ErrorPage(res, "No Such User", getuErr) {
		return
	}

	http.SetCookie(res, MakeSessionCookie(u))
	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusFound)
}

// Call: /user
// TEMPORARY CALL
// For debug uses only.
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func AUTH_UserInfo(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Temporary GET
	// This is an excellent way to see just what session info we have and to verify login.
	u, err := GetUserFromSession(res, req)
	fmt.Fprint(res, `<html><plaintext>`)
	fmt.Fprint(res, "Error?\n", err, "::\n")
	fmt.Fprint(res, "User?\n", u, "::\n")
}
