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
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"net/http"
	// "strings"
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
	u := user.Current(ctx) // we may be serving them an oauth user.
	ServeTemplateWithParams(res, "register.html", u)
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
	fmt.Fprint(res, "<html><plaintext>")
	if req.FormValue("Oauth") == "yes" {
		fmt.Fprintln(res, "Oauth has been requested, at this point")
		fmt.Fprintln(res, "in time a url has been submitted to google.")
		fmt.Fprintln(res, "Refresh the page and assume a google oauth token has been created.")

		fmt.Fprintln(res, "\n", GetOAuthURL(ctx, "/register"))
		return
	}
	fmt.Fprintln(res, "In this case, information is coming in:")
	fmt.Fprintln(res, req.FormValue("Name"), "\n", req.FormValue("Email"), "\n", req.FormValue("Password"), "\n")
	fmt.Fprintln(res, user.Current(ctx))
	fmt.Fprintln(res, "We would preform register actions and log in.")
	// http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
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
