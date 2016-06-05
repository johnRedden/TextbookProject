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
	// "google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
	"net/http"
	// "strings"
	// "time"
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

func NewLoginUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

//// --------------------------
// Login Process,
// Login, Logout, and registration handlers.
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
	// Attempt to gather user info.
	ServeTemplateWithParams(res, "login.html", nil)
	// http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
}

func AUTH_LOGIN_POST(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	fmt.Fprint(res, "<html><plaintext>", req.FormValue("Name"), "\n", req.FormValue("Email"), "\n", req.FormValue("Password"), "\n")

	// http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
}

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
	// DeleteCookie(res, CookieKey) // Have the user invalidate their local login token.
	// DeleteCookie(res, "ACSID")   // To be nice, we'll also delete the oauth token from google.
	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
}

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

	registerUUID, cErr := FromCookie(req, "UUID")
	if cErr != nil {
		registerUUID = NewLoginUUID() // No uuid, we'll make one
		ToCookie(res, "UUID", registerUUID, StorageDuration)
	} else {
		getErr := retrievable.GetFromMemcache(ctx, registerUUID, &screen)
		if getErr != nil { // get what we stored for this uuid
			fmt.Fprintln(res, "<html><plaintext>UUID ERROR!", getErr)
			return
		}
	}

	if req.FormValue("Oauth") == "now" {
		lgn, _ := user.LoginURL(ctx, "/register?Oauth=yes")
		lgo, _ := user.LogoutURL(ctx, lgn)
		http.Redirect(res, req, lgo, http.StatusFound)
		return
	}

	if req.FormValue("Oauth") == "yes" {
		u := user.Current(ctx)
		screen.Email = u.Email
		if u.Admin {
			screen.Permission = AdminPermissions
		}
	}

	retrievable.PlaceInMemcache(ctx, registerUUID, screen, 0)
	ServeTemplateWithParams(res, "register.html", screen)
}
func AUTH_Register_GET_USINGUUID(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ToCookie(res, "UUID", params.ByName("UUID"), 0)
	http.Redirect(res, req, "/register", http.StatusFound)
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
	fmt.Fprint(res, "<html><plaintext>")
	ctx := appengine.NewContext(req)

	name := req.FormValue("Name")
	email := req.FormValue("Email")
	pswd := req.FormValue("Password")

	ruuid, cErr := FromCookie(req, "UUID")
	if cErr != nil {
		fmt.Fprintln(res, "Error with uuid cookie!")
		fmt.Fprintln(res, cErr, "\n")
	}

	uNew := &User{}
	getErr := retrievable.GetFromMemcache(ctx, ruuid, uNew)
	if getErr != nil {
		fmt.Fprintln(res, "Error with UUID mamcache!", getErr, "\n")
	}

	uNew.Email = email
	uNew.Name = name

	fmt.Fprintln(res, "Name:", name)
	fmt.Fprintln(res, "Email:", email)
	fmt.Fprintln(res, "Password:", pswd)
	fmt.Fprintln(res, "Permission:", uNew.Permission)
	fmt.Fprintln(res, "Now Preforming Creation actions.")

	uNew, createErr := CreateUserFromLogin(ctx, uNew.Email, pswd, uNew)
	if createErr != nil {
		fmt.Fprintln(res, "Creation Error!", createErr)
		return
	}

	fmt.Fprintln(res, "Creation Done:", uNew)

	permErr := SetPermission(ctx, uNew.ID, uNew.Permission)
	if permErr != nil {
		fmt.Fprintln(res, "Permission Error!", uNew)
		return
	}
	fmt.Fprintln(res, "Permission Done:", uNew)

	sess, sErr := NewSession(res, req, uNew.ID)
	if sErr != nil {
		fmt.Fprintln(res, "Sesssion Error!", sErr)
		return
	}
	ToCookie(res, SessionCookie, fmt.Sprint(sess.ID), StorageDuration)
	fmt.Fprintln(res, "Session Done:", sess)
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
