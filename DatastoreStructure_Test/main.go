package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template // This is the storage location for all of our html files

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", home)
	r.GET("/createCatalog", createCatalog) // initalizes catalog
	r.GET("/readCatalog", getCatalog)      // pulls catalog from datastore
	r.GET("/createBook", createBook)       // expects a query of /createBook?name=<Name of new book>
	r.GET("/readBooks", getBooks)          // will query current catalog for all books.

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	pages = template.Must(pages.ParseGlob("html/*.html"))
}

// **************************************
// URL Handlers

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	err := pages.ExecuteTemplate(res, "index.html", nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func serveTemplateWithParams(res http.ResponseWriter, req *http.Request, templateName string, params interface{}) {
	// simple func to cut down on repeating code.
	err := pages.ExecuteTemplate(res, templateName, &params)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func makeCatalogKey(ctx context.Context) *datastore.Key {
	// we only use a single 'catalog' in this example. the default_catalog of namespace Catalogs
	// This is a simple function to make a recallable default_catalog key
	return datastore.NewKey(ctx, "Catalogs", "default_catalog", 0, nil)
}

func createCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Really hard to call your ancestor if they don't exist.
	// Use this as a one time call to make the default_catalog.
	// Nothing really different than a normal Put call
	var rootCatalog Catalog
	rootCatalog.Name = "This is the root catalog name"

	ctx := appengine.NewContext(req)
	catKey := makeCatalogKey(ctx)

	_, err := datastore.Put(ctx, catKey, &rootCatalog)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	serveTemplateWithParams(res, req, "printme.html", "Catalog is now created!")
}

func getCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// simple datastore.Get for the default_catalog
	var rootCatalog Catalog

	ctx := appengine.NewContext(req)
	catKey := makeCatalogKey(ctx)

	err := datastore.Get(ctx, catKey, &rootCatalog)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	serveTemplateWithParams(res, req, "printme.html", rootCatalog)
}

func createBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Now we're getting into the magic of the task.
	// We will make a new book of title name taken in from the url string.
	// the specific query string is URL/createBook?name=<name-your-using>
	var newBook Book
	newBook.Name = req.URL.Query().Get("name") // read book name from name value in query string

	ctx := appengine.NewContext(req)
	catKey := makeCatalogKey(ctx) // we need to know the ancestor to tie this value to.

	bookKey := datastore.NewKey(ctx, "Books", newBook.Name, 0, catKey) // make a key for the namespace Books with the ID of the book name and an ancestor of the default_catalog

	_, err := datastore.Put(ctx, bookKey, &newBook) // from there. put the data in the datastore using the key.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	serveTemplateWithParams(res, req, "printme.html", "Book is now created!")
}

func getBooks(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Okay. True magic here. Querying based on Ancestor
	ctx := appengine.NewContext(req)
	catKey := makeCatalogKey(ctx)

	q := datastore.NewQuery("Books").Ancestor(catKey) // we will query namespace Books for all books with the ancestor of our default_catalog

	booklist := make([]Book, 0) // make a list of books. we're filling this out.
	for t := q.Run(ctx); ; {    // for values within the query as it's running
		var x Book
		_, qErr := t.Next(&x)       // read one query value into a temporary location
		if qErr == datastore.Done { // if no value was read but it called exit
			break // then exit.
		} else if qErr != nil { // if there was a real error
			http.Error(res, qErr.Error(), http.StatusInternalServerError) // raise that error
		}
		booklist = append(booklist, x) // add the successful book found onto our output list
	}

	msg := ScreenMessage{"Books Gathered", booklist} // make the screen information

	serveTemplateWithParams(res, req, "printme.html", &msg) // send it to screen.
}

type Catalog struct { // simple catalog struct. will be increased later
	Name string
}

type Book struct { // same for books
	Name string
}

type ScreenMessage struct { // sinple message out to screen. I'm using this to read books made.
	Message string
	data    []Book
}

// Links used for information gathering:
// https://cloud.google.com/appengine/docs/go/datastore/entities
// https://cloud.google.com/appengine/docs/go/datastore/reference
// https://golang.org/pkg/net/url/#URL
