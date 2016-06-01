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
	"strings"
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

	ctx := appengine.NewContext(req)
	if sessErr := MaintainSession(res, req); sessErr == ErrNotLoggedIn || req.FormValue("changeuser") == "yes" {
		// User is not logged in.
		// Force them to the google login page before coming back here.

		http.Redirect(res, req, GetLoginURL(ctx, "/login?redirect="+req.FormValue("redirect")), http.StatusTemporaryRedirect)
		return
	} else if sessErr == ErrTimedOut {
		// User has an oauth key.
		// Likely returned from ouath.
		u := user.Current(ctx)
		pu, getErr := GetPermissionUserFromDatastore(ctx, strings.ToLower(u.Email))
		if getErr != nil {
			// They do not have a registered permission user.
			// Kick them over to register.
			http.Redirect(res, req, "/register?redirect="+req.FormValue("redirect"), http.StatusTemporaryRedirect)
			return
		}
		// we now have their user information.
		sessErr := CreateSession(res, req, pu.ToString())
		if sessErr != nil {
			http.Error(res, sessErr.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Session is live.
	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
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
	DeleteCookie(res, CookieKey) // Have the user invalidate their local login token.
	DeleteCookie(res, "ACSID")   // To be nice, we'll also delete the oauth token from google.
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
	u := user.Current(ctx)
	if u == nil {
		// They are not logged in!
		// We'll kick them over to google for ouath.
		http.Redirect(res, req, GetLoginURL(ctx, "/register?redirect="+req.FormValue("redirect")), http.StatusTemporaryRedirect)
		return
	}
	ServeTemplateWithParams(res, "registerUser.html", strings.ToLower(u.Email))
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
	u := user.Current(ctx)
	if u == nil {
		// They are not logged in!
		// No cross site attacks!
		http.Error(res, ErrNotLoggedIn.Error(), http.StatusTeapot)
		return
	}

	// Now that we're all satisfied. Lets grab that info.
	uName := req.FormValue("Name")

	// TODO: Require that a name is more than empty?

	// Permissions Module
	perms := ReadPermissions // Default Permissions.
	if pl, getErr := GetPermissionLevelFromDatastore(ctx, strings.ToLower(u.Email)); getErr == nil {
		perms = pl // use the already determined permission level.
	} else {
		// Ensure that user is not a new administrator.
		if u.Admin {
			perms = AdminPermissions
		}
	}
	putLErr := PutPermissionLevelToDatastore(ctx, strings.ToLower(u.Email), perms)
	HandleError(res, putLErr)

	// Make user and add them to the datastore.
	permU := MakePermissionUser(uName, perms, u)
	putErr := PutPermissionUserToDatastore(ctx, strings.ToLower(u.Email), &permU)
	HandleError(res, putErr)

	// Now we make that session
	sessErr := CreateSession(res, req, permU.ToString())
	HandleError(res, sessErr)

	http.Redirect(res, req, "/"+req.FormValue("redirect"), http.StatusSeeOther)
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
	if err := MaintainSession(res, req); err == nil {
		ctx := appengine.NewContext(req)
		if pVal, err := GetPermissionUserFromSession(ctx); err == nil {
			fmt.Fprint(res, `<p>`, pVal, `</p><br>`)
		} else {
			fmt.Fprint(res, `<p>`, err.Error(), `</p><br>`)
		}
		return
	} else {
		fmt.Fprint(res, `<!DOCTYPE html><html><head><title></title></head><body> Cannot Maintain session`+err.Error()+`</body></html>`)
		return
	}

}
