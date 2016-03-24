package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
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
func API_MakeCatalog() {}
func API_MakeBook()    {}
func API_MakeSection() {}

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

func API_GetSectionData() {}
