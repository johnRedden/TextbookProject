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
	api_Make_Permission = WritePermissions
)

// -------------------------------------------------------------------
// Creation Data calls
// API calls for singular objects.
// Please read each call for expected input/output
////////

func API_MakeCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call for making a catalog, we would check for a signed in user here.
	// Mandatory Options: CatalogName
	// Optional Options: Company, Version, Description
	// Version should be a well formed stringed float
	// Codes:
	//      0 - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	catalogName := req.FormValue("CatalogName")
	if catalogName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name","code":400}`)
		return
	}

	catalogForDatastore, getErr := GetCatalogFromDatastore(req, catalogName)
	HandleError(res, getErr) // If this catalog already exists. We should go get that information to update it.
	catalogForDatastore.Name = catalogName

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
	HandleError(res, putErr)
	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Here is where we get a bit more complicated.
	// One stop shop for everything related to making a book
	// If you feed in an ID it will update that specific id
	// If no id is given, will make a new one.
	// Mandatory Options: CatalogName, BookName OR ID
	// Optional Options: Author, Version, Tags, Description
	// Codes:
	//      0 - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	bookID, _ := strconv.Atoi(req.FormValue("ID"))

	bookForDatastore, getErr := GetBookFromDatastore(req, int64(bookID))
	HandleError(res, getErr)

	if req.FormValue("CatalogName") != "" { // if your giving me a catalog, we're good
		bookForDatastore.Parent = req.FormValue("CatalogName")
	} else if bookID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name","code":400}`)
		return
	}

	if req.FormValue("BookName") != "" { // if your giving me a title, we're good
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

func API_MakeChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for chapter creation, same structure as above.
	// Mandatory Options: BookID, ChapterName OR ID
	// Optional Options: Version, Description
	// Codes:
	//      0 - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, api_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	chapterID, _ := strconv.Atoi(req.FormValue("ID"))

	chapterForDatastore, getErr := GetChapterFromDatastore(req, int64(chapterID))
	HandleError(res, getErr)

	bookID, numErr2 := strconv.Atoi(req.FormValue("BookID"))
	if numErr2 == nil { // if your giving me a catalog, we're good
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

func API_MakeSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Section creation, same structure as above.
	// Mandatory Options: ChapterID, SectionName OR ID
	// Optional Options:  Version, Description
	// Codes:
	//      0 - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data missing

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

	if req.FormValue("SectionName") != "" { // if your giving me a title, we're good
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

func API_MakeObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Objective creation, same structure as above.
	// Mandatory Options: ObjectiveName, SectionID OR ID
	// Optional Options: Version, Content, KeyTakeaways, Author
	// Codes:
	//      0 - Success, All completed
	//      418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	//      400 - Failure, Expected data missing
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
	if numErr2 == nil { // if your giving me a catalog, we're good
		objectiveForDatastore.Parent = int64(sectionID)
	} else if ObjectiveID == 0 { // new books must have a catalog
		fmt.Fprint(res, `{"result":"failure","reason":"Empty SectionID","code":400}`)
		return
	}

	if req.FormValue("ObjectiveName") != "" { // if your giving me a title, we're good
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
