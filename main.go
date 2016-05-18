package main

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
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

	// Module: Core Structure
	// Files: main.go, STRUCT_Handlers.go
	/*************************************/
	r.GET("/", home)                     // <user> Root page
	r.GET("/catalogs", getCatalogsPage)  // <user> Catalogs listing
	r.GET("/select", selectBookFromForm) // <user> select objective based on information
	r.GET("/toc/:ID", getSimpleTOC)      // <user> user viewable toc for a book
	r.GET("/about", getAboutPage)        // <user> About Page
	r.GET("/favicon.ico", favIcon)       // <user> favicon

	// Module: Authentication/Session
	// Files: AUTH_authentication.go
	/****************************************/
	r.GET("/login", AUTH_Login_GET)         // <user> User Login
	r.GET("/logout", AUTH_Logout_GET)       // <user> User Logout
	r.GET("/register", AUTH_Register_GET)   // <user> Register New Users/Modify existing users
	r.POST("/register", AUTH_Register_POST) // <user><auth> Post to make the new user
	r.GET("/user", AUTH_UserInfo)           // <user><auth><DEBUG> DEBUG user info

	// Module: Structure Readers
	// Files: main.go, STRUCT_Handlers.go
	/*****************************************************/
	r.GET("/read/exercise/:ID", getSimpleExerciseReader)   // <user> read exercise given id
	r.GET("/read/objective/:ID", getSimpleObjectiveReader) // <user> read objective given id
	r.GET("/preview", getObjectivePreview)                 // <user><OUTDATED> preview objective given id

	// Module: Structure Modifiers
	// Files: STRUCT_Handlers.go
	/***********************************************/
	r.GET("/edit/catalog/:ID", getCatalogEditor)           // <user><auth> Modify Catalog Information
	r.GET("/edit/book/:ID", getBookEditor)                 // <user><auth> Modify Book Information
	r.GET("/edit/chapter/:ID", getChapterEditor)           // <user><auth> Modify Chapter Information
	r.GET("/edit/section/:ID", getSectionEditor)           // <user><auth> Modify Section Information
	r.GET("/edit/exercise/:ID", getExerciseEditor)         // <user><auth> Modify Exercise Information
	r.GET("/edit/objective/:ID", getSimpleObjectiveEditor) // <user><auth> edit objective given id

	// Module: Structure Parser
	// Files: PARSE_BookParser.go
	/************************************************/
	r.GET("/export/:ID", exportBookToScreen)        // <user><DEBUG>
	r.GET("/import/book", PARSE_GET_FileUploader)   // <DEBUG>
	r.POST("/import/book", PARSE_POST_FileUploader) // <DEBUG

	// Module: Images
	// Files: Images.go
	/********************************************************************/
	r.GET("/image", IMAGE_API_GetImageFromCS)                           // <user> image requester /api/getImage
	r.GET("/api/getImage", IMAGE_API_GetImageFromCS)                    // <DEBUG> Duplicate of /image, *outdated*
	r.GET("/image/browser", IMAGE_BrowserForm)                          // <user> image browser
	r.GET("/image/uploader", IMAGE_PostUploadForm)                      // <user-internal><auth> image uploader
	r.POST("/api/create/image", IMAGE_API_PlaceImageIntoCS)             // <api><auth> image creator
	r.POST("/api/delete/image", IMAGE_API_RemoveImageFromCS)            // <api><auth> image deleter
	r.POST("/api/ckeditor/create", IMAGE_API_CKEDITOR_PlaceImageIntoCS) // <api><auth> ckEditor, image creator

	// Module: API-Readers, Collection
	// Files: API_Readers.go
	/*************************************************/
	r.GET("/api/catalogs.json", API_GetCatalogs)     // <api> read datastore, catalogs
	r.GET("/api/books.json", API_GetBooks)           // <api> read datastore, books
	r.GET("/api/chapters.json", API_GetChapters)     // <api> read datastore, chapters
	r.GET("/api/sections.json", API_GetSections)     // <api> read datastore, sections
	r.GET("/api/objectives.json", API_GetObjectives) // <api> read datastore, objectives
	r.GET("/api/exercises.json", API_GetExercises)   // <api> read datatore, exercises
	r.GET("/api/toc.xml", API_getTOC)                // <api> xml toc for a book

	// Module: API-Readers, Singular
	// Files: API_Readers.go
	/***************************************************/
	r.GET("/api/catalog.xml", API_GetCatalog)          // <api> read datastore, catalog as xml
	r.GET("/api/book.xml", API_GetBook)                // <api> read datastore, book as xml
	r.GET("/api/chapter.xml", API_GetChapter)          // <api> read datastore, chapter as xml
	r.GET("/api/section.xml", API_GetSection)          // <api> read datastore, section as xml
	r.GET("/api/objective.html", API_GetObjectiveHTML) // <api> read datastore, objective as html
	r.GET("/api/exercise.xml", API_GetExercise)        // <api> read datastore, exercise as xml

	// Module: API-Writers
	// Files: API_Writers.go
	/************************************************/
	r.POST("/api/create/catalog", API_MakeCatalog)     // <api><auth> create datastore, catalog
	r.POST("/api/create/book", API_MakeBook)           // <api><auth> create datastore, book
	r.POST("/api/create/chapter", API_MakeChapter)     // <api><auth> create datastore, chapter
	r.POST("/api/create/section", API_MakeSection)     // <api><auth> create datastore, section
	r.POST("/api/create/objective", API_MakeObjective) // <api><auth> create datastore, objective
	r.POST("/api/create/exercise", API_MakeExercise)   // <api><auth> create datastore, exercise

	// Module: API-Deleters
	// Files: API_Deleters.go
	/****************************************************/
	r.POST("/api/delete/catalog", API_DeleteCatalog)     // <api><auth> delete datastore, catalog
	r.POST("/api/delete/book", API_DeleteBook)           // <api><auth> delete datastore, book
	r.POST("/api/delete/chapter", API_DeleteChapter)     // <api><auth> delete datastore, chapter
	r.POST("/api/delete/section", API_DeleteSection)     // <api><auth> delete datastore, section
	r.POST("/api/delete/objective", API_DeleteObjective) // <api><auth> delete datastore, objective
	r.POST("/api/delete/exercise", API_DeleteExercise)   // <api><auth> delete datastore, exercise

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
// Core Functionality, Handlers
/////

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

