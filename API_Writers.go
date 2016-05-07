// # API_Writers
//
// Source Project: https://github.com/johnRedden/TextbookProject
//
// This package holds all api handlers with regards to structure that perform write operations.
// Permission requirement for these api calls: Writer
// For more information, please visit: https://github.com/johnRedden/TextbookProject/wiki
//
// This module shares a collective set of error codes described below:
//    Code: Message
//      0 - Success: All actions completed.
//    400 - Failure: Mandatory parameter missing; check reason for missing/invalid parameter.
//    418 - Failure: Authentication Error; check login status and permission level.
//
package main

/*
API_Writers.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"html/template"
	"net/http"
	"strconv"
)

var (
	// Local Permission Variable: make
	// This variable holds the minimum required permission level to use this module.
	api_Make_Permission = WritePermissions
)

// -------------------------------------------------------------------
// Creation Data calls
// API calls for singular objects.
// Please read each call for expected input/output
////////

// Call: /api/makeCatalog
// Description:
// This call is for the creation or update of a catalog. If given Mandatory:ID this call is in update mode. Mandatory:CatalogName must be a well-formed non-nil string. Mandatory:ID must be a well-formatted integer. Option:Version should be a well-formatted float value.
//
// Method: POST
// Results: JSON
// Mandatory Options: {CatalogName} OR {ID}
// Optional Options: Company, Version, Description
// Codes: See Above.
func API_MakeCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	catID, _ := strconv.Atoi(req.FormValue("ID"))

	catalogForDatastore, getErr := GetCatalogFromDatastore(req, int64(catID))
	if getErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Retrivial Error: `+getErr.Error()+`","code":500}`)
		return
	}
	// HandleError(res, getErr) // If this catalog already exists. We should go get that information to update it.

	if req.FormValue("CatalogName") != "" { // if you're giving me a title, we're good
		catalogForDatastore.Title = req.FormValue("CatalogName")
	} else if catID == 0 { // new catalogs must have a title
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name","code":400}`)
		return
	}

	if req.FormValue("Company") != "" {
		catalogForDatastore.Company = req.FormValue("Company")
	}

	if req.FormValue("Version") != "" {
		ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 64)
		if errFloat == nil {
			catalogForDatastore.Version = ver64
		}
	}

	if req.FormValue("Description") != "" {
		catalogForDatastore.Description = template.HTML(req.FormValue("Description"))
	}

	// Get the datastore up and running!
	_, putErr := PutCatalogIntoDatastore(req, catalogForDatastore)
	if putErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Placement Error: `+putErr.Error()+`","code":500}`)
		return
	}
	// HandleError(res, putErr)
	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/makeBook
// Description:
// This call will create or update book information. If Mandatory:ID is given, all parameters are set as update mode, otherwise Mandatory:CatalogName, Mandatory:BookName must be given. Option:Version should be a well-formatted float.
//
// Method: POST
// Results: JSON
// Mandatory Options: {CatalogID, BookName} OR {ID}
// Optional Options: Author, Version, Tags, Description
// Codes: See Above.
func API_MakeBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	bookID, _ := strconv.Atoi(req.FormValue("ID"))

	bookForDatastore, getErr := GetBookFromDatastore(req, int64(bookID))
	HandleError(res, getErr)

	if catKey, parseErr := strconv.ParseInt(req.FormValue("CatalogID"), 10, 64); parseErr == nil && catKey != int64(0) { // if you're giving me a catalog, we're good
		bookForDatastore.Parent = catKey
	} else if bookID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name","code":400}`)
		return
	}

	if req.FormValue("BookName") != "" { // if you're giving me a title, we're good
		bookForDatastore.Title = req.FormValue("BookName")
	} else if bookID == 0 { // new books must have a title
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Book Name","code":400}`)
		return
	}

	if req.FormValue("Author") != "" { // if updating author
		bookForDatastore.Author = req.FormValue("Author")
	}

	if req.FormValue("Version") != "" { // if updating version
		ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 64)
		if errFloat == nil {
			bookForDatastore.Version = ver64
		}
	}

	if req.FormValue("Tags") != "" { // if updating tags
		bookForDatastore.Tags = req.FormValue("Tags")
	}

	if req.FormValue("Description") != "" {
		bookForDatastore.Description = template.HTML(req.FormValue("Description"))
	}

	_, putErr := PutBookIntoDatastore(req, bookForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/makeChapter
// Description:
// This call will create or update chapter information. If Mandatory:ID is given, all parameters are set as update mode, otherwise Mandatory:BookID, Mandatory:ChapterName must be given. Option:Version should be a well-formatted float. Mandatory:BookID should be a well-formatted integer.
//
// Method: POST
// Results: JSON
// Mandatory Options: {BookID, ChapterName} OR {ID}
// Optional Options: Version, Description
// Codes: See Above.
func API_MakeChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	chapterID, _ := strconv.Atoi(req.FormValue("ID"))

	chapterForDatastore, getErr := GetChapterFromDatastore(req, int64(chapterID))
	HandleError(res, getErr)

	bookID, numErr2 := strconv.Atoi(req.FormValue("BookID"))
	if numErr2 == nil { // if you're giving me a catalog, we're good
		chapterForDatastore.Parent = int64(bookID)
	} else if chapterID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty BookID","code":400}`)
		return
	}

	if req.FormValue("ChapterName") != "" { // if your giving me a title, we're good
		chapterForDatastore.Title = req.FormValue("ChapterName")
	} else if chapterID == 0 { // new books must have a title
		fmt.Fprint(res, `{"result":"failure","reason":"Empty chapter Name","code":400}`)
		return
	}

	if req.FormValue("Version") != "" { // if updating version
		ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 64)
		if errFloat == nil {
			chapterForDatastore.Version = ver64
		}
	}

	if req.FormValue("Description") != "" {
		chapterForDatastore.Description = template.HTML(req.FormValue("Description"))
	}

	_, putErr := PutChapterIntoDatastore(req, chapterForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/makeSection
// Description:
// This call will create or update section information. If Mandatory:ID is given, all parameters are set as update mode, otherwise Mandatory:ChapterID, Mandatory:SectionName must be given. Option:Version should be a well-formatted float. Mandatory:ChapterID should be a well-formatted integer.
//
// Method: POST
// Results: JSON
// Mandatory Options: {ChapterID, SectionName} OR {ID}
// Optional Options: Version, Description
// Codes: See Above.
func API_MakeSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	sectionID, _ := strconv.Atoi(req.FormValue("ID"))

	sectionForDatastore, getErr := GetSectionFromDatastore(req, int64(sectionID))
	HandleError(res, getErr)

	chapterID, numErr2 := strconv.Atoi(req.FormValue("ChapterID"))
	if numErr2 == nil { // if your giving me a catalog, we're good
		sectionForDatastore.Parent = int64(chapterID)
	} else if sectionID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty ChapterID","code":400}`)
		return
	}

	if req.FormValue("SectionName") != "" { // if you're giving me a title, we're good
		sectionForDatastore.Title = req.FormValue("SectionName")
	} else if sectionID == 0 { // new books must have a title
		fmt.Fprint(res, `{"result":"failure","reason":"Empty section Name","code":400}`)
		return
	}

	if req.FormValue("Version") != "" { // if updating version
		ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 64)
		if errFloat == nil {
			sectionForDatastore.Version = ver64
		}
	}

	if req.FormValue("Description") != "" {
		sectionForDatastore.Description = template.HTML(req.FormValue("Description"))
	}

	_, putErr := PutSectionIntoDatastore(req, sectionForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// Call: /api/makeObjective
// Description:
// This call will create or update objective information. If Mandatory:ID is given, all parameters are set as update mode, otherwise Mandatory:SectionID, Mandatory:ObjectiveName must be given. Option:Version should be a well-formatted float. Mandatory:SectionID should be a well-formatted integer.
//
// Method: POST
// Results: JSON
// Mandatory Options: {SectionID, ObjectiveName} OR {ID}
// Optional Options: Version, Content, KeyTakeaways, Author
// Codes: See Above.
func API_MakeObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	ObjectiveID, _ := strconv.Atoi(req.FormValue("ID"))

	objectiveForDatastore, getErr := GetObjectiveFromDatastore(req, int64(ObjectiveID))
	HandleErrorWithLog(res, getErr, "api/makeObjective Error: (GET) ", ctx)

	sectionID, numErr2 := strconv.Atoi(req.FormValue("SectionID"))
	if numErr2 == nil { // if you're giving me a section, we're good
		objectiveForDatastore.Parent = int64(sectionID)
	} else if ObjectiveID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty SectionID","code":400}`)
		return
	}

	if req.FormValue("ObjectiveName") != "" { // if you're giving me a title, we're good
		objectiveForDatastore.Title = req.FormValue("ObjectiveName")
	} else if ObjectiveID == 0 { // new books must have a title
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Objective Name","code":400}`)
		return
	}

	if req.FormValue("Version") != "" { // if updating version
		ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 64)
		if errFloat == nil {
			objectiveForDatastore.Version = ver64
		}
	}

	if req.FormValue("Content") != "" {
		objectiveForDatastore.Content = template.HTML(req.FormValue("Content"))
	}

	if req.FormValue("KeyTakeaways") != "" {
		objectiveForDatastore.KeyTakeaways = template.HTML(req.FormValue("KeyTakeaways"))
	}

	if req.FormValue("Author") != "" {
		objectiveForDatastore.Author = req.FormValue("Author")
	}

	_, putErr := PutObjectiveIntoDatastore(req, objectiveForDatastore)
	HandleErrorWithLog(res, putErr, "api/makeObjective Error: (PUT) ", ctx)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}
