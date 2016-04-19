package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"strconv"
)

var pages *template.Template // This is the storage location for all of our html files

func init() {

	r := httprouter.New()
	http.Handle("/", r)

	// Images.go
	r.GET("/image", IMAGE_BrowserForm)                                  // image browser <user>
	r.GET("/image/uploader", IMAGE_PostUploadForm)                      // image uploader <user-internal>
	r.GET("/api/getImage", IMAGE_API_GetImageFromCS)                    // image requester <user-internal>
	r.POST("/api/makeImage", IMAGE_API_PlaceImageIntoCS)                // image creator <api>
	r.POST("/api/ckeditor/create", IMAGE_API_CKEDITOR_PlaceImageIntoCS) // ckEditor, image creator <api>

	// API.go, readers
	r.GET("/api/catalogs.json", API_GetCatalogs)       // read datastore, catalogs <api>
	r.GET("/api/books.json", API_GetBooks)             // read datastore, books <api>
	r.GET("/api/chapters.json", API_GetChapters)       // read datastore, chapters <api>
	r.GET("/api/sections.json", API_GetSections)       // read datastore, sections <api>
	r.GET("/api/objectives.json", API_GetObjectives)   // read datastore, objectives <api>
	r.GET("/api/objective.html", API_GetObjectiveHTML) // read datastore, objective as html <api>

	// API.go, writers
	r.POST("/api/makeCatalog", API_MakeCatalog)     // create datastore, catalog <api>
	r.POST("/api/makeBook", API_MakeBook)           // create datastore, book <api>
	r.POST("/api/makeChapter", API_MakeChapter)     // create datastore, chapter <api>
	r.POST("/api/makeSection", API_MakeSection)     // create datastore, section <api>
	r.POST("/api/makeObjective", API_MakeObjective) // create datastore, objective <api>

	// main.go
	r.GET("/", home)                         // Root page <user>
	r.GET("/select", selectBookFromForm)     // select objective based on information <user>
	r.GET("/edit", getSimpleObjectiveEditor) // edit objective given id <user>
	r.GET("/read", getSimpleObjectiveReader) // read objective given id <user>

	// main.go, Table of Contents
	r.GET("/toc", getBookTOC)    // xml toc for a book <api>
	r.GET("/toc.html", printTOC) // user viewable toc for a book <user>

	r.GET("/favicon.ico", favIcon) // favicon <user>

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(pages.ParseGlob("templates/*.*"))
}

// ------------------------------------
// Helper Functions
/////

func HandleError(res http.ResponseWriter, e error) {
	// generic error handling for any error we encounter plus a message we've defined.
	if e != nil {
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

func HandleErrorWithLog(res http.ResponseWriter, e error, tag string, ctx context.Context) {
	// generic error handling for any error we encounter plus a message we've defined.
	if e != nil {
		log.Criticalf(ctx, "%s: %v", tag, e)
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

func ServeTemplateWithParams(res http.ResponseWriter, req *http.Request, templateName string, params interface{}) {
	// simple func to cut down on repeating code.
	err := pages.ExecuteTemplate(res, templateName, &params)
	HandleError(res, err)
}

// ------------------------------------
// Pages
/////

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, req, "index.html", nil)
}

func selectBookFromForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /select
	ServeTemplateWithParams(res, req, "simpleSelector.html", nil)
}

func getSimpleObjectiveEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /edit?ID=<Objective ID Number>

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

func getSimpleObjectiveReader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /read
	ServeTemplateWithParams(res, req, "simpleReader.html", nil)
}

func getBookTOC(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /toc?ID=<Book ID Number>

	/// - - - -
	// Initial Check, Ensure information is trivially good
	/////////

	BookID_In, numErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	if numErr != nil || BookID_In == 0 {
		http.Redirect(res, req, "/?status=invalid_id", http.StatusTemporaryRedirect)
	}

	/// - - - -
	// Gather Book information, ensure that book exists.
	////////

	BookTitle, BookCatalog, BookID_Out := func(req *http.Request, id int64) (string, string, int64) { // get book data
		book_to_output, _ := GetBookFromDatastore(req, id)
		return book_to_output.Title, book_to_output.CatalogTitle, book_to_output.ID
	}(req, BookID_In)

	if BookID_In != BookID_Out {
		ServeTemplateWithParams(res, req, "printme.html", "ERROR! Incoming id not found!")
		return
	}

	/// - - - -
	// Prepare to make everything simple.
	//////

	ctx := appengine.NewContext(req)
	type Title_ID struct { // struct for each layer.
		Title string
		ID    int64
	}
	gatherKindGroup := func(ctx context.Context, parentID int64, kind string) []Title_ID {
		q := datastore.NewQuery(kind)
		q = q.Filter("Parent =", parentID)
		q = q.Project("Title")

		output_chapters := make([]Title_ID, 0)
		for t := q.Run(ctx); ; {
			var cName struct{ Title string }
			k, qErr := t.Next(&cName)

			if qErr == datastore.Done {
				break
			} else if qErr != nil {
				http.Error(res, qErr.Error(), http.StatusInternalServerError)
			}

			output_chapters = append(output_chapters, Title_ID{cName.Title, k.IntID()})
		}
		return output_chapters
	}

	/// - - - -
	// Print header/Book information
	//////

	fmt.Fprint(res, `<?xml version="1.0" encoding="UTF-8"?><book>`) // Layer Book Information.
	fmt.Fprintf(res, `<booktitle>%s</booktitle><bookid>%d</bookid><catalog>%s</catalog>`, BookTitle, BookID_Out, BookCatalog)

	/// - - - -
	// Gather & Print Sub information as available
	//////

	for _, singleChapter := range gatherKindGroup(ctx, BookID_Out, "Chapters") {
		fmt.Fprintf(res, `<chapter><chaptertitle>%s</chaptertitle><chapterid>%d</chapterid>`, singleChapter.Title, singleChapter.ID)
		for _, singleSection := range gatherKindGroup(ctx, singleChapter.ID, "Sections") {
			fmt.Fprintf(res, `<section><sectiontitle>%s</sectiontitle><sectionid>%d</sectionid>`, singleSection.Title, singleSection.ID)
			for _, singleObjective := range gatherKindGroup(ctx, singleSection.ID, "Objectives") {
				fmt.Fprintf(res, `<objective><objectivetitle>%s</objectivetitle><objectiveid>%d</objectiveid></objective>`, singleObjective.Title, singleObjective.ID)
			}
			fmt.Fprint(res, `</section>`)
		}
		fmt.Fprint(res, `</chapter>`)
	}

	/// - - - -
	// Close Book
	//////

	fmt.Fprint(res, `</book>`) // Layer book close
}

func printTOC(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /toc.html?ID=<Book ID Number>
	ServeTemplateWithParams(res, req, "toc.html", req.FormValue("ID"))
}
