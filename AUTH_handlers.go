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
	screen := struct {
		User
		Redirect string `datastore:"-"`
	}{
		User{},
		req.FormValue("redirect"),
	}

	if req.FormValue("Oauth") == "now" {
		ctx := appengine.NewContext(req)
		lgn, _ := user.LoginURL(ctx, "/login?Oauth=yes&redirect="+req.FormValue("redirect"))
		lgo, _ := user.LogoutURL(ctx, lgn)
		http.Redirect(res, req, lgo, http.StatusFound)
		return
	}

	if req.FormValue("Oauth") == "yes" {
		ctx := appengine.NewContext(req)
		u := user.Current(ctx)
		screen.Email = u.Email
		screen.Redirect, _ = user.LogoutURL(ctx, "/"+req.FormValue("redirect"))
	}

	ServeTemplateWithParams(res, "login.html", screen)
}

func AUTH_LOGIN_POST(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	email := req.FormValue("Email")
	pswd := req.FormValue("Password")

	uid, loginErr := GetUserIDFromLogin(ctx, email, pswd)
	if ErrorPage(res, "Login Validation Error", loginErr) {
		return
	}

	_, sErr := NewSession(res, req, uid)
	if ErrorPage(res, "Session Creation Error!", sErr) {
		return
	}

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
	DeleteSession(res, req) // Remove our local session information.
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
	screen := struct {
		User
		Redirect string `datastore:"-"`
	}{
		User{},
		"",
	}
	// Step 1: UUID
	registerUUID, cErr := FromCookie(req, "UUID")
	if cErr != nil {
		registerUUID = NewLoginUUID() // No uuid, we'll make one
		ToCookie(res, "UUID", registerUUID, StorageDuration)
	} else {
		getErr := retrievable.GetFromMemcache(ctx, registerUUID, &screen)
		if ErrorPage(res, "UUID Error!", getErr) {
			return
		}
	}

	screen.Redirect = req.FormValue("redirect")

	// Step 2.1: Oauth asked for?
	if req.FormValue("Oauth") == "now" {
		lgn, _ := user.LoginURL(ctx, "/register?Oauth=yes&redirect="+req.FormValue("redirect"))
		lgo, _ := user.LogoutURL(ctx, lgn)
		http.Redirect(res, req, lgo, http.StatusFound)
		return
	}

	// Step 2.2: Oauth returned?
	if req.FormValue("Oauth") == "yes" {
		u := user.Current(ctx)
		screen.Email = u.Email
		// A bug has been possibly found with this? Users that did not sign in with google *SOMEHOW* were being given Admin privleges.
		// if u.Admin {
		// 	screen.Permission = AdminPermissions
		// }
		screen.Redirect, _ = user.LogoutURL(ctx, "/"+req.FormValue("redirect"))
	}

	// Step 3: Hold onto the temporary info and go.
	retrievable.PlaceInMemcache(ctx, registerUUID, screen, 0)
	ServeTemplateWithParams(res, "register.html", screen)
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

	name := req.FormValue("Name")
	email := req.FormValue("Email")
	pswd := req.FormValue("Password")

	// Validate info

	ruuid, cErr := FromCookie(req, "UUID")
	if ErrorPage(res, "UUID Error: Cannot find UUID cookie.", cErr) {
		return
	}

	uNew := &User{}
	getErr := retrievable.GetFromMemcache(ctx, ruuid, uNew)
	if ErrorPage(res, "UUID Error: Cannot find UUID memcache value.", getErr) {
		return
	}

	uNew.Email = email
	uNew.Name = name

	uNew, createErr := CreateUserFromLogin(ctx, uNew.Email, pswd, uNew)
	if ErrorPage(res, "User Creation Error!", createErr) {
		return
	}

	// FUTURE: Advanced user permissions
	// permErr := SetPermission(ctx, uNew.ID, uNew.Permission)
	// if ErrorPage(res, "User Permission Error: Cannot create a permission value for user.", permErr) {
	// 	return
	// }

	_, sErr := NewSession(res, req, uNew.ID)
	if ErrorPage(res, "User Session Error: Cannot create session.", sErr) {
		return
	}
	DeleteCookie(res, "UUID")
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
	u, err := GetUserFromSession(res, req)
	fmt.Fprint(res, `<html><plaintext>`)
	fmt.Fprint(res, "Error?\n", err, "::\n")
	fmt.Fprint(res, "User?\n", u, "::\n")
}
