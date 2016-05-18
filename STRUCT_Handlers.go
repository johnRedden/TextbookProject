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
	// "strings"
	// "html/template"
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetCatalogFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_Catalog.html", screenOutput)
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetBookFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_Book.html", screenOutput)
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetChapterFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_Chapter.html", screenOutput)
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetSectionFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_Section.html", screenOutput)
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetObjectiveFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_Objective.html", screenOutput)
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
		HandleErrorIntoPage(res, req, permErr, errorMessages[http.StatusUnauthorized])
		return
	}

	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetExerciseFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
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

	ServeTemplateWithParams(res, req, "editor_exercise.html", screenOutput)
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
	ServeTemplateWithParams(res, req, "simpleReader.html", readID)
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
	if HandleErrorIntoPage(res, req, parseErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	if i == 0 {
		HandleErrorIntoPage(res, req, errors.New("Invalid ID Given: Incoming parameter ID is a zero value."), "ID cannot be 0. Please ensure that the url is correct.")
		return
	}

	itemToScreen, getErr := GetExerciseFromDatastore(req, int64(i))
	if HandleErrorIntoPage(res, req, getErr, errorMessages[500]) {
		return
	}
	ServeTemplateWithParams(res, req, "reader_exercise.html", itemToScreen)
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
	objKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if HandleErrorIntoPage(res, req, convErr, "Invalid ID Given: Please ensure that the url is correct.") {
		return
	}

	objToScreen, err := GetObjectiveFromDatastore(req, objKey)
	if HandleErrorIntoPage(res, req, err, errorMessages[500]) {
		return
	}
	ServeTemplateWithParams(res, req, "preview.html", objToScreen)
}
