package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	//"google.golang.org/appengine/datastore"
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

	// post message onto form webpage.
	err := pages.ExecuteTemplate(res, "showMessage.html", "Test Message")
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

	// ctx := appengine.NewContext(req)
	// key := datastore.NewKey(ctx, "Users", req.FormValue("Message"), 0, nil)

	fmt.Println("Data Taken from submit", req.FormValue("Message"))
	// upload form data to datastore

	// redirect to home
	err := pages.ExecuteTemplate(res, "showMessage.html", req.FormValue("Message"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	//http.Redirect(res, req, "/", http.StatusTemporaryRedirect)

}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}
