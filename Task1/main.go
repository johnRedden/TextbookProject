package main

import (
	//"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", home)
	r.GET("/showMessage", showMsg)
	r.GET("/makeMessage", makeMsg)
	r.POST("/makeMessage", uploadMsg)
	r.GET("/favicon.ico", favIcon)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(pages.ParseGlob("html/*.html"))
}

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	err := pages.ExecuteTemplate(res, "index.html", nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func showMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Retrive form data from datastore
	var messageFromDatastore MessageStructure

	ctx := appengine.NewContext(req)
	messageKey := datastore.NewKey(ctx, "Messages", "MessageID", 0, nil)
	datastoreErr := datastore.Get(ctx, messageKey, &messageFromDatastore)
	if datastoreErr != nil {
		messageFromDatastore.Data = "NO MESSAGE FOUND - " + datastoreErr.Error()
	}

	// post message onto form webpage.
	err := pages.ExecuteTemplate(res, "showMessage.html", messageFromDatastore.Data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

}

func makeMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post the form to edit data
	err := pages.ExecuteTemplate(res, "makeMessage.html", nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func uploadMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var messageToUpload MessageStructure
	messageToUpload.Data = req.FormValue("Message")

	// upload form data to datastore
	ctx := appengine.NewContext(req)
	messageKey := datastore.NewKey(ctx, "Messages", "MessageID", 0, nil)
	_, datastoreErr := datastore.Put(ctx, messageKey, &messageToUpload)

	if datastoreErr != nil {
		http.Error(res, datastoreErr.Error(), http.StatusInternalServerError)
	}

	// redirect to home
	http.Redirect(res, req, "/", http.StatusFound)
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

// *******************************
// Structures

type MessageStructure struct {
	Data string
}
