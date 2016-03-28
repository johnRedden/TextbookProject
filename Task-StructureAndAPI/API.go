package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
)

type catalogWithID struct {
	ID string
	Catalog
}
type bookWithID struct {
	ID int64
	Book
}
type chapterWithID struct {
	ID int64
	Chapter
}

// Post data calls. There should be relavant form data in these calls
func API_MakeCatalog(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Post call for making a catalog, we would check for a signed in user here.
	// Expects data from at minimum CatalogName
	// Also has data in from Company and Version
	// Version should be a well formed stringed float

	catalogName := req.FormValue("CatalogName")
	if catalogName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name"}`)
		return
	}
	// handle incoming data Company
	comp := req.FormValue("Company")
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}
	// Make our catalog
	catalogForDatastore := Catalog{}
	catalogForDatastore.Name = catalogName
	catalogForDatastore.Company = comp
	catalogForDatastore.Version = ver32
	// Get the datastore up and running!
	ctx := appengine.NewContext(req)

	ck := MakeCatalogKey(ctx, catalogName)
	_, errDatastore := datastore.Put(ctx, ck, &catalogForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":""}`)
}
func API_MakeBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Here is where we get a bit more complicated.
	// One stop shop for everything related to making a book
	// If you feed in a BookID it will update that specific id
	// If no id is given, will make a new one.
	// Taking in mandatory options CatalogName and BookName
	// Optional options Author,Version
	// TODO: Add in tags functionality

	bookID, numErr := strconv.Atoi(req.FormValue("BookID"))
	if numErr != nil {
		bookID = 0
	}

	catalogName := req.FormValue("CatalogName")
	if catalogName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Catalog Name"}`)
		return
	}

	bookName := req.FormValue("BookName")
	if bookName == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"Empty Book Name"}`)
		return
	}
	// handle incoming data Company
	auth := req.FormValue("Author")
	// handle incoming data Version
	ver64, errFloat := strconv.ParseFloat(req.FormValue("Version"), 32)
	var ver32 float32
	if errFloat != nil {
		ver32 = float32(ver64)
	}

	// TODO: Add in something to allow for tags.
	// string parsing maybe?

	bookForDatastore := Book{}
	bookForDatastore.Title = bookName
	bookForDatastore.CatalogTitle = catalogName
	bookForDatastore.Author = auth
	bookForDatastore.Version = ver32

	ctx := appengine.NewContext(req)

	bk := MakeBookKey(ctx, int64(bookID))
	_, errDatastore := datastore.Put(ctx, bk, &bookForDatastore)
	HandleError(res, errDatastore)

	fmt.Fprint(res, `{"result":"success","reason":""}`)

}
func API_MakeChapter(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

}
func API_MakeSection(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

}
func API_MakeObjective(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

}

// Get data calls, these will Fprint for reading
// Using this as ref: https://github.com/FelixVicis/f15_advWeb_finalProject/blob/master/04_Backend_Functionality/public/templates/signup.html
func API_GetCatalogData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Catalogs")
	cataloglist := make([]catalogWithID, 0)
	for t := q.Run(ctx); ; {
		var x Catalog
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		cataloglist = append(cataloglist, catalogWithID{k.StringID(), x})
	}

	ServeTemplateWithParams(res, req, "Catalogs.json", cataloglist)
}

func API_GetBookData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Books")

	queryCatalogName := req.FormValue("Catalog")
	if queryCatalogName != "" {
		q = q.Filter("CatalogTitle =", queryCatalogName)
	}

	booklist := make([]bookWithID, 0)
	for t := q.Run(ctx); ; {
		var x Book
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		booklist = append(booklist, bookWithID{k.IntID(), x})
	}

	ServeTemplateWithParams(res, req, "Books.json", booklist)
}

func API_GetChapterData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Chapters")

	queryBookID := req.FormValue("BookID")
	if queryBookID != "" { // Ensure that a BookID was indeed sent.
		i, numErr := strconv.Atoi(queryBookID) // does that BookID contain a number?
		HandleError(res, numErr)
		q = q.Filter("BookID =", int64(i))
	}

	chapterList := make([]chapterWithID, 0)
	for t := q.Run(ctx); ; {
		var x Chapter
		k, qErr := t.Next(&x)
		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			http.Error(res, qErr.Error(), http.StatusInternalServerError)
		}
		chapterList = append(chapterList, chapterWithID{k.IntID(), x})
	}

	ServeTemplateWithParams(res, req, "Chapters.json", chapterList)
}

func API_GetSectionData()   {}
func API_GetObjectiveData() {}
