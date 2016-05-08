package main

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"strconv"
)

// pages is a local storage variable for all of our executable templates.
// These templates should be called using ServeTemplateWithParams()
var pages *template.Template

func init() {
	r := httprouter.New()
	http.Handle("/", r)

	//// ---------------------------------------------------------- //
	// Handlers
	//
	// This segment below is the master list of all url handlers this
	// website accepts. If it is not here, it is internal.
	//
	// Tags:
	//  * <user>: A page the user can interact with
	//  * <user-internal>: A page that the user does not directly interact with but depends on.
	//  * <api>: A request that preforms background actions.
	//  * <auth>: (Modifier) This request requires user authentication.
	//  * <DEBUG>: (Modifier) This handler is to be treated as temporary, used in development only.
	//// ---------------------------------------------------------- //

	// Module: Images
	// Files: Images.go
	/********************************************************************/
	r.GET("/image", IMAGE_API_GetImageFromCS)                           // <user> image requester /api/getImage
	r.GET("/api/getImage", IMAGE_API_GetImageFromCS)                    // <DEBUG> Duplicate of /image, *outdated*
	r.GET("/image/browser", IMAGE_BrowserForm)                          // <user> image browser
	r.GET("/image/uploader", IMAGE_PostUploadForm)                      // <user-internal><auth> image uploader
	r.POST("/api/makeImage", IMAGE_API_PlaceImageIntoCS)                // <api><auth> image creator
	r.POST("/api/deleteImage", IMAGE_API_RemoveImageFromCS)             // <api><auth> image deleter
	r.POST("/api/ckeditor/create", IMAGE_API_CKEDITOR_PlaceImageIntoCS) // <api><auth> ckEditor, image creator

	// Module: API-Readers, Collection
	// Files: API_Readers.go
	/*************************************************/
	r.GET("/api/catalogs.json", API_GetCatalogs)     // <api> read datastore, catalogs
	r.GET("/api/books.json", API_GetBooks)           // <api> read datastore, books
	r.GET("/api/chapters.json", API_GetChapters)     // <api> read datastore, chapters
	r.GET("/api/sections.json", API_GetSections)     // <api> read datastore, sections
	r.GET("/api/objectives.json", API_GetObjectives) // <api> read datastore, objectives
	r.GET("/toc", API_getTOC)                        // <api> xml toc for a book

	// Module: API-Readers, Singular
	// Files: API_Readers.go
	/***************************************************/
	r.GET("/api/catalog.xml", API_GetCatalog)          // <api> read datastore, catalog as xml
	r.GET("/api/book.xml", API_GetBook)                // <api> read datastore, book as xml
	r.GET("/api/chapter.xml", API_GetChapter)          // <api> read datastore, chapter as xml
	r.GET("/api/section.xml", API_GetSection)          // <api> read datastore, section as xml
	r.GET("/api/objective.html", API_GetObjectiveHTML) // <api> read datastore, objective as html

	// Module: API-Writers
	// Files: API_Writers.go
	/************************************************/
	r.POST("/api/makeCatalog", API_MakeCatalog)     // <api><auth> create datastore, catalog
	r.POST("/api/makeBook", API_MakeBook)           // <api><auth> create datastore, book
	r.POST("/api/makeChapter", API_MakeChapter)     // <api><auth> create datastore, chapter
	r.POST("/api/makeSection", API_MakeSection)     // <api><auth> create datastore, section
	r.POST("/api/makeObjective", API_MakeObjective) // <api><auth> create datastore, objective

	// Module: API-Deleters
	// Files: API_Deleters.go
	/****************************************************/
	r.POST("/api/deleteCatalog", API_DeleteCatalog)     // <api><auth> delete datastore, catalog
	r.POST("/api/deleteBook", API_DeleteBook)           // <api><auth> delete datastore, book
	r.POST("/api/deleteChapter", API_DeleteChapter)     // <api><auth> delete datastore, chapter
	r.POST("/api/deleteSection", API_DeleteSection)     // <api><auth> delete datastore, section
	r.POST("/api/deleteObjective", API_DeleteObjective) // <api><auth> delete datastore, objective

	// Module: Structure Modifiers
	// Files: main.go
	/*********************************************/
	r.GET("/edit/Catalog/:ID", getCatalogEditor) // <user><auth> Modify Catalog Information
	r.GET("/edit/Book/:ID", getBookEditor)       // <user><auth> Modify Book Information
	r.GET("/edit/Chapter/:ID", getChapterEditor) // <user><auth> Modify Chapter Information
	r.GET("/edit/Section/:ID", getSectionEditor) // <user><auth> Modify Section Information

	// Module: Core Structure
	// Files: main.go
	/*****************************************/
	r.GET("/", home)                         // <user> Root page
	r.GET("/select", selectBookFromForm)     // <user> select objective based on information
	r.GET("/edit", getSimpleObjectiveEditor) // <user><auth> edit objective given id
	r.GET("/read", getSimpleObjectiveReader) // <user> read objective given id
	r.GET("/preview", getObjectivePreview)   // <user> preview objective given id
	r.GET("/toc.html/:ID", getSimpleTOC)     // <user> user viewable toc for a book
	r.GET("/favicon.ico", favIcon)           // <user> favicon

	// Module: Authentication/Session
	// Files: AUTH_authentication.go
	/****************************************/
	r.GET("/login", AUTH_Login_GET)         // <user> User Login
	r.GET("/logout", AUTH_Logout_GET)       // <user> User Logout
	r.GET("/register", AUTH_Register_GET)   // <user> Register New Users/Modify existing users
	r.POST("/register", AUTH_Register_POST) // <user><auth> Post to make the new user
	r.GET("/user", AUTH_UserInfo)           // <user><auth><DEBUG> DEBUG user info

	// Module: Administration, Console and Commands
	// Files: ADMIN_administration.go
	/************************************************************/
	r.GET("/admin", ADMIN_AdministrationConsole)                // <user><auth> Admin Console
	r.POST("/admin/changeUsrPerm", ADMIN_POST_ELEVATEUSER)      // <api><auth> Admin: Change User Permissions
	r.GET("/admin/getUsrPerm", ADMIN_GET_USERPERM)              // <api><auth> Admin: Retrive User Permissions
	r.POST("/admin/forceUsrLogout", ADMIN_POST_ForceUserLogout) // <api><auth> Admin: Force a user to log out.
	r.POST("/admin/deleteUsr", ADMIN_POST_DELETEUSER)           // <api><auth> Admin: Delete a user. Will require said user to re-register.

	// Public file handling.
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	// Prepare templates.
	pages = template.Must(pages.ParseGlob("templates/*.*"))
}

