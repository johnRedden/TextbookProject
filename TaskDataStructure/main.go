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
	r.GET("/", home) 
	r.POST("/", homeAgain)                                                                  
                                                       
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
	ctx := appengine.NewContext(req)                                    
	catKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil) 
	_, datastoreErr := datastore.Put(ctx, catKey, &catalog)  
	if datastoreErr != nil { 
		http.Error(res, datastoreErr.Error(), http.StatusInternalServerError)
	}
	
	err := pages.ExecuteTemplate(res, "index.html", catalog) 
	if err != nil {                                      
		http.Error(res, err.Error(), http.StatusInternalServerError) 
	}
}

func homeAgain(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	var dog Catalog  // dog just to show that the get works here

	ctx := appengine.NewContext(req)                                   
	mKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil)  
	datastoreErr := datastore.Get(ctx, mKey, &dog) 
	if datastoreErr != nil {    
			dog.Name = "NO MESSAGE FOUND - " + datastoreErr.Error()
	}

	var book Book
	book.Title = req.FormValue("BookName") 
	
	//q := datastore.NewQuery("Catalog").Ancestor(mKey).Order("-Name").Limit(10)
	// I am confused!!
	// 1. I do not know how to add child book to the catalog
	// 2. Not sure how to pass more than one parameter to the template.
	// 3. Not sure the best way to debug, test and code with GO.  I need to see the objects somehow?  In javascript I use console.log(obj) to find GO equivalent.

	err := pages.ExecuteTemplate(res, "index.html", dog) 
	if err != nil {                                      
		http.Error(res, err.Error(), http.StatusInternalServerError) 
	}
}


func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

// *******************************
// Frist try at a datastore design  Parent->child->down the list to objective level.
type Catalog struct { 
	Name string
	Version float32
	Company string
}
type Book struct { 
	Title string
	Version float32
	Author string  // or array of strings
}
type Chapter struct { 
	Title string
	Version float32
}
type Section struct { 
	Title string
	Version float32
}
type Objective struct { 
	objective string
	Version float32
	Author string  //or array of strings
	Content string
	KeyTakeaways string // or array of strings
	rating int  // out of 5 stars
}
