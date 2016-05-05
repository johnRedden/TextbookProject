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
	"html/template"
	"net/http"
	"strconv"
)

var (
	// Local Variables to set module permission levels.
	api_Delete_Permission = AdminPermissions
)

// -------------------------------------------------------------------
// Deletion Data calls
// API calls for singular objects.
// Please read each section for expected input/output
/////////////

func API_DeleteCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for deletion of a Catalog
	//  Will also delete all data pointing to it.
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	//      0   - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data Missing
	//      500 - Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, api_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	catalogKey := req.FormValue("ID")
	if catalogKey == "" {
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

	fmt.Fprint(res, `{"result":"success","reason":","code":0}`)
}
func API_DeleteBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for deletion of a Book
	//  Will also delete all data pointing to it.
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	//      0   - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data Missing
	//      500 - Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

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

	fmt.Fprint(res, `{"result":"success","reason":","code":0}`)
}
func API_DeleteChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for deletion of a Chapter
	//  Will also delete all data pointing to it.
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	//      0   - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data Missing
	//      500 - Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

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

	fmt.Fprint(res, `{"result":"success","reason":","code":0}`)
}
func API_DeleteSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for deletion of a Section
	//  Will also delete all data pointing to it.
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	//      0   - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data Missing
	//      500 - Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

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

	fmt.Fprint(res, `{"result":"success","reason":","code":0}`)
}
func API_DeleteObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for deletion of an objective
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	//      0   - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data Missing
	//      500 - Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

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

	fmt.Fprint(res, `{"result":"success","reason":","code":0}`)
}