// Call: /about
// Description:
// Our about page
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func getAboutPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))
	ServeTemplateWithParams(res, req, "about.html", pu)
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

// Call: /catalogs
// Description:
// Our catalogs page
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options:
func getCatalogsPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	pu, _ := GetPermissionUserFromSession(appengine.NewContext(req))
	ServeTemplateWithParams(res, req, "catalogs.html", pu)
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

var (
	errorMessages = map[int]string{
		400: "Expected information missing. Ensure that all form information has values.",
		401: "You must login to complete this action.",
		406: "Incoming information invalid. Please try again.",
		500: "Internal Server Error. Try again in 30 seconds. If issue persists, please report bug.",
	}
)

func HandleErrorIntoPage(res http.ResponseWriter, req *http.Request, e error, action string) bool {
	if e != nil {
		screenOutput := struct {
			Recommend template.HTML
			MoreInfo  string
		}{
			template.HTML(action),
			e.Error(),
		}
		ServeTemplateWithParams(res, req, "error.gohtml", screenOutput)
		return true
	}
	return false
}

// Internal Function
// Passes along any information to templates and then executes them.
func ServeTemplateWithParams(res http.ResponseWriter, req *http.Request, templateName string, params interface{}) {
	err := pages.ExecuteTemplate(res, templateName, &params)
	HandleError(res, err)
}
