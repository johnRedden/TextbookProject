package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template // This is the storage location for all of our html files
var apiPages *template.Template

func init() {

	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", home)
	r.GET("/init", initalizeData)
	r.GET("/vb", viewBookKeys)
	r.GET("/vbs", viewBook)
	r.GET("/api/books.json", API_GetBookData)
	r.GET("/api/catalogs.json", API_GetCatalogData)

	r.GET("/favicon.ico", favIcon)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	pages = template.Must(pages.ParseGlob("templates/*.*"))
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

func ServeTemplateWithParams(res http.ResponseWriter, req *http.Request, templateName string, params interface{}) {
	// simple func to cut down on repeating code.
	err := pages.ExecuteTemplate(res, templateName, &params)
	HandleError(res, err)
}

func HandleError(res http.ResponseWriter, e error) {
	// generic error handling for any error we encounter plus a message we've defined.
	if e != nil {
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, req, "index.html", nil)
}

// *************************************
func makeCatalogKey(ctx context.Context, keyname string) *datastore.Key {
	return datastore.NewKey(ctx, "Catalogs", keyname, 0, nil)
}
func makeBookKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Books", "", id, nil)
}
func makeSectionKey(ctx context.Context, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(ctx, "Sections", "", 0, parent)
}

func initalizeData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	catalogTitles := []string{"default_catalog", "Math", "Science", "***"}

	for _, k := range catalogTitles {
		ck := makeCatalogKey(ctx, k)
		cc := Catalog{"Basic Catalog", 0, "eduNet"}
		_, err := datastore.Put(ctx, ck, &cc)
		HandleError(res, err)
	}

	for i, title := range []string{"Hello ", "World", "A list", "Of titles", "The Hobbit", "Lord of the Trees", "A brand new cat", "Gone with the start", "Not on your life", "Bores", "Party Time with Party Pete: A Ride Of Your Life: Not for your pets!", "Marko Polo, Silly Game or Deadly Secret?", "Starbucks, The REAL addiction"} {
		bookInput := Book{}
		bookInput.Title = title
		bookInput.CatalogTitle = catalogTitles[(i % 4)]
		_, err2 := datastore.Put(ctx, makeBookKey(ctx, 0), &bookInput)
		HandleError(res, err2)
	}

	ServeTemplateWithParams(res, req, "printme.html", "Datastore has been initalized!")
}

type BookKeys struct {
	Key int64
	Book
}

func viewBookKeys(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Books")

	booklist := make([]BookKeys, 0)
	for t := q.Run(ctx); ; {
		var x Book
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		booklist = append(booklist, BookKeys{k.IntID(), x})
	}

	ServeTemplateWithParams(res, req, "printme.html", booklist)
}

func viewBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	i, err := strconv.Atoi(req.FormValue("Key"))
	HandleError(res, err)

	ctx := appengine.NewContext(req)

	bookKey := makeBookKey(ctx, int64(i))
	var x Book
	e := datastore.Get(ctx, bookKey, &x)
	HandleError(res, e)

	ServeTemplateWithParams(res, req, "printme.html", x)

	// ctx := appengine.NewContext(req)
	// q := datastore.NewQuery("Books")

	// booklist := make([]BookKeys, 0)
	// for t := q.Run(ctx); ; {
	// 	var x Book
	// 	k, qErr := t.Next(&x)
	// 	if qErr == datastore.Done {
	// 		break
	// 	} else if qErr != nil {
	// 		http.Error(res, qErr.Error(), http.StatusInternalServerError)
	// 	}
	// 	booklist = append(booklist, BookKeys{k.IntID(), x})
	// }

	// ServeTemplateWithParams(res, req, "printme.html", booklist)
}

// **************************************
// URL Handlers

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
