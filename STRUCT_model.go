package main

/*
Model.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"html/template"
)

type Catalog struct {
	Title       string
	Version     float64 `datastore:",noindex"`
	Company     string
	Description template.HTML `datastore:",noindex"`

	ID int64 `datastore:"-"`
}

func (c *Catalog) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Catalogs", "", id.(int64), nil)
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

func (b *Book) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Books", "", id.(int64), nil)
}

type Chapter struct { // Chapter has an ancestor in Book. Chapter only has meaning from book.
	Title       string
	Version     float64       `datastore:",noindex"`
	Description template.HTML `datastore:",noindex"`
	Order       int

	Parent int64 // key.intID for Book
	ID     int64 `datastore:"-"` // self.ID assigned when pulled from datastore.
}

func (c *Chapter) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Chapters", "", id.(int64), nil)
}

type Section struct {
	Title       string
	Version     float64       `datastore:",noindex"`
	Description template.HTML `datastore:",noindex"`
	Order       int

	Parent int64 // key.intID for Chapter
	ID     int64 `datastore:"-"`
}

func (s *Section) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Sections", "", id.(int64), nil)
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

func (o *Objective) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Objectives", "", id.(int64), nil)
}

type Exercise struct {
	Instruction string
	Question    template.HTML `datastore:",noindex"`
	Solution    template.HTML `datastore:",noindex"`
	Order       int

	Parent int64
	ID     int64 `datastore:"-"`
}

func (e *Exercise) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Exercises", "", id.(int64), nil)
}
