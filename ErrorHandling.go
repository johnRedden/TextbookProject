package main

import (
	"errors"
	"net/http"
)

var (
	ErrNotImplemented = errors.New("Struture Error: Function Not Implemented!")
)

// Internal Function
// generic error handling for any error we encounter.
func HandleError(res http.ResponseWriter, e error) {
	if e != nil {
		http.Error(res, e.Error(), http.StatusInternalServerError)
	}
}

// Internal Function: ErrorPage
/// Prints an error page and returns a boolean representation of the function executing.
/// Results: Boolean Value
////  True: Parent should cease execution, error has been found.
////  False: No Error, Parent may ignore this function.
func ErrorPage(res http.ResponseWriter, ErrorTitle string, e error) bool {
	if e != nil {
		serveErr := pages.ExecuteTemplate(res, "error.gohtml", struct {
			Title   string
			Details error
		}{ErrorTitle, e}) // Execute the error page with the anonymous struct.

		if serveErr != nil {
			panic("Func: ErrorPage\n" + serveErr.Error())
		}
		return true
	}
	return false
}
