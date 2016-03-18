package main

/*
Model.go by Allen J. Mills
    mm.d.yy

    Description
*/

import ()

// Catalog is the root structure, Everything below this will inherit from a Catalog.
type Catalog struct {
	Name    string
	Version float32 `datastore:,noindex`
	Company string
	// Company-Website string
}
type Book struct { // Book has an ancestor in catalog, searchable based on catalog that it was a part of.
	Title   string
	Version float32  `datastore:,noindex` // we will not query on versions. Do not need to store in a searchable way.
	Author  string   // or array of strings
	Tags    []string // searchable tags to describe the book
	// ESBN-10 string
	// ESBN-13 string
	// Copyright date
}
type Chapter struct { // Chapter has an ancestor in Book. Chapter only has meaning from book.
	Title   string
	Version float32 `datastore:,noindex`
}
type Section struct {
	Title   string
	Version float32
	// Text string `datastore:,noindex`
}
type Objective struct {
	Objective string
	Version   float32 `datastore:,noindex`
	// Author       string  //or array of strings // doesnt make sense to have this here. the book knows it's author.
	Content      string `datastore:,noindex`
	KeyTakeaways string // or array of strings
	Rating       int    // out of 5 stars
}

type VIEW_CatalogData struct { // structure to view catalog information
	Catalogs []KeyValuePair
	Books    []KeyValuePair
}

type VIEW_BookData struct { // structure to view book information
	Catalog_Name KeyValuePair
	Book_Name    KeyValuePair
	Sections     KeyValuePair
}

type KeyValuePair struct {
	Key   string
	Value string
}
