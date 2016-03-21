package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template // This is the storage location for all of our html files

func init() {

	r := httprouter.New()
	http.Handle("/", r)

	r.GET("/", home)
	r.POST("/postData", savePostedObjectiveData)
	r.GET("/preview", showObjective)

	r.GET("/favicon.ico", favIcon)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(pages.ParseGlob("html/*.html"))

}

// **************************************
// URL Handlers

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {


	err := pages.ExecuteTemplate(res, "editor.html", nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func savePostedObjectiveData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var objectiveToUpload Objective  
	objectiveToUpload.Objective = req.FormValue("objective")
	objectiveToUpload.Content = req.FormValue("content")
	objectiveToUpload.Author = req.FormValue("author")
	objectiveToUpload.Version = req.FormValue("version")

	ctx := appengine.NewContext(req)
	messageKey := datastore.NewKey(ctx, "Objective", "ObjectiveID", 0, nil) 
	_, datastoreErr := datastore.Put(ctx, messageKey, &objectiveToUpload)

	if datastoreErr != nil { // if datastore throws an error, handle it.
		http.Error(res, datastoreErr.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(res, objectiveToUpload.Content)

	
}
func showObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var obj Objective 

	ctx := appengine.NewContext(req)
	messageKey := datastore.NewKey(ctx, "Objective", "ObjectiveID", 0, nil) 
	datastoreErr := datastore.Get(ctx, messageKey, &obj) 
	if datastoreErr != nil {                                            
		http.Error(res, datastoreErr.Error(), http.StatusInternalServerError)
	}
	
	err := pages.ExecuteTemplate(res, "preview.html", obj)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

// *******************************
type Objective struct {
	Objective string
	Version   string `datastore:,noindex`
	Author       string  //separate objectives may have different authors.
	Content      string `datastore:,noindex`
	KeyTakeaways string // or array of strings
	Rating       int    // out of 5 stars
}
