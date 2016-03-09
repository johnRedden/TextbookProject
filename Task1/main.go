package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var pages *template.Template // This is the storage location for all of our html files

func init() { // init is our equivalent to main in normal files, this will be called first.
	r := httprouter.New()                                                                      // We will use the http router to efficently handle incoming and outgoing requests.
	http.Handle("/", r)                                                                        // Tie this router into all requests from OUR-URL/*
	r.GET("/", home)                                                                           // Handle the default request
	r.GET("/showMessage", showMsg)                                                             // User has requested to see the stored datastore message
	r.GET("/makeMessage", makeMsg)                                                             // User has requested to edit the stored datastore message
	r.POST("/makeMessage", uploadMsg)                                                          // User has submitted a message to send to the datastore
	r.GET("/favicon.ico", favIcon)                                                             // Handle favicon.ico for clients. This is just a reroute to the public file.
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/")))) // Mark all files within the local folder ./public/ to be accessable by OUR-URL/public/filename.ext

	pages = template.Must(pages.ParseGlob("html/*.html")) // Look within ./html/ for any files of extention html and store them within variable pages
}

// **************************************
// URL Handlers

func home(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Home is our reqest handler for any request from OUR-URL/
	err := pages.ExecuteTemplate(res, "index.html", nil) // For this, we will attempt to execute index.html with nil additional information.
	if err != nil {                                      // if for any reason this action throws an error message
		http.Error(res, err.Error(), http.StatusInternalServerError) // post to the user that an error has occured.
	}
}

func showMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Request handler for OUR-URL/showMessage
	// This will get the stored message from datastore and post it in a webpage.

	var messageFromDatastore MessageStructure // Our message storage, You MUST have a struct to store and retrive information from datastore. Simple types will not cover it.

	ctx := appengine.NewContext(req)                                      // Make an appengine context, this will allow us to contact our datastore.
	messageKey := datastore.NewKey(ctx, "Messages", "MessageID", 0, nil)  // make a key for the specific information we're looking for. Looking inside namespace Messages, key of MessageID
	datastoreErr := datastore.Get(ctx, messageKey, &messageFromDatastore) // Attempt to get the information from datastore using context and key, unpackage into var messageFromDatastore
	if datastoreErr != nil {                                              // if this throws an error, override message.data to put this error on the webpage. It's a nicer way than just returning an error in the case that no message has been created yet.
		messageFromDatastore.Data = "NO MESSAGE FOUND - " + datastoreErr.Error()
	}

	err := pages.ExecuteTemplate(res, "showMessage.html", messageFromDatastore.Data) // Now that we have some information, execute our showMessage webpage and send the info into the template.
	if err != nil {                                                                  // Handle all errors. Good Practice.
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

}

func makeMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Request handler for OUR-URL/makeMessage
	// This will post the message editing form to the user, we will handle the info coming back elsewhere.
	err := pages.ExecuteTemplate(res, "makeMessage.html", nil) // Simple post and handle errors
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func uploadMsg(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Request for POST from forms onto OUR-URL/makeMessage
	// Now we're catching the information coming back. We will get the relavant infomation and store it onto the datastore.
	var messageToUpload MessageStructure            // Same as before, we must have a struct to submit data.
	messageToUpload.Data = req.FormValue("Message") // get the form value of Message into the MessageStructure

	// Upload form data to datastore
	ctx := appengine.NewContext(req)                                     // Make the context to know what datastore we're talking to.
	messageKey := datastore.NewKey(ctx, "Messages", "MessageID", 0, nil) // Use the same key as before as we are overriding the previous information not makeing a new entry.
	_, datastoreErr := datastore.Put(ctx, messageKey, &messageToUpload)  // Submit that data to datastore

	if datastoreErr != nil { // if datastore throws an error, handle it.
		http.Error(res, datastoreErr.Error(), http.StatusInternalServerError)
	}

	http.Redirect(res, req, "/", http.StatusFound) // redirect to home, the user can decide what to do from there.
}

func favIcon(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Simple redirect to the relavant public file for our icon. This is only for browsers ease of access.
	http.Redirect(res, req, "public/images/favicon.ico", http.StatusTemporaryRedirect)
}

// *******************************
// Structures

type MessageStructure struct { // A very simple structure to submit a string variable to datastore. This could become more complicated as need demands.
	Data string
}
