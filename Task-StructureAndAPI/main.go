package main

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"html/template"
	"net/http"
	"strconv"
)

var pages *template.Template // This is the storage location for all of our html files

func init() {

	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", home)
	r.GET("/init", initalizeData)
	r.GET("/api/catalogs.json", API_GetCatalogs)
	r.GET("/api/books.json", API_GetBooks)
	r.GET("/api/chapters.json", API_GetChapters)
	r.GET("/api/sections.json", API_GetSections)
	r.GET("/api/objectives.json", API_GetObjectives)

	r.GET("/select", selectBookFromForm)
	r.GET("/edit", getSimpleObjectiveEditor)

	r.GET("/api/makeCatalog", API_MakeCatalog)
	r.GET("/api/makeBook", API_MakeBook)
	r.GET("/api/makeChapter", API_MakeChapter)
	r.GET("/api/makeSection", API_MakeSection)
	r.GET("/api/makeObjective", API_MakeObjective)

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
func MakeSectionKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Sections", "", id, nil)
}
func MakeObjectiveKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Objectives", "", id, nil)
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

func getSimpleObjectiveEditor(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Fixed point of edit. we will preform requests for this data later.
	// ctx := appengine.NewContext(req)

	ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
	if numErr != nil || ObjectiveID == 0 {
		http.Redirect(res, req, "/select?status=invalid_id", http.StatusTemporaryRedirect)
	}
	ctx := appengine.NewContext(req)

	objKey := MakeObjectiveKey(ctx, int64(ObjectiveID))
	obj_temp := Objective{}
	objectiveGetErr := datastore.Get(ctx, objKey, &obj_temp)
	HandleError(res, objectiveGetErr)

	ve := VIEW_Editor{}
	ve.ObjectiveID = objKey.IntID()
	ve.ObjectiveTitle = obj_temp.Title
	ve.ObjectiveVersion = obj_temp.Version
	ve.Content = obj_temp.Content
	ve.KeyTakeaways = obj_temp.KeyTakeaways
	ve.SectionID = obj_temp.SectionID

	ServeTemplateWithParams(res, req, "addData.html", ve)
}

// ***************************************************************

// catKey := MakeCatalogKey(ctx, "Fixed-Data-Catalog")
// fixedCatalog := Catalog{"Fixed-Data-Catalog"}
// _, catPutErr := datastore.Put(ctx, catKey, &fixedCatalog)
// HandleError(res, catPutErr)

// bookKey := MakeBookKey(ctx, 0)
// fixedBook := Book{}
// fixedBook.Title = "Fixed-Data-Book"
// fixedBook.CatalogTitle = fixedCatalog.Name
// fixedBookID, bookPutErr := datastore.Put(ctx, bookKey, &fixedBook)
// HandleError(res, bookPutErr)

// chapterKey := MakeChapterKey(ctx, 0)
// fixedChapter := Chapter{}
// fixedChapter.Title = "Fixed-Data-Chapter"
// fixedChapter.BookID = fixedBookID.IntID()
// fixedChapterID, chatperPutErr := datastore.Put(ctx, chapterKey, &fixedChapter)
// HandleError(res, chatperPutErr)

// sectionKey := MakeSectionKey(ctx, 0)
// fixedSection := Section{}
// fixedSection.Title = "Fixed-Data-Section"
// fixedSection.ChapterID = fixedChapterID.IntID()
// fixedSectionID, sectionPutErr := datastore.Put(ctx, sectionKey, &fixedSection)
// HandleError(res, sectionPutErr)

// ObjectiveID, numErr := strconv.Atoi(req.FormValue("ID"))
// var ObjectiveKey datastore.Key
// fixedObjective := Objective{}
// if numErr != nil {
// 	ObjectiveID = 0
// 	objectiveKey = MakeObjectiveKey(ctx, int64(ObjectiveID))
// } else {
// 	objectiveKey = MakeObjectiveKey(ctx, int64(ObjectiveID))
// 	objectiveGetErr := datastore.Get(ctx, objectiveKey, &fixedObjective)
// 	HandleError(res, objectiveGetErr)
// }

