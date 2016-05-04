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
	"html/template"
	"net/http"
	"strconv"
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
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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
	// 		0 - Success, All completed
	// 		418 - Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400 - Failure, Expected data missing
	ctx := appengine.NewContext(req)

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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

// -------------------------------------------------------------------
// Query/Collection Data calls
// API calls for multiple objects.
///////

func API_GetCatalogs(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Call for a json list of catalogs
	// Mandatory Options:
	// Optional Options:
	// Codes:
	// 		None, Data is either served or an http.Error is returned.
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
	// Call for a json list of books, can limit on originating catalog.
	// Mandatory Options:
	// Optional Options: Catalog
	// Codes:
	// 		None, Data is either served or an http.Error is returned.
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Books")

	queryCatalogName := req.FormValue("Catalog")
	if queryCatalogName != "" {
		q = q.Filter("Parent =", queryCatalogName)
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
	// Call for a json list of chapters, can limit on originating book.
	// Mandatory Options:
	// Optional Options: BookID
	// Codes:
	// 		None, Data is either served or an http.Error is returned.
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
	// Call for a json list of sections, can limit on originating chapter.
	// Mandatory Options:
	// Optional Options: ChapterID
	// Codes:
	// 		None, Data is either served or an http.Error is returned.
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
	// Call for a json list of objectives, can limit on originating section.
	// Mandatory Options:
	// Optional Options: SectionID
	// Codes:
	// 		None, Data is either served or an http.Error is returned.
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

func API_getTOC(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /toc?ID=<Book ID Number>
	// Generates a well formed xml document detailing all the identifing information of a book.
	// Mandatory Options: ID
	// Optional Options:
	// Codes:
	// 		XML<status> Failure - read <message> of error for more information
	// 		Success will return the well formed xml.

	/// - - - -
	// Initial Check, Ensure information is trivially good
	/////////

	BookID_In, numErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if numErr != nil || BookID_In == 0 {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
		// http.Redirect(res, req, "/?status=invalid_id", http.StatusTemporaryRedirect)
	}

	/// - - - -
	// Gather Book information, ensure that book exists.
	////////

	BookTitle, BookCatalog, BookID_Out := func(req *http.Request, id int64) (string, string, int64) { // get book data
		book_to_output, _ := GetBookFromDatastore(req, id)
		return book_to_output.Title, book_to_output.Parent, book_to_output.ID
	}(req, BookID_In)

	if BookID_In != BookID_Out {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Book Not Found!</message></error>`)
		// ServeTemplateWithParams(res, req, "printme.html", "ERROR! Incoming id not found!")
		return
	}

	/// - - - -
	// Prepare to make everything simple.
	//////

	ctx := appengine.NewContext(req)
	gatherKindGroup := Get_Name_ID_From_Parent // alias new function with old name.
	/// - - - -
	// Print header/Book information
	//////

	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><book>`)
	fmt.Fprintf(res, `<booktitle>%s</booktitle><bookid>%d</bookid><catalog>%s</catalog>`, BookTitle, BookID_Out, BookCatalog)

	/// - - - -
	// Gather & Print Sub information as available
	//////

	for _, singleChapter := range gatherKindGroup(ctx, BookID_Out, "Chapters") { // Sub-Layer Chapters
		fmt.Fprintf(res, `<chapter><chaptertitle>%s</chaptertitle><chapterid>%d</chapterid>`, singleChapter.Title, singleChapter.ID)

		for _, singleSection := range gatherKindGroup(ctx, singleChapter.ID, "Sections") {
			fmt.Fprintf(res, `<section><sectiontitle>%s</sectiontitle><sectionid>%d</sectionid>`, singleSection.Title, singleSection.ID)

			for _, singleObjective := range gatherKindGroup(ctx, singleSection.ID, "Objectives") {
				fmt.Fprintf(res, `<objective><objectivetitle>%s</objectivetitle><objectiveid>%d</objectiveid></objective>`, singleObjective.Title, singleObjective.ID)
			}
			fmt.Fprint(res, `</section>`) // Close this section
		}
		fmt.Fprint(res, `</chapter>`) // Close this chapter
	}

	/// - - - -
	// Close Book
	//////

	fmt.Fprint(res, `</book>`)
}

// -------------------------------------------------------------------
// Singular Data calls
// API calls for singular objects.
// Please read each section for expected input/output
/////////////

func API_GetCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving an xml view on Catalog
	// Mandatory Option: ID
	// Optional Options:

	CatalogID := req.FormValue("ID")
	if CatalogID == "" {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	Catalog_to_Output, geterr := GetCatalogFromDatastore(req, CatalogID)
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><catalog>`)
	fmt.Fprintf(res, `<title>%s</title>`, Catalog_to_Output.Name)
	fmt.Fprintf(res, `<version>%f</version>`, Catalog_to_Output.Version)
	fmt.Fprintf(res, `<parentid>%s</parentid>`, Catalog_to_Output.Company)
	fmt.Fprint(res, `<description>`+Catalog_to_Output.Description+`</description>`)
	fmt.Fprint(res, `</catalog>`)
}

func API_GetBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving an xml view on Book
	// Mandatory Option: ID
	// Optional Options:

	BookID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	Book_to_Output, geterr := GetBookFromDatastore(req, BookID)
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><book>`)
	fmt.Fprintf(res, `<title>%s</title>`, Book_to_Output.Title)
	fmt.Fprintf(res, `<author>%s</author>`, Book_to_Output.Author)
	fmt.Fprintf(res, `<version>%f</version>`, Book_to_Output.Version)
	fmt.Fprintf(res, `<catalog>%s</catalog>`, Book_to_Output.Parent)
	fmt.Fprintf(res, `<id>%d</id>`, Book_to_Output.ID)
	fmt.Fprintf(res, `<tags>%s</tags>`, Book_to_Output.Tags)
	fmt.Fprint(res, `<description>`+Book_to_Output.Description+`</description>`)
	fmt.Fprint(res, `</book>`)
}
func API_GetChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving an xml view on Chapter
	// Mandatory Option: ID
	// Optional Options:

	ChapterID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	Chapter_to_Output, geterr := GetChapterFromDatastore(req, ChapterID)
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><chapter>`)
	fmt.Fprintf(res, `<title>%s</title>`, Chapter_to_Output.Title)
	fmt.Fprintf(res, `<version>%f</version>`, Chapter_to_Output.Version)
	fmt.Fprintf(res, `<parentid>%d</parentid>`, Chapter_to_Output.Parent)
	fmt.Fprintf(res, `<id>%d</id>`, Chapter_to_Output.ID)
	fmt.Fprint(res, `<description>`+Chapter_to_Output.Description+`</description>`)
	fmt.Fprint(res, `</chapter>`)
}
func API_GetSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving an xml view on Section
	// Mandatory Option: ID
	// Optional Options:

	SectionID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	Section_to_Output, geterr := GetSectionFromDatastore(req, SectionID)
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><section>`)
	fmt.Fprintf(res, `<title>%s</title>`, Section_to_Output.Title)
	fmt.Fprintf(res, `<version>%f</version>`, Section_to_Output.Version)
	fmt.Fprintf(res, `<parentid>%d</parentid>`, Section_to_Output.Parent)
	fmt.Fprintf(res, `<id>%d</id>`, Section_to_Output.ID)
	fmt.Fprint(res, `<description>`+Section_to_Output.Description+`</description>`)
	fmt.Fprint(res, `</section>`)
}

