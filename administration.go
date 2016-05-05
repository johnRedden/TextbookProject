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
	"net/http"
)

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

func ADMIN_POST_ELEVATEUSER(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := req.FormValue("UEmail")

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

func ADMIN_GET_USERPERM(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := req.FormValue("UEmail")

	ctx := appengine.NewContext(req)
	pl, getErr := GetPermissionLevelFromDatastore(ctx, uEmail)
	if getErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"Internal Services Error: `+getErr.Error()+`","Code":500}`)
		return
	}
	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0,"Response":"`, pl, `"}`)
}

func ADMIN_POST_DELETEUSER(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := req.FormValue("UEmail")

	ctx := appengine.NewContext(req)

	DeleteMemchache(ctx, uEmail)
	RemovePermissionUserFromDatastore(ctx, uEmail)
	RemovePermissionLevelFromDatastore(ctx, uEmail)

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0}`)
}

func ADMIN_POST_ForceUserLogout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
		// User Must be at least Admin.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	uEmail := req.FormValue("UEmail")

	ctx := appengine.NewContext(req)

	memErr := DeleteMemchache(ctx, uEmail)
	if memErr != nil {
		fmt.Fprint(res, `{"Status":"Failure","Reason":"`+memErr.Error()+`","Code":500}`)
		return
	}

	fmt.Fprint(res, `{"Status":"Success","Reason":"","Code":0}`)
}
