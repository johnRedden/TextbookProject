package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"crypto/sha1"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/cloud/storage"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// ------------------------------------
// Image Handlers
/////

func IMAGE_PostUploadForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)
	imgl, _ := filesFromCS(ctx, nil)
	ServeTemplateWithParams(res, req, "simpleImageUploader.html", imgl)
}

// func IMAGE_RecieveFormData(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
// 	// Very Temporary arrangement.
// 	// TODO: self.delete()
// 	multipartFile, multipartHeader, fileError := req.FormFile("incomingImage")
// 	HandleError(res, fileError)
// 	defer multipartFile.Close()

// 	fileName, prepareError := IMAGE_API_SendToCloudStorage(req, multipartFile, multipartHeader)
// 	HandleError(res, prepareError)
// 	res.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	io.WriteString(res, `<img src="`+`https://storage.googleapis.com/`+GCS_BucketID+`\`+fileName+`" />`)
// }

// ------------------------------------
// API - Parse/Prepare Image
/////

func IMAGE_API_SendToCloudStorage(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader) (string, error) {
	ext, extErr := filterExtension(req, hdr) // ensure that file's extention is an image
	if extErr != nil {                       // if it is not, exit, returning error
		return "", extErr
	}

	uploadName := makeSHA(mpf) + "." + ext // build new filename based on the image data instead. this will keep us from making multiple files of the same data.
	mpf.Seek(0, 0)                         // makeSHA moved the reader, move it back.

	ctx := appengine.NewContext(req)
	return uploadName, fileToCS(ctx, uploadName, mpf) // upload the file and name. if there is an error, our parent will catch it.
}

func makeSHA(src multipart.File) string {
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func filterExtension(req *http.Request, hdr *multipart.FileHeader) (string, error) {
	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:] // parse through the fileheader for it's extention.

	for _, allowedExt := range Allowed_Filetypes { // for all allowed filetypes
		if allowedExt == ext { // found it? Excellent!
			return ext, nil
		}
	}
	// It was not a part of the allowed extentions, return an error.
	return ext, fmt.Errorf("Filetype %s is not allowed by server.", ext)
}

// ------------------------------------
// API - Post/Return Image to Cloud Storage
/////

func IMAGE_API_PlaceImageIntoCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// POST: url/API/UploadImage
	multipartFile, multipartHeader, fileError := req.FormFile("incomingImage")
	HandleError(res, fileError)
	defer multipartFile.Close()

	fileName, prepareError := IMAGE_API_SendToCloudStorage(req, multipartFile, multipartHeader)
	HandleError(res, prepareError)
	http.Redirect(res, req, "/api/getImage?id="+fileName, http.StatusSeeOther)
	// fmt.Fprint(res, `{"result":"success","reason":"","code":0,"ID":"`+fileName+`"}`)
}
func IMAGE_API_GetImageFromCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := req.FormValue("id")
	if id == "" {
		fmt.Fprint(res, `{"result":"failure","reason":"missing image id","code":400}`)
		return
	}

	ctx := appengine.NewContext(req)
	client, clientErr := storage.NewClient(ctx)
	HandleError(res, clientErr)
	defer client.Close()

	obj := client.Bucket(GCS_BucketID).Object(id)

	// We'll just copy the image data onto the response, letting the browser know that we're sending an image.
	res.Header().Set("Content-Type", "image/jpeg; charset=utf-8")
	rdr, _ := obj.NewReader(ctx)
	io.Copy(res, rdr)
}

// ------------------------------------
// API - Internal Cloud Storage functions
/////

func fileToCS(ctx context.Context, filename string, freader io.Reader) error {
	client, clientErr := storage.NewClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer client.Close()

	csWriter := client.Bucket(GCS_BucketID).Object(filename).NewWriter(ctx)

	// Cloud Storage Writer - Permissions
	csWriter.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}

	csWriter.ContentType = "image/jpeg"
	io.Copy(csWriter, freader)
	return csWriter.Close()
}

func filesFromCS(ctx context.Context, q *storage.Query) ([]string, error) {
	results := make([]string, 0)

	client, clientErr := storage.NewClient(ctx)
	if clientErr != nil {
		return results, clientErr
	}
	defer client.Close()

	objectList, errList := client.Bucket(GCS_BucketID).List(ctx, q)
	if errList != nil {
		return results, errList
	}

	for _, elem := range objectList.Results {
		results = append(results, elem.Name)
	}
	return results, nil
}