// ve := VIEW_Editor{}
// ve.BookID = fixedBookID.IntID()
// ve.BookTitle = fixedBook.Title
// ve.ChapterID = fixedChapterID.IntID()
// ve.ChapterTitle = fixedChapter.Title
// ve.Content = fixedObjective.Content
// ve.KeyTakeaways = fixedObjective.KeyTakeaways
// ve.ObjectiveID = ObjectiveKey.IntID()
// ve.ObjectiveTitle = fixedObjective.Title
// ve.ObjectiveVersion = fixedObjective.Version
// ve.SectionID = fixedSectionID.IntID()
// ve.SectionTitle = fixedSection.Title

// ***************************************************************

// func bookSelectedFromForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	i, numErr := strconv.Atoi(req.FormValue("BookID"))
// 	HandleError(res, numErr)
// 	// bookKey := MakeBookKey(ctx, int64(i))
// 	ServeTemplateWithParams(res, req, "printme.html", i)
// }

// **************************************
// URL Handlers

// func test(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	//var deleteBook string
// 	var x string
// 	var newBook Book
// 	newBook.Title = req.FormValue("BookName")
// 	x = req.FormValue("delete")
// 	if x == "yes" {
// 		//delete book (this is totall insecure!)
// 		var dog Catalog // dog just to show that the get works here
// 		ctx := appengine.NewContext(req)

// 		catKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil)
// 		datastoreErr := datastore.Get(ctx, catKey, &dog)
// 		if datastoreErr != nil {
// 			dog.Name = "NO MESSAGE FOUND - " + datastoreErr.Error()
// 		}

// 		bookKey := datastore.NewKey(ctx, "Books", newBook.Title, 0, catKey)

// 		datastore.Delete(ctx, bookKey) // from there. put the data in the datastore using the key.
// 		/*			if err2 != nil {
// 						http.Error(res, err2.Error(), http.StatusInternalServerError)
// 					}
// 		*/

// 		x = "smack"
// 	}

// 	err := pages.ExecuteTemplate(res, "test.html", x)
// 	if err != nil {
// 		http.Error(res, err.Error(), http.StatusInternalServerError)
// 	}
// }
// func homeAgain(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

// 	var dog Catalog // dog just to show that the get works here

// 	ctx := appengine.NewContext(req)

// 	catKey := datastore.NewKey(ctx, "Catalog", "CatalogID", 0, nil)
// 	datastoreErr := datastore.Get(ctx, catKey, &dog)
// 	if datastoreErr != nil {
// 		dog.Name = "NO MESSAGE FOUND - " + datastoreErr.Error()
// 	}

// 	//var deleteBook string
// 	var newBook Book
// 	newBook.Title = req.FormValue("BookName")
// 	//deleteBook = req.FormValue("delete")

// 	bookKey := datastore.NewKey(ctx, "Books", newBook.Title, 0, catKey)

// 	_, err2 := datastore.Put(ctx, bookKey, &newBook) // from there. put the data in the datastore using the key.
// 	if err2 != nil {
// 		http.Error(res, err2.Error(), http.StatusInternalServerError)
// 	}

// 	q := datastore.NewQuery("Books").Ancestor(catKey)

// 	booklist := make([]Book, 0) // make a list of books. we're filling this out.
// 	for t := q.Run(ctx); ; {    // for values within the query as it's running
// 		var x Book
// 		_, qErr := t.Next(&x)       // read one query value into a temporary location
// 		if qErr == datastore.Done { // if no value was read but it called exit
// 			break // then exit.
// 		} else if qErr != nil { // if there was a real error
// 			http.Error(res, qErr.Error(), http.StatusInternalServerError) // raise that error
// 		}
// 		x.Author = "me"
// 		booklist = append(booklist, x) // add the successful book found onto our output list
// 	}

// 	err := pages.ExecuteTemplate(res, "index.html", booklist)
// 	if err != nil {
// 		http.Error(res, err.Error(), http.StatusInternalServerError)
// 	}
// }
