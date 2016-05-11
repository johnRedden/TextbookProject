package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
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
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetCatalogFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Catalog.html", itemToScreen)
}

// Call: /edit/Book/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getBookEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetBookFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Book.html", itemToScreen)
}

// Call: /edit/Chapter/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getChapterEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetChapterFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Chapter.html", itemToScreen)
}

// Call: /edit/Section/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSectionEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetSectionFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Section.html", itemToScreen)
}

// Call: /edit
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
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil || ObjectiveID == 0 {
		http.Redirect(res, req, "/select?status=invalid_id", http.StatusTemporaryRedirect)
	}
	ctx := appengine.NewContext(req)

	objKey := MakeObjectiveKey(ctx, int64(ObjectiveID))
	obj_temp := Objective{}
	objectiveGetErr := datastore.Get(ctx, objKey, &obj_temp)
	HandleError(res, objectiveGetErr)

	sect_temp := Section{}
	sectionGetErr := datastore.Get(ctx, MakeSectionKey(ctx, obj_temp.Parent), &sect_temp)
	HandleError(res, sectionGetErr)

	chap_temp := Chapter{}
	chapterGetErr := datastore.Get(ctx, MakeChapterKey(ctx, sect_temp.Parent), &chap_temp)
	HandleError(res, chapterGetErr)

	book_temp := Book{}
	bookGetErr := datastore.Get(ctx, MakeBookKey(ctx, chap_temp.Parent), &book_temp)
	HandleError(res, bookGetErr)

	ve := VIEW_Editor{}
	ve.ObjectiveID = objKey.IntID()
	ve.SectionID = obj_temp.Parent
	ve.ChapterID = sect_temp.Parent
	ve.BookID = chap_temp.Parent

	ve.ObjectiveTitle = obj_temp.Title
	ve.SectionTitle = sect_temp.Title
	ve.ChapterTitle = chap_temp.Title
	ve.BookTitle = book_temp.Title

	ve.ObjectiveVersion = obj_temp.Version
	ve.Content = obj_temp.Content
	ve.KeyTakeaways = obj_temp.KeyTakeaways
	ve.Author = obj_temp.Author

	ServeTemplateWithParams(res, req, "simpleEditor.html", ve)
}

// Call: /edit/Exercise/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getExerciseEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetExerciseFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_exercise.html", itemToScreen)
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
	readID := req.FormValue("ID")
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
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetExerciseFromDatastore(req, int64(i))
	HandleError(res, getErr)
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
	HandleError(res, convErr)
	objToScreen, err := GetObjectiveFromDatastore(req, objKey)
	HandleError(res, err)
	ServeTemplateWithParams(res, req, "preview.html", objToScreen)
}
