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
// Structures
// Internal use structures for ID handling

type catalogWithID struct {
	ID string
	Catalog
}
type bookWithID struct {
	ID int64
	Book
}
type chapterWithID struct {
	ID int64
	Chapter
}
type sectionWithID struct {
	ID int64
	Section
}
type objectiveWithID struct {
	ID int64
	Objective
}

// -------------------------------------------------------------------
// Post Data calls
// API calls for singular objects.
// Please read each call for expected input/output

func API_MakeCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call for making a catalog, we would check for a signed in user here.
	// Expects data from CatalogName
	// Also has data in from Company and Version
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
	// handle incoming data Company
	comp := req.FormValue("Company")
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}
	// Make our catalog
	catalogForDatastore := Catalog{}
	catalogForDatastore.Name = catalogName
	catalogForDatastore.Company = comp
	catalogForDatastore.Version = ver32
	// Get the datastore up and running!
	ctx := appengine.NewContext(req)

	ck := MakeCatalogKey(ctx, catalogName)
	_, errDatastore := datastore.Put(ctx, ck, &catalogForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Here is where we get a bit more complicated.
	// One stop shop for everything related to making a book
	// If you feed in an ID it will update that specific id
	// If no id is given, will make a new one.
	// Taking in mandatory options CatalogName and BookName
	// Optional options Author,Version
	// TODO: Add in tags functionality
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	bookID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		bookID = 0
	}

	catalogName := req.FormValue("CatalogName")
	if catalogName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name","code":400}`)
		return
	}

	bookName := req.FormValue("BookName")
	if bookName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Book Name","code":400}`)
		return
	}
	// handle incoming data Company
	auth := req.FormValue("Author")
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}

	// TODO: Add in something to allow for tags.
	// string parsing maybe?

	bookForDatastore := Book{}
	bookForDatastore.Title = bookName
	bookForDatastore.CatalogTitle = catalogName
	bookForDatastore.Author = auth
	bookForDatastore.Version = ver32

	ctx := appengine.NewContext(req)

	bk := MakeBookKey(ctx, int64(bookID))
	_, errDatastore := datastore.Put(ctx, bk, &bookForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for chapter creation, same structure as above.
	// Mandatory Options: BookID, ChapterName
	// Optional: ID, Version
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	chapterID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		chapterID = 0
	}

	bookID, numErr2 := strconv.Atoi(req.FormValue("BookID"))
	if numErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty BookID","code":400}`)
		return
	}

	chapterName := req.FormValue("ChapterName")
	if chapterName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty chapter Name","code":400}`)
		return
	}
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}

	chapterForDatastore := Chapter{}
	chapterForDatastore.Title = chapterName
	chapterForDatastore.Version = ver32
	chapterForDatastore.BookID = int64(bookID)

	ctx := appengine.NewContext(req)

	ck := MakeChapterKey(ctx, int64(chapterID))
	_, errDatastore := datastore.Put(ctx, ck, &chapterForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Section creation, same structure as above.
	// Mandatory Options: ChapterID, SectionName
	// Optional: ID, Version
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	sectionID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		sectionID = 0
	}

	chapterID, numErr2 := strconv.Atoi(req.FormValue("ChapterID"))
	if numErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty ChapterID","code":400}`)
		return
	}

	sectionName := req.FormValue("SectionName")
	if sectionName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty section Name","code":400}`)
		return
	}
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}

	sectionForDatastore := Section{}
	sectionForDatastore.Title = sectionName
	sectionForDatastore.Version = ver32
	sectionForDatastore.ChapterID = int64(chapterID)

	ctx := appengine.NewContext(req)

	ck := MakeSectionKey(ctx, int64(sectionID))
	_, errDatastore := datastore.Put(ctx, ck, &sectionForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

func API_MakeObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call  for Objective creation, same structure as above.
	// Mandatory Options: ObjectiveName, SectionID
	// Optional: ID, Version, Content, KeyTakeaways
	// Codes:
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		ObjectiveID = 0
	}

	sectionID, numErr2 := strconv.Atoi(req.FormValue("SectionID"))
	if numErr2 != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty SectionID","code":400}`)
		return
	}

	objectiveName := req.FormValue("ObjectiveName")
	if objectiveName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Objective Name","code":400}`)
		return
	}
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}

	objectiveForDatastore := Objective{}
	objectiveForDatastore.Title = objectiveName
	objectiveForDatastore.SectionID = int64(sectionID)
	objectiveForDatastore.Version = ver32
	objectiveForDatastore.Content = req.FormValue("Content")
	objectiveForDatastore.KeyTakeaways = req.FormValue("KeyTakeaways")

	ctx := appengine.NewContext(req)

	ck := MakeObjectiveKey(ctx, int64(ObjectiveID))
	_, errDatastore := datastore.Put(ctx, ck, &objectiveForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// -------------------------------------------------------------------
// Query Data calls
// API calls for multiple objects.
// Will extend this later to detect singular calls

func API_GetCatalogs(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Catalogs")
	cataloglist := make([]catalogWithID, 0)
	for t := q.Run(ctx); ; {
		var x Catalog
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		cataloglist = append(cataloglist, catalogWithID{k.StringID(), x})
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

	booklist := make([]bookWithID, 0)
	for t := q.Run(ctx); ; {
		var x Book
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		booklist = append(booklist, bookWithID{k.IntID(), x})
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
		q = q.Filter("BookID =", int64(i))
	}

	chapterList := make([]chapterWithID, 0)
	for t := q.Run(ctx); ; {
		var x Chapter
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		chapterList = append(chapterList, chapterWithID{k.IntID(), x})
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
		q = q.Filter("ChapterID =", int64(i))
	}

	sectionList := make([]sectionWithID, 0)
	for t := q.Run(ctx); ; {
		var x Section
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		sectionList = append(sectionList, sectionWithID{k.IntID(), x})
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
		q = q.Filter("SectionID =", int64(i))
	}

	objectiveList := make([]objectiveWithID, 0)
	for t := q.Run(ctx); ; {
		var x Objective
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		objectiveList = append(objectiveList, objectiveWithID{k.IntID(), x})
	}

	ServeTemplateWithParams(res, req, "Objectives.json", objectiveList)
}

// -------------------------------------------------------------------
// Query Data calls
// API calls for singular objects.
// Please read each section for expected input/outpu

func API_GetCatalog()   {}
func API_GetBook()      {}
func API_GetChapter()   {}
func API_GetSection()   {}
func API_GetObjective() {}
