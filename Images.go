package main

/*
Images.go by Allen J. Mills
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

var (
	// Permission Requirements for Image
	image_Make_Permission   = WritePermissions
	image_Delete_Permission = AdminPermissions
)

// ------------------------------------
// Form/Frame Handlers
/////

// Call: /image/uploader
// Description:
// This handler will serve an html form to upload an
// image if you have at minimum make permissions for this module.
// Will pass along an objective id through oid.
//
// Method: GET
// Results: HTML/JSON
// Mandatory Options:
// Optional Options: oid
// Codes:
//      418 : Invalid Authorization; Check your login status and permission level.
func IMAGE_PostUploadForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, image_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}
	// ACTION: Give the user an internal permissions key?

	ServeTemplateWithParams(res, req, "simpleImageUploader.html", req.FormValue("oid"))
}

// Call: /image/browser
// Description:
// This handler will serve an html file browser
// to allow a user to browse images uploaded to the server.
// Option:oid will limit the returned images to only those
// with the objective id.
// Option:CKEditorFuncNum will alert the browser that it is a child
// of a CKEditor instance and should attempt to let the editor know of
// any selection the user makes.
//
// Method: GET
// Results: HTML
// Mandatory Options:
// Optional Options: oid, CKEditorFuncNum
func IMAGE_BrowserForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(req)

	// ACTION: Give the user an internal permissions key?

	prefixQuery := storage.Query{}
	if req.FormValue("oid") == "" {
		prefixQuery.Prefix = "global"
	} else {
		prefixQuery.Prefix = req.FormValue("oid")
	}

	imgl, _ := getFileFromGCS(ctx, &prefixQuery) // get a list of files out of the CS

	imageBrowser := struct { // make a struct on the fly for the page
		CKEditorFuncNum string
		Images          []string
		ID              string
	}{
		req.FormValue("CKEditorFuncNum"),
		imgl,
		req.FormValue("oid"),
	}

	ServeTemplateWithParams(res, req, "simpleImageBrowser.html", imageBrowser)
}

// ------------------------------------
// CKEditor Specific Handlers
/////

// Call: /api/ckeditor/create
// Description:
// This handler will take in an image in
// Mandatory:upload and send it to cloud storage.
//
// Option:oid tie the uploaded image to an objective id.
// Option:CKEditorFuncNum will alert the browser that it is a child
// of a CKEditor instance and should attempt to let the editor know of
// any selection the user makes.
//
// Method: POST
// Results: HTML
// Mandatory Options: upload
// Optional Options: oid, CKEditorFuncNum
func IMAGE_API_CKEDITOR_PlaceImageIntoCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, permErr := HasPermission(res, req, image_Make_Permission); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `<!DOCTYPE html><html><body><script type="text/javascript">window.parent.CKEDITOR.tools.callFunction('`+req.FormValue("CKEditorFuncNum")+`',"","`+permErr.Error()+`");//window.close();</script></body></html>`)
		return
	}

	multipartFile, multipartHeader, fileError := req.FormFile("upload") // pull the uploaded image out of the request
	if fileError != nil {                                               // if there was an issue with the request, exit and note that to ckeditor
		fmt.Fprint(res, `<!DOCTYPE html><html><body><script type="text/javascript">window.parent.CKEDITOR.tools.callFunction('`+req.FormValue("CKEditorFuncNum")+`',"","`+fileError.Error()+`");//window.close();</script></body></html>`)
		return
	}
	defer multipartFile.Close()

	prefix := req.FormValue("oid")
	if prefix == "" {
		prefix = "global"
	}

	fileName, prepareError := IMAGE_API_SendToCloudStorage(req, multipartFile, multipartHeader, prefix) // send the image out to the cloudstore.
	if prepareError != nil {
		fmt.Fprint(res, `<!DOCTYPE html><html><body><script type="text/javascript">window.parent.CKEDITOR.tools.callFunction('`+req.FormValue("CKEditorFuncNum")+`',"","`+prepareError.Error()+`");//window.close();</script></body></html>`)
		return
	}

	// image successfully sent, let CK know the final url.
	fmt.Fprint(res, `<!DOCTYPE html><html><body><script type="text/javascript">window.parent.CKEDITOR.tools.callFunction('`+req.FormValue("CKEditorFuncNum")+`', "`+"/image?id="+fileName+`","");//window.close();</script></body></html>`)
	return
}

// ------------------------------------
// API - Post/Return/Delete Image to Cloud Storage
/////

// Call: /api/create/image
// Description:
// This handler will take in an image in
// Mandatory:upload and send it to cloud storage.
//
// Option:oid tie the uploaded image to an objective id..
//
// Method: POST
// Results: HTTP Redirect
// Mandatory Options: upload
// Optional Options: oid
// Codes:
//      Success, redirect to image/uploader with status of success
//      Failure, redirect to image/uploader with status of failure
func IMAGE_API_PlaceImageIntoCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if validPerm, _ := HasPermission(res, req, image_Make_Permission); !validPerm {
		// User Must be at least Writer.
		http.Redirect(res, req, "/image/uploader?status=failure&message=invalid_login", http.StatusSeeOther)
		return
	}

	multipartFile, multipartHeader, fileError := req.FormFile("upload") // pull uploaded image.
	if fileError != nil {                                               // handle error in a stable way, this will be a part of another page.
		http.Redirect(res, req, "/image/uploader?status=failure", http.StatusSeeOther)
		return
	}
	defer multipartFile.Close()

	prefix := req.FormValue("oid")
	if prefix == "" {
		prefix = "global"
	}

	_, prepareError := IMAGE_API_SendToCloudStorage(req, multipartFile, multipartHeader, prefix)
	if prepareError != nil { // send to CS and same as above.
		http.Redirect(res, req, "/image/uploader?status=failure", http.StatusSeeOther)
		return
	}
	// success, let user know that their image is waiting.
	http.Redirect(res, req, "/image/uploader?status=success", http.StatusSeeOther)
}

// Call: /image
// Description:
// This handler will retrieve an image from cloud storage
//
// Method: GET
// Results: Image Binary/JSON
// Mandatory Options: id
// Optional Options:
// Codes:
//      Success, Image in response
//      400 - Missing Parameter
func IMAGE_API_GetImageFromCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := req.FormValue("id") // this is an id request only.
	if id == "" {             // if no id, exit with failure.
		fmt.Fprint(res, `{"result":"failure","reason":"missing image id","code":400}`)
		return
	}

	ctx := appengine.NewContext(req) // quickly get a handle into CS
	client, clientErr := storage.NewClient(ctx)
	HandleError(res, clientErr)
	defer client.Close()

	obj := client.Bucket(GCS_BucketID).Object(id) // pull the object from cs

	// We'll just copy the image data onto the response, letting the browser know that we're sending an image.
	res.Header().Set("Content-Type", "image/jpeg; charset=utf-8")
	rdr, _ := obj.NewReader(ctx)
	io.Copy(res, rdr)
}

// Call: /api/delete/image
// Description:
// This handler will delete an image from cloud storages
//
// Method: POST
// Results: JSON
// Mandatory Options: id
// Optional Options:
// Codes:
//        0 - Success, All actions completed
//      400 - Failure, Missing Parameter
//      418 - Failure, Invalid Authorization
//      500 - Failure, Internal Services Error
//
func IMAGE_API_RemoveImageFromCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := req.FormValue("id") // this is an id request only.
	if id == "" {             // if no id, exit with failure.
		fmt.Fprint(res, `{"result":"failure","reason":"missing image id","code":400}`)
		return
	}

	if validPerm, permErr := HasPermission(res, req, image_Delete_Permission); !validPerm {
		// User Must be at least Admin.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}

	ctx := appengine.NewContext(req)
	csRemoveErr := removeFileFromGCS(ctx, id)

	if csRemoveErr != nil {
		fmt.Fprint(res, `{"result":"failure","reason":"Internal Error:`+csRemoveErr.Error()+`","code":500}`)
		return
	}
	fmt.Fprint(res, `{"result":"success","reason":"","code":0}`)
}

// ------------------------------------
// API - Parse/Prepare Image
/////

// Internal Function
// Description:
// This function will prepare and send file to GCS then return the key.
//
// Returns:
//      key(string) - name of GCS key. Item is now in GCS
//      failure?(error) - If any errors occur they exist here.
func IMAGE_API_SendToCloudStorage(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader, prefix string) (string, error) {
	ext, extErr := filterExtension(req, hdr) // ensure that file's extension is an image
	if extErr != nil {                       // if it is not, exit, returning error
		return "", extErr
	}

	uploadName := prefix + makeSHA(mpf) + "." + ext // build new filename based on the image data instead. this will keep us from making multiple files of the same data.
	mpf.Seek(0, 0)                                  // makeSHA moved the reader, move it back.

	ctx := appengine.NewContext(req)
	return uploadName, addFileToGCS(ctx, uploadName, mpf) // upload the file and name. if there is an error, our parent will catch it.}
}

// Internal Function
// Description:
// This function will create a SHA name of a file's contents.
// This will ensure that duplicate items have the same key.
//
// Returns:
//      key(string) - SHA of contents.
func makeSHA(src multipart.File) string {
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Internal Function
// Description:
// This function will ensure that the extension of an incoming filename
// is allowed by this server. Will return the isolated extension if yes.
//
// Returns:
//      extension(string) - Extension of file.
//      failure?(error) - Error if filetype is not allowed.
func filterExtension(req *http.Request, hdr *multipart.FileHeader) (string, error) {
	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:] // parse through the fileheader for its extension.
	ext = strings.ToLower(ext)                                   // uppercase, lowercase. all the same here.

	for _, allowedExt := range Allowed_Filetypes { // for all allowed filetypes
		if allowedExt == ext { // found it? Excellent!
			return ext, nil
		}
	}
	// It was not a part of the allowed extensions, return an error.
	return ext, fmt.Errorf("Filetype %s is not allowed by server.", ext)
}

// ------------------------------------
// API - Internal Cloud Storage functions
// Local Only!
/////

// Internal Function
// Description:
// This function will add a file to GCS at filename.
//
// Returns:
//      failure?(error) - Error if storage fails.
func addFileToGCS(ctx context.Context, filename string, freader io.Reader) error {
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

// Internal Function
// Description:
// This function will remove a file to GCS at filename.
//
// Returns:
//      failure?(error) - Error if deletion fails.
func removeFileFromGCS(ctx context.Context, filename string) error {
	client, clientErr := storage.NewClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer client.Close()
	return client.Bucket(GCS_BucketID).Object(filename).Delete(ctx)
}

// Internal Function
// Description:
// This function will retrive filenames from GCS per a storage Query
//
// Returns:
//      files []string - list of filenames.
//      failure?(error) - Error if storage fails.
func getFileFromGCS(ctx context.Context, q *storage.Query) ([]string, error) {
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
