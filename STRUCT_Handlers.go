package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"net/http"
	"strconv"
)

// ------------------------------------
// Editors
/////

// Call: /edit/Catalog/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getCatalogEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetCatalogFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Catalog
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_Catalog.html", screenOutput)
}

// Call: /edit/Book/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getBookEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetBookFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Book
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_Book.html", screenOutput)
}

// Call: /edit/Chapter/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getChapterEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetChapterFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Chapter
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_Chapter.html", screenOutput)
}

// Call: /edit/Section/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSectionEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetSectionFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Section
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_Section.html", screenOutput)
}

// Call: /edit/objective/:ID
// Description:
// Our editor page for objectives given a valid objective id.
// Mandatory:ID must be a well-formatted integer of an existing objective id.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleObjectiveEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetObjectiveFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Objective
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_Objective.html", screenOutput)
}

// Call: /edit/Exercise/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getExerciseEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		ErrorPage(res, "Invalid Permission", permErr)
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetExerciseFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}

	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		Exercise
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		itemToScreen,
	}

	ServeTemplateWithParams(res, "editor_exercise.html", screenOutput)
}

// ------------------------------------
// Readers/Preview
/////

// Call: /read
// Description:
// Our reader page for objectives.
// Mandatory:ID has no requirements on this level. Sub levels will
// require that objective ID exists and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleObjectiveReader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	readID := params.ByName("ID")
	ServeTemplateWithParams(res, "simpleReader.html", readID)
}

// Call: /read/exercise/:ID
// Description:
// Our reader page for objectives.
// Mandatory:ID has no requirements on this level. Sub levels will
// require that objective ID exists and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleExerciseReader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	i, parseErr := strconv.Atoi(params.ByName("ID"))
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	if i == 0 {
		ErrorPage(res, "ID cannot be 0. Please ensure that the url is correct.", errors.New("Invalid ID Given: Incoming parameter ID is a zero value."))
		return
	}

	ctx := appengine.NewContext(req)
	itemToScreen, getErr := GetExerciseFromDatastore(ctx, int64(i))
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}
	ServeTemplateWithParams(res, "reader_exercise.html", itemToScreen)
}

// Call: /preview
// Description:
// Our simple page preview.
//
// Mandatory:ID must be an existing objective ID and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getObjectivePreview(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	objKey, parseErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if ErrorPage(res, "Invalid ID Given: Please ensure that the url is correct.", parseErr) {
		return
	}

	ctx := appengine.NewContext(req)
	objToScreen, getErr := GetObjectiveFromDatastore(ctx, objKey)
	if ErrorPage(res, "Internal Services Error", getErr) {
		return
	}
	ServeTemplateWithParams(res, "preview.html", objToScreen)
}
