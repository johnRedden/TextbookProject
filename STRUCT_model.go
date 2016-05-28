package main

/*
Model.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"html/template"
)

type Catalog struct {
	Title       string
	Version     float64 `datastore:",noindex"`
	Company     string
	Description template.HTML `datastore:",noindex"`

	ID int64 `datastore:"-"`
}

type Book struct { // Book has an ancestor in catalog, searchable based on catalog that it was a part of.
	Title       string
	Version     float64       `datastore:",noindex"` // we will not query on versions. Do not need to store in a searchable way.
	Author      string        // or array of strings
	Tags        string        // searchable tags to describe the book, We can search based on substring.
	Description template.HTML `datastore:",noindex"`

	Parent int64 // This is the key.string for Catalog
	ID     int64 `datastore:"-"` // self.ID, assigned when pulled from datastore.
}

type Chapter struct { // Chapter has an ancestor in Book. Chapter only has meaning from book.
	Title       string
	Version     float64       `datastore:",noindex"`
	Description template.HTML `datastore:",noindex"`
	Order       int

	Parent int64 // key.intID for Book
	ID     int64 `datastore:"-"` // self.ID assigned when pulled from datastore.
}

type Section struct {
	Title       string
	Version     float64       `datastore:",noindex"`
	Description template.HTML `datastore:",noindex"`
	Order       int

	Parent int64 // key.intID for Chapter
	ID     int64 `datastore:"-"`
}

type Objective struct {
	Title   string
	Version float64 `datastore:",noindex"`
	Author  string  //or array of strings
	Order   int

	Content      template.HTML `datastore:",noindex"`
	KeyTakeaways template.HTML `datastore:",noindex"` // or array of strings

	Parent int64 // key.intID for Section
	ID     int64 `datastore:"-"`
}

type Exercise struct {
	Instruction string
	Question    template.HTML `datastore:",noindex"`
	Solution    template.HTML `datastore:",noindex"`
	Order       int

	Parent int64
	ID     int64 `datastore:"-"`
}
