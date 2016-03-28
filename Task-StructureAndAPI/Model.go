package main

/*
Model.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
// "html/template"
)

// Catalog is the root structure, Everything below this will inherit from a Catalog.
type Catalog struct {
	Name    string
	Version float32 `datastore:,noindex`
	Company string
	// Company-Website string
}
type Book struct { // Book has an ancestor in catalog, searchable based on catalog that it was a part of.
	Title        string
	Version      float32  `datastore:,noindex` // we will not query on versions. Do not need to store in a searchable way.
	Author       string   // or array of strings
	Tags         []string // searchable tags to describe the book
	CatalogTitle string   // This is the key.string for Catalog
	// ESBN-10 string
	// ESBN-13 string
	// Copyright date
}
type Chapter struct { // Chapter has an ancestor in Book. Chapter only has meaning from book.
	Title   string
	Version float32 `datastore:,noindex`
	BookID  int64   // key.intID for Book
	// OrderNumber int
}
type Section struct {
	Title     string
	Version   float32 `datastore:,noindex`
	ChapterID int64   // key.intID for Chapter
	// OrderNumber int
	// Text string `datastore:,noindex`
}
type Objective struct {
	Title     string
	Version   float32 `datastore:,noindex`
	SectionID int64   // key.intID for Section
	// Author       string  //or array of strings // doesnt make sense to have this here. the book knows it's author.
	Content      string `datastore:,noindex`
	KeyTakeaways string `datastore:,noindex` // or array of strings
	// Rating       int    `datastore:,noindex` // out of 5 stars
}

type VIEW_Editor struct {
	// This is the key daddy. This is the editor page for the book.
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

	ObjectiveVersion float32
	Content          string
	KeyTakeaways     string
}
