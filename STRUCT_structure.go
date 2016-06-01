package main

/*
structure.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// ------------------------------
// Datastore Keys for structure objects.
//
// Using Tables: Catalogs, Books, Chapters, Sections, and Objectives
// for our structure objects.
/////

func MakeCatalogKey(ctx context.Context, id int64) *datastore.Key {
	return (&Catalog{}).Key(ctx, id)
}
func MakeBookKey(ctx context.Context, id int64) *datastore.Key {
	return (&Book{}).Key(ctx, id)
}
func MakeChapterKey(ctx context.Context, id int64) *datastore.Key {
	return (&Chapter{}).Key(ctx, id)
}
func MakeSectionKey(ctx context.Context, id int64) *datastore.Key {
	return (&Section{}).Key(ctx, id)
}
func MakeObjectiveKey(ctx context.Context, id int64) *datastore.Key {
	return (&Objective{}).Key(ctx, id)
}
func MakeExerciseKey(ctx context.Context, id int64) *datastore.Key {
	return (&Exercise{}).Key(ctx, id)
}

// ------------------------------
// Struct:Catalog, Get from datastore
/////
func GetCatalogFromDatastore(ctx context.Context, key int64) (Catalog, error) {
	if key == 0 {
		return Catalog{}, nil
	}

	c := Catalog{}
	getErr := GetFromDatastore(ctx, key, &c)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	c.ID = key
	return c, getErr
}

// ------------------------------
// Struct:Book, Get from datastore
/////
func GetBookFromDatastore(ctx context.Context, key int64) (Book, error) {
	if key == 0 {
		return Book{}, nil
	}

	b := Book{}
	getErr := GetFromDatastore(ctx, key, &b)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	b.ID = key
	return b, getErr
}

// ------------------------------
// Struct:Chapter, Get,Put, and Remove from datastore
/////
func GetChapterFromDatastore(ctx context.Context, key int64) (Chapter, error) {
	if key == 0 {
		return Chapter{}, nil
	}

	e := Chapter{}
	getErr := GetFromDatastore(ctx, key, &e)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	e.ID = key
	return e, getErr
}

// ------------------------------
// Struct:Section, Get from datastore
/////
func GetSectionFromDatastore(ctx context.Context, key int64) (Section, error) {
	if key == 0 {
		return Section{}, nil
	}

	e := Section{}
	getErr := GetFromDatastore(ctx, key, &e)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	e.ID = key
	return e, getErr
}

// ------------------------------
// Struct:Objective, Get,Put, and Remove from datastore
/////
func GetObjectiveFromDatastore(ctx context.Context, key int64) (Objective, error) {
	if key == 0 {
		return Objective{}, nil
	}

	e := Objective{}
	getErr := GetFromDatastore(ctx, key, &e)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	e.ID = key
	return e, getErr
}

// ------------------------------
// Struct:Exercise, Get,Put, and Remove from datastore
/////
func GetExerciseFromDatastore(ctx context.Context, key int64) (Exercise, error) {
	if key == 0 {
		return Exercise{}, nil
	}

	e := Exercise{}
	getErr := GetFromDatastore(ctx, key, &e)
	if getErr == datastore.ErrNoSuchEntity {
		getErr = nil
	}
	e.ID = key
	return e, getErr
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
