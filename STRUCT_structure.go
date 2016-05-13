package main

/*
structure.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
)

// ------------------------------
// Datastore Keys for structure objects.
//
// Using Tables: Catalogs, Books, Chapters, Sections, and Objectives
// for our structure objects.
/////
func MakeCatalogKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Catalogs", "", id, nil)
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
func MakeExerciseKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Exercises", "", id, nil)
}

// ------------------------------
// Struct:Catalog, Get,Put, and Remove from datastore
/////
func GetCatalogFromDatastore(req *http.Request, key int64) (Catalog, error) {
	if key == 0 {
		return Catalog{}, nil
	}
	ctx := appengine.NewContext(req)

	catalogToReturn := Catalog{}
	ck := MakeCatalogKey(ctx, key)
	getErr := datastore.Get(ctx, ck, &catalogToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	catalogToReturn.ID = key
	return catalogToReturn, getErr
}
func PutCatalogIntoDatastore(req *http.Request, c Catalog) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	ck := MakeCatalogKey(ctx, c.ID)
	rk, putErr := datastore.Put(ctx, ck, &c)
	return rk, putErr
}
func RemoveCatalogFromDatastore(req *http.Request, catalogKey int64) error {
	ctx := appengine.NewContext(req)
	ck := MakeCatalogKey(ctx, catalogKey)
	return datastore.Delete(ctx, ck)
}

// ------------------------------
// Struct:Book, Get,Put, and Remove from datastore
/////
func GetBookFromDatastore(req *http.Request, key int64) (Book, error) {
	if key == 0 {
		return Book{}, nil
	}
	ctx := appengine.NewContext(req)

	bookToReturn := Book{}
	bk := MakeBookKey(ctx, key)

	getErr := datastore.Get(ctx, bk, &bookToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		return Book{}, nil
	}
	bookToReturn.ID = key
	return bookToReturn, getErr
}
func PutBookIntoDatastore(req *http.Request, b Book) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	bk := MakeBookKey(ctx, b.ID)
	rk, putErr := datastore.Put(ctx, bk, &b)
	return rk, putErr
}
func RemoveBookFromDatastore(req *http.Request, bookKey int64) error {
	ctx := appengine.NewContext(req)
	bk := MakeBookKey(ctx, bookKey)
	return datastore.Delete(ctx, bk)
}

// ------------------------------
// Struct:Chapter, Get,Put, and Remove from datastore
/////
func GetChapterFromDatastore(req *http.Request, key int64) (Chapter, error) {
	if key == 0 {
		return Chapter{}, nil
	}
	ctx := appengine.NewContext(req)

	chatperToReturn := Chapter{}
	ck := MakeChapterKey(ctx, key)
	getErr := datastore.Get(ctx, ck, &chatperToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		return Chapter{}, nil
	}
	chatperToReturn.ID = key
	return chatperToReturn, getErr
}
func PutChapterIntoDatastore(req *http.Request, c Chapter) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	ck := MakeChapterKey(ctx, c.ID)
	rk, putErr := datastore.Put(ctx, ck, &c)
	return rk, putErr
}
func RemoveChapterFromDatastore(req *http.Request, chapterKey int64) error {
	ctx := appengine.NewContext(req)
	ck := MakeChapterKey(ctx, chapterKey)
	return datastore.Delete(ctx, ck)
}

// ------------------------------
// Struct:Section, Get,Put, and Remove from datastore
/////
func GetSectionFromDatastore(req *http.Request, key int64) (Section, error) {
	if key == 0 {
		return Section{}, nil
	}
	ctx := appengine.NewContext(req)

	sectionToReturn := Section{}
	sk := MakeSectionKey(ctx, key)
	getErr := datastore.Get(ctx, sk, &sectionToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		return Section{}, nil
	}
	sectionToReturn.ID = key
	return sectionToReturn, getErr
}
func PutSectionIntoDatastore(req *http.Request, s Section) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	sk := MakeSectionKey(ctx, s.ID)
	rk, putErr := datastore.Put(ctx, sk, &s)
	return rk, putErr
}
func RemoveSectionFromDatastore(req *http.Request, sectionKey int64) error {
	ctx := appengine.NewContext(req)
	sk := MakeSectionKey(ctx, sectionKey)
	return datastore.Delete(ctx, sk)
}

// ------------------------------
// Struct:Objective, Get,Put, and Remove from datastore
/////
func GetObjectiveFromDatastore(req *http.Request, key int64) (Objective, error) {
	if key == 0 {
		return Objective{}, nil
	}

	ctx := appengine.NewContext(req)

	objectiveToReturn := Objective{}
	ok := MakeObjectiveKey(ctx, key)
	getErr := datastore.Get(ctx, ok, &objectiveToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		return Objective{}, nil
	}
	objectiveToReturn.ID = key
	return objectiveToReturn, getErr
}
func PutObjectiveIntoDatastore(req *http.Request, o Objective) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	ok := MakeObjectiveKey(ctx, o.ID)
	rk, putErr := datastore.Put(ctx, ok, &o)
	return rk, putErr
}
func RemoveObjectiveFromDatastore(req *http.Request, objectiveKey int64) error {
	ctx := appengine.NewContext(req)
	ok := MakeObjectiveKey(ctx, objectiveKey)
	return datastore.Delete(ctx, ok)
}

// ------------------------------
// Struct:Exercise, Get,Put, and Remove from datastore
/////
func GetExerciseFromDatastore(req *http.Request, key int64) (Exercise, error) {
	if key == 0 { // 0 is a new blank ID. In this context, return a blank struct.
		return Exercise{}, nil
	}
	ctx := appengine.NewContext(req)

	exerciseToReturn := Exercise{}
	ek := MakeExerciseKey(ctx, key)
	getErr := datastore.Get(ctx, ek, &exerciseToReturn)
	exerciseToReturn.ID = key
	return exerciseToReturn, getErr
}
func PutExerciseIntoDatastore(req *http.Request, e Exercise) (*datastore.Key, error) {
	ctx := appengine.NewContext(req)

	ek := MakeExerciseKey(ctx, e.ID)
	rk, putErr := datastore.Put(ctx, ek, &e)
	return rk, putErr
}
func RemoveExerciseFromDatastore(req *http.Request, exerciseKey int64) error {
	ctx := appengine.NewContext(req)
	ek := MakeExerciseKey(ctx, exerciseKey)
	return datastore.Delete(ctx, ek)
}

// ------------------------------
// Additional Functionality
//
// This section includes functions to handle misc
// interactions with the structure outside of
// previously defined behavior
/////

// Internal Function
// Description:
// This function's main purpose is to return a list of structs
// that have Title and ID for any kind in the datastore.
// Objects to be returned must have the parameters
// Parent(interface{}) and Title(string)
func Get_Name_ID_From_Parent(ctx context.Context, parentID interface{}, kind string) []struct {
	Title string
	ID    int64
} {
	// function Get_Name_ID_From_Parent to collect Title/Key information for each given kind
	q := datastore.NewQuery(kind)      // Make a query into the given kind
	q = q.Filter("Parent =", parentID) // Limit to only the parent ID
	q = q.Project("Title")             // return a struct containing only {Title string}

	output_chapters := make([]struct {
		Title string
		ID    int64
	}, 0)
	for t := q.Run(ctx); ; { // standard query run.
		var cName struct{ Title string }
		k, qErr := t.Next(&cName)

		if qErr == datastore.Done {
			break
		} else if qErr != nil {
			return output_chapters
		}

		output_chapters = append(output_chapters, struct {
			Title string
			ID    int64
		}{cName.Title, k.IntID()})
	}
	return output_chapters
}

// Internal Function
// Description:
//
func Get_Child_Key_From_Parent(ctx context.Context, parentID interface{}, kind string) []*datastore.Key {
	// function Get_Name_ID_From_Parent to collect Title/Key information for each given kind
	q := datastore.NewQuery(kind)      // Make a query into the given kind
	q = q.Filter("Parent =", parentID) // Limit to only the parent ID
	q = q.KeysOnly()

	cks, _ := q.GetAll(ctx, nil)
	return cks
}
