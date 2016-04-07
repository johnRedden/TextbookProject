package main

/*
Model.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"html/template"
)

// Catalog is the root structure, Everything below this will inherit from a Catalog.
type Catalog struct {
	Name    string
	Version float64 `datastore:,noindex`
	Company string
	ID      string `datastore:"-"`
	// Company-Website string
}
type Book struct { // Book has an ancestor in catalog, searchable based on catalog that it was a part of.
	Title   string
	Version float64 `datastore:,noindex` // we will not query on versions. Do not need to store in a searchable way.
	Author  string  // or array of strings
	Tags    string  // searchable tags to describe the book, We can search based on substring.

	// ESBN-10 string
	// ESBN-13 string
	// Copyright date

	CatalogTitle string // This is the key.string for Catalog
	ID           int64  `datastore:"-"` // self.ID, assigned when pulled from datastore.
}

type Chapter struct { // Chapter has an ancestor in Book. Chapter only has meaning from book.
	Title   string
	Version float64 `datastore:,noindex`
	Parent  int64   // key.intID for Book
	ID      int64   `datastore:"-"` // self.ID assigned when pulled from datastore.
}

type Section struct {
	Title   string
	Version float64 `datastore:,noindex`
	Parent  int64   // key.intID for Chapter
	ID      int64   `datastore:"-"`
}

type Objective struct {
	Title   string
	Version float64 `datastore:,noindex`
	Author  string  //or array of strings

	Content      template.HTML `datastore:,noindex`
	KeyTakeaways template.HTML `datastore:,noindex` // or array of strings

	Parent int64 // key.intID for Section
	ID     int64 `datastore:"-"`
}

type VIEW_Editor struct {
	// This is the key; the editor page for the book.
	// This struct technically contains Objective but for the sake of verbosity
	//  all included datapoints are spelled out here.
	// This struct will be the data submitted to the editor so a template for the editor
	// should access this information through {{.DataPoint}}
	// Any form that will have submission should include hidden fields for any ID.
	// Please submit all form data with the same name as the datapoint.
	BookID      int64
	ChapterID   int64
	SectionID   int64
	ObjectiveID int64

	BookTitle      string
	ChapterTitle   string
	SectionTitle   string
	ObjectiveTitle string

	ObjectiveVersion float64
	Content          template.HTML
	KeyTakeaways     template.HTML
	Author           string
}
