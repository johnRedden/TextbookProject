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

// Call: /api/deleteCatalog
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

	catalogKey, _ := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if catalogKey == 0 {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	keyCollection := make([]*datastore.Key, 0) // From here on out, we will use the batch version of delete to ensure all or none of the objects are deleted.

	ctx := appengine.NewContext(req)
	keyCollection = append(keyCollection, MakeCatalogKey(ctx, catalogKey))

	for _, bookNameKey := range Get_Name_ID_From_Parent(ctx, catalogKey, "Books") {
		keyCollection = append(keyCollection, MakeBookKey(ctx, bookNameKey.ID))
		for _, chaptNameKey := range Get_Name_ID_From_Parent(ctx, bookNameKey.ID, "Chapters") {
			keyCollection = append(keyCollection, MakeChapterKey(ctx, chaptNameKey.ID))
			for _, sectNameKey := range Get_Name_ID_From_Parent(ctx, chaptNameKey.ID, "Sections") {
				keyCollection = append(keyCollection, MakeSectionKey(ctx, sectNameKey.ID))
				for _, objeNameKey := range Get_Name_ID_From_Parent(ctx, sectNameKey.ID, "Objectives") {
					keyCollection = append(keyCollection, MakeObjectiveKey(ctx, objeNameKey.ID))
				}
			}
		}
	}

	remvErr := datastore.DeleteMulti(ctx, keyCollection)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error:`+remvErr.Error()+`","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/deleteBook
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

	bookKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	keyCollection := make([]*datastore.Key, 0) // From here on out, we will use the batch version of delete to ensure all or none of the objects are deleted.

	ctx := appengine.NewContext(req)
	keyCollection = append(keyCollection, MakeBookKey(ctx, bookKey))
	for _, chaptNameKey := range Get_Name_ID_From_Parent(ctx, bookKey, "Chapters") {
		keyCollection = append(keyCollection, MakeChapterKey(ctx, chaptNameKey.ID))
		for _, sectNameKey := range Get_Name_ID_From_Parent(ctx, chaptNameKey.ID, "Sections") {
			keyCollection = append(keyCollection, MakeSectionKey(ctx, sectNameKey.ID))
			for _, objeNameKey := range Get_Name_ID_From_Parent(ctx, sectNameKey.ID, "Objectives") {
				keyCollection = append(keyCollection, MakeObjectiveKey(ctx, objeNameKey.ID))
			}
		}
	}

	remvErr := datastore.DeleteMulti(ctx, keyCollection)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error:`+remvErr.Error()+`","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/deleteChapter
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

	chaptKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	keyCollection := make([]*datastore.Key, 0) // From here on out, we will use the batch version of delete to ensure all or none of the objects are deleted.

	ctx := appengine.NewContext(req)
	keyCollection = append(keyCollection, MakeChapterKey(ctx, chaptKey))

	for _, sectNameKey := range Get_Name_ID_From_Parent(ctx, chaptKey, "Sections") {
		keyCollection = append(keyCollection, MakeSectionKey(ctx, sectNameKey.ID))

		for _, objeNameKey := range Get_Name_ID_From_Parent(ctx, sectNameKey.ID, "Objectives") {
			keyCollection = append(keyCollection, MakeObjectiveKey(ctx, objeNameKey.ID))
		}
	}

	remvErr := datastore.DeleteMulti(ctx, keyCollection)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error:`+remvErr.Error()+`","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/deleteSection
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

	sectKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	keyCollection := make([]*datastore.Key, 0) // From here on out, we will use the batch version of delete to ensure all or none of the objects are deleted.

	ctx := appengine.NewContext(req)
	keyCollection = append(keyCollection, MakeSectionKey(ctx, sectKey))

	for _, objeNameKey := range Get_Name_ID_From_Parent(ctx, sectKey, "Objectives") {
		keyCollection = append(keyCollection, MakeObjectiveKey(ctx, objeNameKey.ID))
	}

	remvErr := datastore.DeleteMulti(ctx, keyCollection)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error:`+remvErr.Error()+`","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/deleteObjective
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

	objKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid ID","code":400}`)
		return
	}

	remvErr := RemoveObjectiveFromDatastore(req, objKey)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/deleteExercise
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

	remvErr := RemoveExerciseFromDatastore(req, exerID)
	if remvErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error","code":500}`)
	}

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}
