package main

/*
filename.go by Allen J. Mills
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
	"strings"
)

//// --------------------------
// Permission User,
// A user with a permission value.
// Has a true name.
////

var (
	ErrPermissionUserMarshall = errors.New("MarshallPermissionUser: Cannot Marshall String, Too Few Values")
)

type PermissionUser struct {
	Name       string
	Permission string
	ID         string
	Email      string
}

func (u PermissionUser) ToString() string {
	return fmt.Sprintf("%s�%s�%s�%s", u.Name, u.Email, u.Permission, u.ID)
}

func MarshallPermissionUser(p string) (PermissionUser, error) {
	data := strings.Split(p, "�")
	if len(data) < 4 {
		return PermissionUser{}, ErrPermissionUserMarshall
	}
	return PermissionUser{
		Name:       data[0],
		Email:      data[1],
		Permission: data[2],
		ID:         data[3],
	}, nil
}

func MakePermissionUser(name, permission string, u *user.User) PermissionUser {
	return PermissionUser{
		Name:       name,
		Permission: permission,
		Email:      u.Email,
		ID:         u.ID,
	}
}

func GetPermissionUserFromSession(ctx context.Context) (PermissionUser, error) {
	u := user.Current(ctx)
	if u != nil {
		if mVal, err := FromMemcache(ctx, MemcacheKey(u)); err == nil {
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
// Permisison User, Datastore
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
// Login Process,
////

func AUTH_Login_GET(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// User has requested a login procedure.
	// Attempt to gather user info.
	redirectBack := ""
	if req.FormValue("redirect") != "" {
		redirectBack = "?redirect=" + req.FormValue("redirect")
	}
	ctx := appengine.NewContext(req)
	if sessErr := MaintainSession(res, req); sessErr == ErrNotLoggedIn || req.FormValue("changeuser") == "yes" {
		// User is not logged in.
		// Force them to the google login page before coming back here.

		// TODO: Somewhere in this process, double logins are occuring. Find this please!
		// We've traced it down to somewhere here.

		http.Redirect(res, req, GetLoginURL(ctx, "/login"+redirectBack), http.StatusTemporaryRedirect)
		return
	} else if sessErr == ErrTimedOut {
		// User has an oauth key.
		// Likely returned from ouath.
		u := user.Current(ctx)
		pu, getErr := GetPermissionUserFromDatastore(ctx, u.Email)
		if getErr != nil {
			// They do not have a registered permission user.
			// Kick them over to register.
			http.Redirect(res, req, "/register"+redirectBack, http.StatusTemporaryRedirect)
			return
		}
		// we now have their user information.
		sessErr := CreateSession(res, req, func(u *user.User) string { return pu.ToString() })
		if sessErr != nil {
			http.Error(res, sessErr.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Session is live.
	redirectTo := req.FormValue("redirect")
	if redirectTo == "" {
		redirectTo = "/"
	}
	http.Redirect(res, req, redirectTo, http.StatusSeeOther)
}

func AUTH_Register_GET(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)
	if u == nil {
		// They are not logged in!
		// We'll kick them over to google for ouath.
		http.Redirect(res, req, GetLoginURL(ctx, "/register?redirect="+req.FormValue("redirect")), http.StatusTemporaryRedirect)
		return
	}
	// TODO: Create an actual login page and serve that.
	page := `<!DOCTYPE html><html><body>
<form method="POST">
    <p>Name: <input name="Name" autofocus></input></p>
    <input type="submit">
</form>
</body></html>`
	fmt.Fprint(res, page)
}

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
	// TODO: Replace permission with something significant.
	uName := req.FormValue("Name")
	perms := "none"
	if u.Admin {
		perms = "admin"
	}

	// Make user and add them to the datastore.
	permU := MakePermissionUser(uName, perms, u)
	putErr := PutPermissionUserToDatastore(ctx, u.Email, &permU)
	if putErr != nil {
		http.Error(res, putErr.Error(), http.StatusInternalServerError)
		return
	}

	// Now we make that session
	memValue := func(u *user.User) string {
		return permU.ToString()
	}
	sessErr := CreateSession(res, req, memValue)
	if sessErr != nil {
		http.Error(res, sessErr.Error(), http.StatusInternalServerError)
		return
	}

	redirectTo := req.FormValue("redirect")
	if redirectTo == "" {
		redirectTo = "/"
	}
	http.Redirect(res, req, redirectTo, http.StatusSeeOther)
}

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