package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	r.get("/favicon.ico", favIcon)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(tpl.ParseGlob("html/*.html"))
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
}

func makeMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post the form to edit data
}

func uploadMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// upload form data to datastore
	// redirect to home
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	http.Redirect(res, req, "public/images/favicon.ico", 302)
}
