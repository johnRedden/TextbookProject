// # API_Readers
//
// Source Project: https://github.com/johnRedden/TextbookProject
//
// This package holds all api handlers with regards to structure that perform read operations.
// No requirement currently exists in respect to permissions.
// For more information, please visit: https://github.com/johnRedden/TextbookProject/wiki
//
package main

/*
API_Readers.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
)

// -------------------------------------------------------------------
// Query/Collection Data calls
// API calls for multiple objects.
///////

// Call: /api/catalogs.json
// Description:
// This call will return a complete list of catalogs. There are no options to limit results.
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options:
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetCatalogs(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Catalogs")
	q = q.Order("Title")
	cataloglist := make([]Catalog, 0)
	for t := q.Run(ctx); ; {
		var x Catalog
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		cataloglist = append(cataloglist, x)
	}
	ServeTemplateWithParams(res, "Catalogs.json", cataloglist)
}

// Call: /api/books.json
// Description:
// This call will return a list of currently available books. Results may be limited by parent catalog title. Option:Catalog must be a well-formed non-nil string.
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options: Catalog
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetBooks(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Books")
	q = q.Order("Title")

	queryCatID := req.FormValue("Catalog")
	if queryCatID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(queryCatID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
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
	ServeTemplateWithParams(res, "Books.json", booklist)
}

// Call: /api/chapters.json
// Description:
// This call will return a list of currently available chapters. May limit results based on parent book id. Option:BookID must be a well-formatted integer number.
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options: BookID
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetChapters(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Chapters")

	queryBookID := req.FormValue("BookID")
	if queryBookID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(queryBookID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	q = q.Order("Order")
	q = q.Order("Title")

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

	ServeTemplateWithParams(res, "Chapters.json", chapterList)
}

// Call: /api/sections.json
// Description:
// This call will return a list of currently available sections. May limit results based on parent chapter id. Option:ChapterID must be a well-formatted integer number.
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options: ChapterID
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetSections(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Sections")

	queryChapterID := req.FormValue("ChapterID")
	if queryChapterID != "" { // Ensure that a ChapterID was indeed sent.
		i, numErr := strconv.Atoi(queryChapterID) // does that ChapterID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	q = q.Order("Order").Order("Title")

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

	ServeTemplateWithParams(res, "Sections.json", sectionList)
}

// Call: /api/sections.json
// Description:
// This call will return a list of currently available objectives. May limit results based on parent section id. Option:SectionID must be a well-formatted integer number.
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options: SectionID
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetObjectives(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Objectives")

	querySectionID := req.FormValue("SectionID")
	if querySectionID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(querySectionID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	q = q.Order("Order").Order("Title")

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

	ServeTemplateWithParams(res, "Objectives.json", objectiveList)
}

// Call: /api/exercises.json
// Description:
// This call will return a complete list of exercises.
// Limit results by parent objective by ObjectiveID
// Limit results by Instruction kind by IKind
//
// Method: GET
// Results: JSON
// Mandatory Options:
// Optional Options: ObjectiveID, IKind
// Codes:
//      None, Data is either served or an http.Error is returned.
func API_GetExercises(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Exercises")

	if req.FormValue("ObjectiveID") != "" {
		i, numErr := strconv.Atoi(req.FormValue("ObjectiveID"))
		HandleError(res, numErr)
		q = q.Filter("Parent =", int64(i))
	}

	if req.FormValue("IKind") != "" {
		q = q.Filter("Instruction =", req.FormValue("IKind"))
	}

	q = q.Order("Order").Order("Instruction")

	exerciselist := make([]Exercise, 0)
	for t := q.Run(ctx); ; {
		var x Exercise
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		x.ID = k.IntID()
		exerciselist = append(exerciselist, x)
	}
	ServeTemplateWithParams(res, "Exercises.json", exerciselist)
}

// Call: /api/toc.xml
// Description:
// This call will return an xml formatted view of an entire book by ID. ID must be a well-formatted integer id of an existing book.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_getTOC(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

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
	ctx := appengine.NewContext(req)

	BookTitle, BookCatalog, BookID_Out := func(req *http.Request, id int64) (string, int64, int64) { // get book data
		book_to_output, _ := GetBookFromDatastore(ctx, id)
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

	gatherKindGroup := Get_Name_ID_From_Parent // alias new function with old name.
	/// - - - -
	// Print header/Book information
	//////

	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><book>`)
	fmt.Fprintf(res, `<booktitle>%s</booktitle><bookid>%d</bookid><catalog>%d</catalog>`, BookTitle, BookID_Out, BookCatalog)

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

type JsonOptions struct {
	Status, Reason string
	Code           int
	Results        interface{}
}

func ServeJsonOfStruct(res http.ResponseWriter, opt JsonOptions, param interface{}) error {
	opt.Results = param
	output, merr := json.Marshal(opt)
	fmt.Fprint(res, string(output))
	return merr
}

// Call: /api/parent/:KIND/:ID
// Description:
//
// Method: GET
// Results: JSON
// Mandatory Options: ID
// Optional Options:
// Codes:
func API_GetParent(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	screen := struct {
		BookTitle, ChapterTitle                    string // titles
		SectionTitle, ObjectiveTitle               string //
		BookID, ChapterID                          int64  //
		SectionID, ObjectiveID                     int64  // ids
		ChapterOrder, SectionOrder, ObjectiveOrder int    // orders
	}{}

	id, parseErr := strconv.ParseInt(params.ByName("ID"), 10, 64)
	if parseErr != nil {
		ServeJsonOfStruct(res, JsonOptions{
			Code:   http.StatusNotAcceptable,
			Status: "Failure",
			Reason: parseErr.Error(),
		}, nil)
		return
	} else if id == 0 {
		ServeJsonOfStruct(res, JsonOptions{
			Code:   http.StatusNotAcceptable,
			Status: "Failure",
			Reason: "ID cannot be zero.",
		}, nil)
		return
	}

	switch params.ByName("KIND") {
	default:
		ServeJsonOfStruct(res, JsonOptions{
			Code:   http.StatusNotAcceptable,
			Status: "Failure",
			Reason: "Invalid Kind given",
		}, nil)
		return
	case "Objective":
		screen.ObjectiveID = id
	case "Section":
		screen.SectionID = id
	case "Chapter":
		screen.ChapterID = id
	case "Book":
		screen.BookID = id
	}

	ctx := appengine.NewContext(req)

	// DO OBJECTIVE
	if screen.ObjectiveID != 0 {
		ob, err := GetObjectiveFromDatastore(ctx, screen.ObjectiveID)
		if err != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Code:   500,
				Status: "Failure",
				Reason: err.Error(),
			}, nil)
			return
		}
		screen.ObjectiveTitle = ob.Title
		screen.ObjectiveOrder = ob.Order
		screen.SectionID = ob.Parent
	}

	if screen.SectionID != 0 {
		sc, err := GetSectionFromDatastore(ctx, screen.SectionID)
		if err != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Code:   500,
				Status: "Failure",
				Reason: err.Error(),
			}, nil)
			return
		}
		screen.SectionTitle = sc.Title
		screen.SectionOrder = sc.Order
		screen.ChapterID = sc.Parent
	}

	if screen.ChapterID != 0 {
		ch, err := GetChapterFromDatastore(ctx, screen.ChapterID)
		if err != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Code:   500,
				Status: "Failure",
				Reason: err.Error(),
			}, nil)
			return
		}
		screen.ChapterTitle = ch.Title
		screen.ChapterOrder = ch.Order
		screen.BookID = ch.Parent
	}

	if screen.BookID != 0 {
		bk, err := GetBookFromDatastore(ctx, screen.BookID)
		if err != nil {
			ServeJsonOfStruct(res, JsonOptions{
				Code:   500,
				Status: "Failure",
				Reason: err.Error(),
			}, nil)
			return
		}
		screen.BookTitle = bk.Title
	}

	ServeJsonOfStruct(res, JsonOptions{
		Code:   0,
		Status: "Success",
		Reason: "",
	}, screen)
	res.Header().Set("Content-Type", "text/json")
}

// -------------------------------------------------------------------
// Singular Data calls
// API calls for singular objects.
// Please read each section for expected input/output
/////////////

// Call: /api/catalog.xml
// Description:
// Call will return an xml view on a singular catalog. ID is a well-formed string of an existing catalog.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_GetCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	CatalogID, _ := strconv.Atoi(req.FormValue("ID"))
	if CatalogID == 0 {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	ctx := appengine.NewContext(req)
	Catalog_to_Output, geterr := GetCatalogFromDatastore(ctx, int64(CatalogID))
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><catalog>`)
	fmt.Fprintf(res, `<title>%s</title>`, Catalog_to_Output.Title)
	fmt.Fprintf(res, `<version>%f</version>`, Catalog_to_Output.Version)
	fmt.Fprintf(res, `<parentid>%s</parentid>`, Catalog_to_Output.Company)
	fmt.Fprint(res, `<description>`+Catalog_to_Output.Description+`</description>`)
	fmt.Fprint(res, `<id>`, Catalog_to_Output.ID, `</id>`)
	fmt.Fprint(res, `</catalog>`)
}

// Call: /api/book.xml
// Description:
// Call will return an xml view on a singular book. ID is a well-formatted integer of an existing book id.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_GetBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	BookID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	ctx := appengine.NewContext(req)
	Book_to_Output, geterr := GetBookFromDatastore(ctx, BookID)
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

// Call: /api/chapter.xml
// Description:
// Call will return an xml view on a singular chapter. ID is a well-formatted integer of an existing chapter id.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_GetChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ChapterID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	ctx := appengine.NewContext(req)
	Chapter_to_Output, geterr := GetChapterFromDatastore(ctx, ChapterID)
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
	fmt.Fprintf(res, `<order>%d</order>`, Chapter_to_Output.Order)
	fmt.Fprint(res, `</chapter>`)
}

// Call: /api/section.xml
// Description:
// Call will return an xml view on a singular section. ID is a well-formatted integer of an existing section id.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_GetSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	SectionID, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if convErr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	ctx := appengine.NewContext(req)
	Section_to_Output, geterr := GetSectionFromDatastore(ctx, SectionID)
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
	fmt.Fprintf(res, `<order>%d</order>`, Section_to_Output.Order)
	fmt.Fprint(res, `</section>`)
}

// Call: /api/objective.html
// Description:
// Call will return an xml view on a singular objective. ID is a well-formatted integer of an existing objective id.
//
// Method: GET
// Results: HTML Snippet
// Mandatory Options: ID
// Optional Options:
// Codes:
//      Failure: HTML<section> that describes the error.
//      Success: HTML<section> of objective information.
func API_GetObjectiveHTML(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil {
		fmt.Fprint(res, `<section><p>Request has failed: Invalid ID.</p></section>`)
		return
	}

	ctx := appengine.NewContext(req)
	objectiveToScreen, getErr := GetObjectiveFromDatastore(ctx, int64(ObjectiveID))
	//HandleError(res, getErr)
	if getErr != nil {
		fmt.Fprint(res, `<section><p>Request has failed: No objective with given ID.</p></section>`)
		return
	}

	ServeTemplateWithParams(res, "ObjectiveHTML.html", objectiveToScreen)
}

// Call: /api/exercise.xml
// Description:
// Call will return an xml view on a singular Exercise.
// ID is a well-formed string of an existing exercise.
//
// Method: GET
// Results: XML
// Mandatory Options: ID
// Optional Options:
// Codes:
//      XML<status> Failure - read <message> of error for more information
//      Success will return the well formed xml.
func API_GetExercise(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ExerciseID, _ := strconv.Atoi(req.FormValue("ID"))
	if ExerciseID == 0 {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>Invalid ID</message></error>`)
		return
	}
	ctx := appengine.NewContext(req)
	Exercise_to_Output, geterr := GetExerciseFromDatastore(ctx, int64(ExerciseID))
	if geterr != nil {
		fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8" ?><error><status>Failure</status><message>ID Not Found!</message></error>`)
		return
	}
	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><exercise>`)
	fmt.Fprintf(res, `<instruction>%s</instruction>`, Exercise_to_Output.Instruction)
	fmt.Fprint(res, `<question>`+Exercise_to_Output.Question+`</question>`)
	fmt.Fprint(res, `<solution>`+Exercise_to_Output.Solution+`</solution>`)
	fmt.Fprint(res, `<answer>`, Exercise_to_Output.Answer, "</answer>")
	fmt.Fprintf(res, `<parent>%d</parent>`, Exercise_to_Output.Parent)
	fmt.Fprintf(res, `<id>%d</id>`, Exercise_to_Output.ID)
	fmt.Fprintf(res, `<order>%d</order>`, Exercise_to_Output.Order)
	fmt.Fprint(res, `</exercise>`)
}
