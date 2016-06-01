// # API_Deleters
//
// Source Project: https://github.com/johnRedden/TextbookProject
//
// This package holds all api handlers with regards to structure that perform deletion operations.
// Permission requirement for these api calls: Admin
// For more information, please visit: https://github.com/johnRedden/TextbookProject/wiki
//
// This module shares a collective set of error codes described below:
//    Code: Message
//      0 - Success: All actions completed.
//    400 - Failure: Mandatory parameter missing; check reason for missing/invalid parameter.
//    418 - Failure: Authentication Error; check login status and permission level.
//    500 - Failure: Internal Services Error; check reason for more information.
//
package main

/*
API_Deleters.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
)

var (
	// Local Permission Variable: delete
	// This variable holds the minimum required permission level to use this module.
	api_Delete_Permission = AdminPermissions
)

// -------------------------------------------------------------------
// Deletion Data calls
// API calls for singular objects.
// Please read each section for expected input/output
/////////////

// Call: /api/delete/catalog
// Description:
// This call will delete a catalog and all child structures.
// ID should be a well-formed non-nil string of an existing catalog name.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	catalogID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil || catalogID == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Initalize Variables
	ctx := appengine.NewContext(req)
	keyCollection := make([]*datastore.Key, 0)
	fileCollection := make([]string, 0)

	// Add Parent(Catalog) to collection
	keyCollection = append(keyCollection, MakeCatalogKey(ctx, catalogID))

	for _, bk := range Get_Child_Key_From_Parent(ctx, catalogID, "Books") {
		keyCollection = append(keyCollection, bk)

		for _, chK := range Get_Child_Key_From_Parent(ctx, bk.IntID(), "Chapters") {
			keyCollection = append(keyCollection, chK)

			for _, sk := range Get_Child_Key_From_Parent(ctx, chK.IntID(), "Sections") {
				keyCollection = append(keyCollection, sk)

				for _, ok := range Get_Child_Key_From_Parent(ctx, sk.IntID(), "Objectives") {
					keyCollection = append(keyCollection, ok)
					fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ok.IntID()))...)

					for _, ek := range Get_Child_Key_From_Parent(ctx, ok.IntID(), "Exercises") {
						keyCollection = append(keyCollection, ek)
						fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ek.IntID()))...)
					}
				}
			}
		}
	}

	if err := datastore.DeleteMulti(ctx, keyCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Datastore Failure","code":500}`)
		return
	}

	if err := RemoveFilesFromGCS(ctx, fileCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Cloudstore Failure","code":500}`)
		return
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/delete/book
// Description:
// This call will delete a book and all child structures.
// ID should be a well-formatted integer of an existing book id.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	bookID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil || bookID == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Initalize Variables
	ctx := appengine.NewContext(req)
	keyCollection := make([]*datastore.Key, 0)
	fileCollection := make([]string, 0)

	// Add Parent(Book) to collection
	keyCollection = append(keyCollection, MakeBookKey(ctx, bookID))

	for _, chK := range Get_Child_Key_From_Parent(ctx, bookID, "Chapters") {
		keyCollection = append(keyCollection, chK)

		for _, sk := range Get_Child_Key_From_Parent(ctx, chK.IntID(), "Sections") {
			keyCollection = append(keyCollection, sk)

			for _, ok := range Get_Child_Key_From_Parent(ctx, sk.IntID(), "Objectives") {
				keyCollection = append(keyCollection, ok)
				fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ok.IntID()))...)

				for _, ek := range Get_Child_Key_From_Parent(ctx, ok.IntID(), "Exercises") {
					keyCollection = append(keyCollection, ek)
					fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ek.IntID()))...)
				}
			}
		}
	}

	if err := datastore.DeleteMulti(ctx, keyCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Datastore Failure","code":500}`)
		return
	}

	if err := RemoveFilesFromGCS(ctx, fileCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Cloudstore Failure","code":500}`)
		return
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/delete/chapter
// Description:
// This call will delete a chapter and all child structures.
// ID should be a well-formatted integer of an existing chapter id.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	chaptID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil || chaptID == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Initalize Variables
	ctx := appengine.NewContext(req)
	keyCollection := make([]*datastore.Key, 0)
	fileCollection := make([]string, 0)

	// Add Parent(Chapter) to collection
	keyCollection = append(keyCollection, MakeChapterKey(ctx, chaptID))

	for _, sk := range Get_Child_Key_From_Parent(ctx, chaptID, "Sections") {
		keyCollection = append(keyCollection, sk)

		for _, ok := range Get_Child_Key_From_Parent(ctx, sk.IntID(), "Objectives") {
			keyCollection = append(keyCollection, ok)
			fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ok.IntID()))...)

			for _, ek := range Get_Child_Key_From_Parent(ctx, ok.IntID(), "Exercises") {
				keyCollection = append(keyCollection, ek)
				fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ek.IntID()))...)
			}
		}
	}

	if err := datastore.DeleteMulti(ctx, keyCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Datastore Failure","code":500}`)
		return
	}

	if err := RemoveFilesFromGCS(ctx, fileCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Cloudstore Failure","code":500}`)
		return
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/delete/section
// Description:
// This call will delete a section and all child structures.
// ID should be a well-formatted integer of an existing section id.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	sectID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil || sectID == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Initalize Variables
	ctx := appengine.NewContext(req)
	keyCollection := make([]*datastore.Key, 0)
	fileCollection := make([]string, 0)

	// Add Parent(Section) to collection
	keyCollection = append(keyCollection, MakeSectionKey(ctx, sectID))

	for _, ok := range Get_Child_Key_From_Parent(ctx, sectID, "Objectives") {
		keyCollection = append(keyCollection, ok)
		fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ok.IntID()))...)

		for _, ek := range Get_Child_Key_From_Parent(ctx, ok.IntID(), "Exercises") {
			keyCollection = append(keyCollection, ek)
			fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ek.IntID()))...)
		}
	}

	if err := datastore.DeleteMulti(ctx, keyCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Datastore Failure","code":500}`)
		return
	}

	if err := RemoveFilesFromGCS(ctx, fileCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Cloudstore Failure","code":500}`)
		return
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/delete/objective
// Description:
// This call will delete an objective and all child structures.
// ID should be a well-formatted integer of an existing objective id.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	objID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil || objID == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Initalize Variables
	ctx := appengine.NewContext(req)
	keyCollection := make([]*datastore.Key, 0)
	fileCollection := make([]string, 0)

	// Add Parent(Objective) to collection
	keyCollection = append(keyCollection, MakeObjectiveKey(ctx, objID))
	fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(objID))...)

	for _, ek := range Get_Child_Key_From_Parent(ctx, objID, "Exercises") {
		keyCollection = append(keyCollection, ek)
		fileCollection = append(fileCollection, GetFilesFromGCS_WithPrefix(ctx, fmt.Sprint(ek.IntID()))...)
	}

	if err := datastore.DeleteMulti(ctx, keyCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Datastore Failure","code":500}`)
		return
	}

	if err := RemoveFilesFromGCS(ctx, fileCollection); err != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error: Cloudstore Failure","code":500}`)
		return
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/delete/exercise
// Description:
// This call will delete an Exercise and all child structures.
// ID should be a well-formatted integer of an existing Exercise id.
//
// Method: POST
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes: See Above.
func API_DeleteExercise(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	exerID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	// Remove this exercise from datastore.
	ctx := appengine.NewContext(req)
	remvErr := DeleteFromDatastore(ctx, exerID, &Exercise{})
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error","code":500}`)
	}

	// Clear this exercises images, if any.
	imagesToDelete := GetFilesFromGCS_WithPrefix(ctx, req.FormValue("ID"))
	RemoveFilesFromGCS(ctx, imagesToDelete)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}