func API_GetObjectiveHTML(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get call for reciving a <section> view on Objective
	// Mandatory Option: ID
	// Optional Options:

	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		fmt.Fprint(res, `<section><p>Request has failed: Invalid ID.</p></section>`)
		return
	}

	objectiveToScreen, getErr := GetObjectiveFromDatastore(req, int64(ObjectiveID))
	//HandleError(res, getErr)
	if getErr != nil {
		fmt.Fprint(res, `<section><p>Request has failed: No objective with given ID.</p></section>`)
		return
	}

	ServeTemplateWithParams(res, req, "ObjectiveHTML.html", objectiveToScreen)
}

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
	// 		0	- Success, All completed
	// 		418	- Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400	- Failure, Expected data Missing
	// 		500	- Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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
	// 		0	- Success, All completed
	// 		418	- Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400	- Failure, Expected data Missing
	// 		500	- Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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
	// 		0	- Success, All completed
	// 		418	- Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400	- Failure, Expected data Missing
	// 		500	- Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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
	// 		0	- Success, All completed
	// 		418	- Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400	- Failure, Expected data Missing
	// 		500	- Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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
	// 		0	- Success, All completed
	// 		418	- Failure, Authentication error, likely caused by a user not signed in or not allowed.
	// 		400	- Failure, Expected data Missing
	// 		500	- Failure, Internal Services Error. Thrown when removal from Datastore cannot be completed.

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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
