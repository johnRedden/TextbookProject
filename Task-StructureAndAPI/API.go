package main

/*
API.go by Allen J. Mills
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

// -------------------------------------------------------------------
// Post Data calls
// API calls for singular objects.
// Please read each call for expected input/output
////////

func API_MakeCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call for making a catalog, we would check for a signed in user here.
	// Mandatory Options: CatalogName
	// Optional Options: Company, Version
	// Version should be a well formed stringed float
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

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
	// Optional Options: Author, Version, Tags
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	bookID, _ := strconv.Atoi(req.FormValue("ID"))

	bookForDatastore, getErr := GetBookFromDatastore(req, int64(bookID))
	HandleError(res, getErr)

	if req.FormValue("CatalogName") != "" { // if your giving me a catalog, we're good
		bookForDatastore.CatalogTitle = req.FormValue("CatalogName")
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

	_, putErr := PutBookIntoDatastore(req, bookForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for chapter creation, same structure as above.
	// Mandatory Options: BookID, ChapterName OR ID
	// Optional Options: Version
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

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

	_, putErr := PutChapterIntoDatastore(req, chapterForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Section creation, same structure as above.
	// Mandatory Options: ChapterID, SectionName OR ID
	// Optional Options:  Version
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

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

	_, putErr := PutSectionIntoDatastore(req, sectionForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Objective creation, same structure as above.
	// Mandatory Options: ObjectiveName, SectionID OR ID
	// Optional Options: Version, Content, KeyTakeaways, Author
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	ObjectiveID, _ := strconv.Atoi(req.FormValue("ID"))

	objectiveForDatastore, getErr := GetObjectiveFromDatastore(req, int64(ObjectiveID))
	HandleError(res, getErr)

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
		objectiveForDatastore.Content = req.FormValue("Content")
	}

	if req.FormValue("KeyTakeaways") != "" {
		objectiveForDatastore.KeyTakeaways = req.FormValue("KeyTakeaways")
	}

	if req.FormValue("Author") != "" {
		objectiveForDatastore.Author = req.FormValue("Author")
	}

	_, putErr := PutObjectiveIntoDatastore(req, objectiveForDatastore)
	HandleError(res, putErr)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// -------------------------------------------------------------------
// Query Data calls
// API calls for multiple objects.
// Will extend this later to detect singular calls
///////

func API_GetCatalogs(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Catalogs")
	cataloglist := make([]Catalog, 0)
	for t := q.Run(ctx); ; {
		var x Catalog
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.StringID()
		cataloglist = append(cataloglist, x)
	}

	ServeTemplateWithParams(res, req, "Catalogs.json", cataloglist)
}

func API_GetBooks(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Books")

	queryCatalogName := req.FormValue("Catalog")
	if queryCatalogName != "" {
		q = q.Filter("CatalogTitle =", queryCatalogName)
	}

	booklist := make([]Book, 0)
	for t := q.Run(ctx); ; {
		var x Book
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		booklist = append(booklist, x)
	}

	ServeTemplateWithParams(res, req, "Books.json", booklist)
}

func API_GetChapters(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Chapters")

	queryBookID := req.FormValue("BookID")
	if queryBookID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(queryBookID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	chapterList := make([]Chapter, 0)
	for t := q.Run(ctx); ; {
		var x Chapter
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		chapterList = append(chapterList, x)
	}

	ServeTemplateWithParams(res, req, "Chapters.json", chapterList)
}

func API_GetSections(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Sections")

	queryChapterID := req.FormValue("ChapterID")
	if queryChapterID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(queryChapterID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	sectionList := make([]Section, 0)
	for t := q.Run(ctx); ; {
		var x Section
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		sectionList = append(sectionList, x)
	}

	ServeTemplateWithParams(res, req, "Sections.json", sectionList)
}

func API_GetObjectives(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Objectives")

	querySectionID := req.FormValue("SectionID")
	if querySectionID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(querySectionID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	objectiveList := make([]Objective, 0)
	for t := q.Run(ctx); ; {
		var x Objective
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		objectiveList = append(objectiveList, x)
	}

	ServeTemplateWithParams(res, req, "Objectives.json", objectiveList)
}

// -------------------------------------------------------------------
// Query Data calls
// API calls for singular objects.
// Please read each section for expected input/output
/////////////

func API_GetTOC_HTML() {}

func API_GetObjectiveHTML(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving a <section> view on Objective
	// Mandatory Option: ID
	// Optional Options:
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty ID","code":400}`)
		return
	}

	objectiveToScreen, getErr := GetObjectiveFromDatastore(req, int64(ObjectiveID))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "ObjectiveHTML.html", objectiveToScreen)
}
