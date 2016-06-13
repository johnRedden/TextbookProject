package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/Esseh/retrievable"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strings"
	"time"
)

// Call: /admin
// Description:
// This is the root page of the Admin Console. Must be an Administrator to access.
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func ADMIN_AdministrationConsole(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}
	u, _ := GetUserFromSession(res, req)

	ServeTemplateWithParams(res, "adminConsole.html", u.Name)
}

// Call: /admin/changeUsrPerm
// Description:
// This call will change the permission level of a user.
//
// Method: POST
// Results: JSON
// Mandatory Options: UEmail, NewPermLevel
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    401 - Failure, Invalid Parameter
//    500 - Failure, Internal Services Error
func ADMIN_POST_ELEVATEUSER(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := strings.ToLower(req.FormValue("UEmail"))

	actualPermLevel := func(incomingPerm string) int {
		switch incomingPerm {
		default:
			return ReadPermissions
		case "Admin":
			return AdminPermissions
		case "Write":
			return WritePermissions
		case "Edit":
			return EditPermissions
		case "Read":
			return ReadPermissions
		}
	}(req.FormValue("NewPermLevel"))

	ctx := appengine.NewContext(req)

	uid, loginErr := GetUIDFromLogin(ctx, uEmail)
	if loginErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Email not found: `+loginErr.Error()+`","Code":500}`)
		return
	}

	u := &User{}
	getErr := retrievable.GetEntity(ctx, u, uid)
	if getErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User cannot be retrived: `+getErr.Error()+`","Code":500}`)
		return
	}

	u.Permission = actualPermLevel

	_, putErr := retrievable.PlaceEntity(ctx, uid, u)
	if putErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User cannot be stored: `+putErr.Error()+`","Code":500}`)
		return
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0}`)
}

// Call: /admin/getUsrPerm
// Description:
// This call will show the permission level of a user.
//
// Method: POST
// Results: JSON
// Mandatory Options: UEmail
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    500 - Failure, Internal Services Error
func ADMIN_GET_USERPERM(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := strings.ToLower(req.FormValue("UEmail"))

	ctx := appengine.NewContext(req)

	uid, loginErr := GetUIDFromLogin(ctx, uEmail)
	if loginErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Email not found: `+loginErr.Error()+`","Code":500}`)
		return
	}

	u := &User{}
	getErr := retrievable.GetEntity(ctx, u, uid)
	if getErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User cannot be retrived: `+getErr.Error()+`","Code":500}`)
		return
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,"Response":"`, u.Permission, `"}`)
}

// Call: /admin/deleteUsr
// Description:
// This call will delete a user.
//
// Method: POST
// Results: JSON
// Mandatory Options: UEmail
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    500 - Failure, Internal Services Error
func ADMIN_POST_DELETEUSER(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := strings.ToLower(req.FormValue("UEmail"))

	ctx := appengine.NewContext(req)
	delErr := DeleteUserAndLogin(ctx, uEmail)
	if delErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"`, delErr.Error(), `","Code":500}`)

	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0}`)
}

// Call: /admin/getEmailsofUSR
// Description:
// This call will return a list of Emails that are of username.
//
// Method: POST
// Results: JSON
// Mandatory Options: Usr
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    500 - Failure, Internal Services Error
func ADMIN_POST_RetriveUserEmails(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	usr := strings.ToLower(req.FormValue("Usr"))
	if usr == "" {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User Name cannot be Empty.","Code":406}`)
		return
	}

	q := datastore.NewQuery(UsersTable)
	collectedEmails := make([]string, 0)

	ctx := appengine.NewContext(req)
	for t := q.Run(ctx); ; { // standard query run.
		var tval User
		_, qErr := t.Next(&tval)

		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			fmt.Fprint(res, `{"Status":"Failure","Reason":"`, qErr.Error(), `","Code":500}`)
			return
		}
		if strings.Contains(strings.ToLower(tval.Name), usr) {
			collectedEmails = append(collectedEmails, tval.Email)
		}
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,`)
	fmt.Fprint(res, `"Results":[`)
	for i, v := range collectedEmails {
		if i != 0 {
			fmt.Fprint(res, `,`)
		}
		fmt.Fprint(res, `"`, v, `"`)
	}
	fmt.Fprint(res, `]}`)
}

// Call: /admin/createInviteUUID
// Description:
//
// Method: POST
// Results: JSON
// Mandatory Options: UName, Perm
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    401 - Failure, Invalid Parameter
//    500 - Failure, Internal Services Error
func ADMIN_POST_INVITEUUID(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	actualPermLevel := func(incomingPerm string) int {
		switch incomingPerm {
		default:
			return ReadPermissions
		case "Admin":
			return AdminPermissions
		case "Write":
			return WritePermissions
		case "Edit":
			return EditPermissions
		case "Read":
			return ReadPermissions
		}
	}(req.FormValue("Perm"))

	ctx := appengine.NewContext(req)

	uuid := NewUUID()
	newU := &User{}
	newU.Name = req.FormValue("UName")
	newU.Permission = actualPermLevel

	registerLimit := time.Hour * time.Duration(72)

	putErr := retrievable.PlaceInMemcache(ctx, uuid, newU, registerLimit)
	if putErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User cannot be stored: `+putErr.Error()+`","Code":500}`)
		return
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,"Results":"`, uuid, `"}`)
}
