// Authentication
// This package handles all advanced Authentication/Session Management.
package main

/*
AUTH_authentication.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrPermissionUserMarshall, this error is thrown when there is an incorrect number of values to unmarshal into a Permission User
	ErrPermissionUserMarshall = errors.New("MarshallPermissionUser: Cannot Unmarshal String, Too Few Values")

	// ErrInvalidPermission, this error is thrown when a user fails a minimum permission level check.
	ErrInvalidPermission = errors.New("Permission Error: User Does Not Have Required Permission Level!")
)

//// --------------------------
// Permission User,
// A user with a permission value.
// Has a true name.
////

// Type: PermissionUser
// An appengine/user.User with Name and Permission
type PermissionUser struct {
	Name       string
	Permission int
	ID         string
	Email      string
}

// Method: ToString
// Converts a permission user into a marshalled string
func (u PermissionUser) ToString() string {
	return fmt.Sprintf("%s�%s�%d�%s", u.Name, u.Email, u.Permission, u.ID)
}

// Internal Function
// Description:
// Function takes a stringed permission user and unmarshalls the values back into their original struct.
//
// Returns:
//      user(PermissionUser) - Unpacked PermissionUser
//      failure?(error) - Any errors are stored here if exists.
func MarshallPermissionUser(p string) (PermissionUser, error) {
	data := strings.Split(p, "�")
	if len(data) < 4 {
		return PermissionUser{}, ErrPermissionUserMarshall
	}
	permLevel, convErr := strconv.Atoi(data[2])
	if convErr != nil {
		return PermissionUser{}, convErr
	}
	return PermissionUser{
		Name:       data[0],
		Email:      data[1],
		Permission: permLevel,
		ID:         data[3],
	}, nil
}

// Internal Function
// Description:
// Given an appengine/user.User, a name, and a permission level, will create a valid permission user.
//
// Returns:
//      user(PermissionUser) - Prepared PermissionUser
func MakePermissionUser(name string, permission int, u *user.User) PermissionUser {
	return PermissionUser{
		Name:       name,
		Permission: permission,
		Email:      u.Email,
		ID:         u.ID,
	}
}

// Internal Function
// Description:
// Given a session context, this will retrieve the current user's PermissionUser.
//
// Returns:
//      user(PermissionUser) - Prepared PermissionUser
//      failure?(error) - Any errors are stored here if exists.
func GetPermissionUserFromSession(ctx context.Context) (PermissionUser, error) {
	u := user.Current(ctx)
	if u != nil {
		if mVal, err := FromMemcache(ctx, u.Email); err == nil {
			if pVal, mErr := MarshallPermissionUser(mVal); mErr == nil {
				return pVal, nil
			} else {
				return PermissionUser{}, mErr
			}
		} else {
			return PermissionUser{}, err
		}
	}
	return PermissionUser{}, ErrNotLoggedIn
}

//// --------------------------
// Permission User, Datastore
// This collection of functions handle the insertion,
// retrieval, and deletion of PermissionUsers from
// datastore.
// All PermissionUsers exist on table Users
////

func PutPermissionUserToDatastore(ctx context.Context, keyname string, pu *PermissionUser) error {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	_, putErr := datastore.Put(ctx, userkey, pu)
	return putErr
}
func GetPermissionUserFromDatastore(ctx context.Context, keyname string) (PermissionUser, error) {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	pu := PermissionUser{}
	getErr := datastore.Get(ctx, userkey, &pu)
	return pu, getErr
}
func RemovePermissionUserFromDatastore(ctx context.Context, keyname string) error {
	userkey := datastore.NewKey(ctx, "Users", keyname, 0, nil)
	return datastore.Delete(ctx, userkey)
}

//// --------------------------
// Permission Levels
// This collection details the different levels
// of permissions a user can hold, verification
// that a user meets a minimum permission requirement
// and insertion, retrieval, and deletion of permission
// levels into datastore.
// Permission Levels in datastore are on table Permissions
////

const (
	// Permission Levels.
	// These are const integers. Please refer to them always by name, never number.
	ReadPermissions  = iota
	EditPermissions  = iota
	WritePermissions = iota
	AdminPermissions = iota
)

// Internal Function
// Description:
// Given a response, request, and minimum permission level.
// This function will return a boolean if the current user
// does or does not meet the requirement.
//
// Returns:
//      valid?(bool) - True/False if user meets requirement
//      failure?(error) - Any errors are stored here if exists.
func HasPermission(res http.ResponseWriter, req *http.Request, minimumRequiredPermission int) (bool, error) {
	if sessErr := MaintainSession(res, req); sessErr != nil { // Must have a session
		return false, sessErr
	} else {
		ctx := appengine.NewContext(req)
		if u, permissionErr := GetPermissionUserFromSession(ctx); permissionErr != nil { // Must have a valid permission user.
			return false, permissionErr
		} else {
			if u.Permission < minimumRequiredPermission { // That permission user must be at least the minimum.
				return false, ErrInvalidPermission
			}
		}
	}
	return true, nil
}

func PutPermissionLevelToDatastore(ctx context.Context, keyname string, permLevel int) error {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	toDatastore := &struct{ PL int }{permLevel}
	_, putErr := datastore.Put(ctx, permkey, toDatastore)
	return putErr
}
func GetPermissionLevelFromDatastore(ctx context.Context, keyname string) (int, error) {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	pl := struct{ PL int }{}
	getErr := datastore.Get(ctx, permkey, &pl)
	return pl.PL, getErr
}
func RemovePermissionLevelFromDatastore(ctx context.Context, keyname string) error {
	permkey := datastore.NewKey(ctx, "Permissions", keyname, 0, nil)
	return datastore.Delete(ctx, permkey)
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
		pu, getErr := GetPermissionUserFromDatastore(ctx, u.Email)
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
	ServeTemplateWithParams(res, req, "registerUser.html", u.Email)
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
	if pl, getErr := GetPermissionLevelFromDatastore(ctx, u.Email); getErr == nil {
		perms = pl // use the already determined permission level.
	} else {
		// Ensure that user is not a new administrator.
		if u.Admin {
			perms = AdminPermissions
		}
	}
	putLErr := PutPermissionLevelToDatastore(ctx, u.Email, perms)
	HandleError(res, putLErr)

	// Make user and add them to the datastore.
	permU := MakePermissionUser(uName, perms, u)
	putErr := PutPermissionUserToDatastore(ctx, u.Email, &permU)
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
