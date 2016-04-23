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
// Keys
/////
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

// ------------------------------
// Datastore Get/Puts
/////

func GetCatalogFromDatastore(req *http.Request, key string) (Catalog, error) {
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
func RemoveCatalogFromDatastore(req *http.Request, catalogKey string) error {
	ctx := appengine.NewContext(req)
	ck := MakeCatalogKey(ctx, catalogKey)
	return datastore.Delete(ctx, ck)
}

func GetBookFromDatastore(req *http.Request, key int64) (Book, error) {
	if key == 0 {
		return Book{}, nil
	}
	ctx := appengine.NewContext(req)

	bookToReturn := Book{}
	bk := MakeBookKey(ctx, key)

	getErr := datastore.Get(ctx, bk, &bookToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
		return bookToReturn, getErr // dont allow the id to be set
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

func GetChapterFromDatastore(req *http.Request, key int64) (Chapter, error) {
	if key == 0 {
		return Chapter{}, nil
	}
	ctx := appengine.NewContext(req)

	chatperToReturn := Chapter{}
	ck := MakeChapterKey(ctx, key)
	getErr := datastore.Get(ctx, ck, &chatperToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
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

func GetSectionFromDatastore(req *http.Request, key int64) (Section, error) {
	if key == 0 {
		return Section{}, nil
	}
	ctx := appengine.NewContext(req)

	sectionToReturn := Section{}
	sk := MakeSectionKey(ctx, key)
	getErr := datastore.Get(ctx, sk, &sectionToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
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

func GetObjectiveFromDatastore(req *http.Request, key int64) (Objective, error) {
	if key == 0 {
		return Objective{}, nil
	}

	ctx := appengine.NewContext(req)

	objectiveToReturn := Objective{}
	ok := MakeObjectiveKey(ctx, key)
	getErr := datastore.Get(ctx, ok, &objectiveToReturn)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
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

// Datastore Simple retrival helper
func Get_Name_ID_From_Parent(ctx context.Context, parentID interface{}, kind string) []struct {
	Title string
	ID    int64
} {
	// function gatherKindGroup to collect Title/Key information for each given kind
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
