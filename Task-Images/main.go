package main

import (
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"html/template"
	"net/http"
	"strconv"
)

var pages *template.Template // This is the storage location for all of our html files

func init() {

	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/image", IMAGE_PostUploadForm)
	r.POST("/image", IMAGE_API_PlaceImageIntoCS)
	r.GET("/api/getImage", IMAGE_API_GetImageFromCS)
	r.POST("/api/makeImage", IMAGE_API_PlaceImageIntoCS)

	r.GET("/api/catalogs.json", API_GetCatalogs) // API GET json calls
	r.GET("/api/books.json", API_GetBooks)
	r.GET("/api/chapters.json", API_GetChapters)
	r.GET("/api/sections.json", API_GetSections)
	r.GET("/api/objectives.json", API_GetObjectives)

	r.GET("/api/objective.html", API_GetObjectiveHTML) // API GET html calls

	r.POST("/api/makeCatalog", API_MakeCatalog) // API POST make calls
	r.POST("/api/makeBook", API_MakeBook)
	r.POST("/api/makeChapter", API_MakeChapter)
	r.POST("/api/makeSection", API_MakeSection)
	r.POST("/api/makeObjective", API_MakeObjective)

	r.GET("/", home)
	r.GET("/select", selectBookFromForm)
	r.GET("/edit", getSimpleObjectiveEditor)
	r.GET("/read", getSimpleObjectiveReader)
	r.GET("/favicon.ico", favIcon)

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
	// Req: GET /select
	ServeTemplateWithParams(res, req, "simpleSelector.html", nil)
}

func getSimpleObjectiveEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Req: GET /edit?ID=<Objective ID Number>
	// ctx := appengine.NewContext(req)

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
	ServeTemplateWithParams(res, req, "simpleReader.html", nil)
}
