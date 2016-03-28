package main

import (
	"html/template"
	"net/http"
	"strconv"

	//"fmt"
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
	r.GET("/api/books.json", API_GetBookData)
	r.GET("/api/catalogs.json", API_GetCatalogData)
	r.GET("/api/chapters.json", API_GetChapterData)
	r.GET("/select", selectBookFromForm)
	r.POST("/api/makeCatalog", API_MakeCatalog)
	r.POST("/api/makeBook", API_MakeBook)

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
// Keys
func MakeCatalogKey(ctx context.Context, keyname string) *datastore.Key {
	return datastore.NewKey(ctx, "Catalogs", keyname, 0, nil)
}
func MakeBookKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Books", "", id, nil)
}
func MakeChapterKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Chapters", "", id, nil)
}

func initalizeData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	catalogTitles := []string{"default_catalog", "Math", "Science", "***"}
	chapterTitles := []string{"Nothing", "Making of", "Readme", "Sometimes always", "Nevermore", "A New Begining", "The Founding of the three states", "Taking over the Tri-State Area!", "Finally", "The End!", "Only when your down", "Over and Out", "Chapter titles are harder than book titles", "Part 1: Part 2", "Part 2: Part 1 again", "Integration", "Newtons Method"}

	for _, k := range catalogTitles {
		ck := MakeCatalogKey(ctx, k)
		cc := Catalog{"Basic Catalog", 0, "eduNet"}
		_, err := datastore.Put(ctx, ck, &cc)
		HandleError(res, err)
	}

	for i, title := range []string{"Hello ", "World", "A list", "Of titles", "The Hobbit", "Lord of the Trees", "A brand new cat", "Gone with the start", "Not on your life", "Bores", "Party Time with Party Pete: A Ride Of Your Life: Not for your pets!", "Marko Polo, Silly Game or Deadly Secret?", "Starbucks, The REAL addiction"} {
		bookInput := Book{}
		bookInput.Title = title
		bookInput.CatalogTitle = catalogTitles[(i % 4)]
		bk := MakeBookKey(ctx, 0)
		k, err2 := datastore.Put(ctx, bk, &bookInput)
		HandleError(res, err2)
		for ii := 0; ii < 3; ii += 1 {
			chapterInput := Chapter{}
			chapterInput.Title = chapterTitles[(int(k.IntID())+ii)%15] // trying some hashing functions to psuedo random the chapter titles.
			chapterInput.BookID = k.IntID()
			ck := MakeChapterKey(ctx, 0)
			_, err3 := datastore.Put(ctx, ck, &chapterInput)
			HandleError(res, err3)
		}
	}

	ServeTemplateWithParams(res, req, "printme.html", "Datastore has been initalized!")
}

func selectBookFromForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ServeTemplateWithParams(res, req, "bookSelection.html", nil)
}
func bookSelectedFromForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	i, numErr := strconv.Atoi(req.FormValue("BookID"))
	HandleError(res, numErr)
	// bookKey := MakeBookKey(ctx, int64(i))
	ServeTemplateWithParams(res, req, "printme.html", i)
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
