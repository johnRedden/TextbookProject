package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template // This is the storage location for all of our html files
var catalog Catalog

func init() {

	r := httprouter.New()
	http.Handle("/", r)
	// do we need both a get and a post?
	// answ: No, get is only for situations where we want to send information to the user, post on the otherhand is when we want information from the user to poll back to us--such as form values.
	r.GET("/", home)
	r.POST("/", homeAgain)
	r.POST("/test", test)

	r.GET("/favicon.ico", favIcon)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(pages.ParseGlob("html/*.html"))

	catalog.Name = "mainCatalog"
	catalog.Company = "eduNetSystems"
	catalog.Version = 1

}

// **************************************
// URL Handlers

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	// Upload data to datastore
	firstMessage := make([]string, 0)
	firstMessage = append(firstMessage, "enter a book")

	err := pages.ExecuteTemplate(res, "index.html", firstMessage)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func test(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//var deleteBook string
	var x string
	var newBook Book
	newBook.Title = req.FormValue("BookName")
	x = req.FormValue("delete")
	if x == "yes" {
		//delete book (this is totall insecure!)
		var dog Catalog // dog just to show that the get works here
		ctx := appengine.NewContext(req)

		catKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil)
		datastoreErr := datastore.Get(ctx, catKey, &dog)
		if datastoreErr != nil {
			dog.Name = "NO MESSAGE FOUND - " + datastoreErr.Error()
		}

		bookKey := datastore.NewKey(ctx, "Books", newBook.Title, 0, catKey)

		datastore.Delete(ctx, bookKey) // from there. put the data in the datastore using the key.
		/*			if err2 != nil {
						http.Error(res, err2.Error(), http.StatusInternalServerError)
					}
		*/

		x = "smack"
	}

	err := pages.ExecuteTemplate(res, "test.html", x)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func homeAgain(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var dog Catalog // dog just to show that the get works here

	ctx := appengine.NewContext(req)

	catKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil)
	datastoreErr := datastore.Get(ctx, catKey, &dog)
	if datastoreErr != nil {
		dog.Name = "NO MESSAGE FOUND - " + datastoreErr.Error()
	}

	//var deleteBook string
	var newBook Book
	newBook.Title = req.FormValue("BookName")
	//deleteBook = req.FormValue("delete")

	bookKey := datastore.NewKey(ctx, "Books", newBook.Title, 0, catKey)

	_, err2 := datastore.Put(ctx, bookKey, &newBook) // from there. put the data in the datastore using the key.
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
	}

	q := datastore.NewQuery("Books").Ancestor(catKey)

	booklist := make([]Book, 0) // make a list of books. we're filling this out.
	for t := q.Run(ctx); ; {    // for values within the query as it's running
		var x Book
		_, qErr := t.Next(&x)       // read one query value into a temporary location
		if qErr == datastore.Done { // if no value was read but it called exit
			break // then exit.
		} else if qErr != nil { // if there was a real error
			http.Error(res, qErr.Error(), http.StatusInternalServerError) // raise that error
		}
		x.Author = "me"
		booklist = append(booklist, x) // add the successful book found onto our output list
	}

	err := pages.ExecuteTemplate(res, "index.html", booklist)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}
