package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strings"
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
	ctx := appengine.NewContext(req)
	pu, _ := GetPermissionUserFromSession(ctx)

	ServeTemplateWithParams(res, req, "adminConsole.html", pu.Name)
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
			return -1
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
	if actualPermLevel < ReadPermissions { // if valid
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Invalid Permission Level","Code":401}`)
		return
	}

	ctx := appengine.NewContext(req)
	putErr := PutPermissionLevelToDatastore(ctx, uEmail, actualPermLevel)
	if putErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Internal Services Error: `+putErr.Error()+`","Code":500}`)
		return
	}

	if pu, getErr := GetPermissionUserFromDatastore(ctx, uEmail); getErr == nil {
		pu.Permission = actualPermLevel
		putErr2 := PutPermissionUserToDatastore(ctx, uEmail, &pu)
		if putErr2 != nil {
			fmt.Fprint(res, `{"Status":"Failure","Reason":"Internal Services Error: `+putErr2.Error()+`","Code":500}`)
			return
		}
	}

	if pustring, getErr := FromMemcache(ctx, uEmail); getErr == nil {
		// Try to update the current user.
		pu, _ := MarshallPermissionUser(pustring)
		pu.Permission = actualPermLevel
		ToMemcache(ctx, uEmail, pu.ToString(), StorageDuration)
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
	pl, getErr := GetPermissionLevelFromDatastore(ctx, uEmail)
	if getErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Internal Services Error: `+getErr.Error()+`","Code":500}`)
		return
	}
	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,"Response":"`, pl, `"}`)
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
func ADMIN_POST_DELETEUSER(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := strings.ToLower(req.FormValue("UEmail"))

	ctx := appengine.NewContext(req)

	DeleteMemcache(ctx, uEmail)
	RemovePermissionUserFromDatastore(ctx, uEmail)
	RemovePermissionLevelFromDatastore(ctx, uEmail)

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0}`)
}

// Call: /admin/forceUsrLogout
// Description:
// This call will force a user to logout.
//
// Method: POST
// Results: JSON
// Mandatory Options: UEmail
// Optional Options:
// Codes:
//      0 - Success, All actions completed
//    500 - Failure, Internal Services Error
func ADMIN_POST_ForceUserLogout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := strings.ToLower(req.FormValue("UEmail"))

	ctx := appengine.NewContext(req)

	memErr := DeleteMemcache(ctx, uEmail)
	if memErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"`+memErr.Error()+`","Code":500}`)
		return
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

	usr := req.FormValue("Usr")
	if usr == "" {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"User Name cannot be Empty.","Code":406}`)
		return
	}

	q := datastore.NewQuery("Users")
	q = q.Project("Name Email")
	collectedEmails := make([]string, 0)

	ctx := appengine.NewContext(req)
	for t := q.Run(ctx); ; { // standard query run.
		var tval struct{ Name, Email string }
		_, qErr := t.Next(&tval)

		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			break
		}
		if strings.HasPrefix(tval.Name, usr) {
			collectedEmails = append(collectedEmails, tval.Email)
		}
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,`)
	fmt.Fprint(res, `"results":[`)
	for i, v := range collectedEmails {
		if i != 0 {
			fmt.Fprint(res, `,`)
		}
		fmt.Fprint(res, `"`, v, `"`)
	}
	fmt.Fprint(res, `]}`)
}