// ------------------------------------
// Helper Functions
/////

// Internal Function
// generic error handling for any error we encounter.
func HandleError(res http.ResponseWriter, e error) {
	if e != nil {
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

// Internal Function
// generic error handling for any error we encounter plus a message we've defined.
// This sends a log out to appengine.
func HandleErrorWithLog(res http.ResponseWriter, e error, tag string, ctx context.Context) {
	if e != nil {
		log.Criticalf(ctx, "%s: %v", tag, e)
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

// Internal Function
// Passes along any information to templates and then executes them.
func ServeTemplateWithParams(res http.ResponseWriter, req *http.Request, templateName string, params interface{}) {
	err := pages.ExecuteTemplate(res, templateName, &params)
	HandleError(res, err)
}

// ------------------------------------
// Core Functionality, Handlers
/////

// Call: /favicon
// Description:
// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
//
// Method: GET
// Results: HTTP.Redirect
// Mandatory Options:
// Optional Options:
func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

// Call: /
// Description:
// Our home page
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))
	ServeTemplateWithParams(res, req, "index.html", pu)
}

// Call: /select
// Description:
// The selector page for site structures.
// Outdated?
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func selectBookFromForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, req, "simpleSelector.html", nil)
}

// Call: /edit
// Description:
// Our editor page for objectives given a valid objective id.
// Mandatory:ID must be a well-formatted integer of an existing objective id.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleObjectiveEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		http.Error(res, permErr.Error(), http.StatusUnauthorized)
		return
	}

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

// Call: /edit/Catalog/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getCatalogEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetCatalogFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Catalog.html", itemToScreen)
}

// Call: /edit/Book/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getBookEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetBookFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Book.html", itemToScreen)
}

// Call: /edit/Chapter/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getChapterEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetChapterFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Chapter.html", itemToScreen)
}

// Call: /edit/Section/:ID
// Description:
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSectionEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	editID := params.ByName("ID")
	i, parseErr := strconv.Atoi(editID)
	HandleError(res, parseErr)

	if i == 0 {
		http.Error(res, "Invalid ID", http.StatusExpectationFailed)
		return
	}

	itemToScreen, getErr := GetSectionFromDatastore(req, int64(i))
	HandleError(res, getErr)

	ServeTemplateWithParams(res, req, "editor_Section.html", itemToScreen)
}

// Call: /read
// Description:
// Our reader page for objectives.
// Mandatory:ID has no requirements on this level. Sub levels will
// require that objective ID exists and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleObjectiveReader(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	readID := req.FormValue("ID")
	ServeTemplateWithParams(res, req, "simpleReader.html", readID)
}

// Call: /toc.html/:ID
// Description:
// Our table of contents page for a book. Currently, this is handling the same
// as a selector for objectives to read/edit.
//
// Mandatory:ID has no requirements on this level. Sub levels will
// require that objective ID exists and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getSimpleTOC(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))

	screenOutput := struct {
		Name       string
		Email      string
		Permission int
		ID         string
	}{
		pu.Name,
		pu.Email,
		pu.Permission,
		params.ByName("ID"),
	}

	ServeTemplateWithParams(res, req, "toc.html", screenOutput)
}

// Call: /preview
// Description:
// Our simple page preview.
//
// Mandatory:ID must be an existing objective ID and is a well-formatted integer.
//
// Method: GET
// Results: HTML
// Mandatory Options: ID
// Optional Options:
func getObjectivePreview(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	objKey, convErr := strconv.ParseInt(req.FormValue("ID"), 10, 64)
	HandleError(res, convErr)
	objToScreen, err := GetObjectiveFromDatastore(req, objKey)
	HandleError(res, err)
	ServeTemplateWithParams(res, req, "preview.html", objToScreen)
}
