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

// ------------------------------------
// Form/Frame Handlers
/////

func IMAGE_PostUploadForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /image/uploader

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
		// User Must be at least Writer.
		fmt.Fprint(res, `{"result":"failure","reason":"Invalid Authorization: `+permErr.Error()+`","code":418}`)
		return
	}
	// ACTION: Give the user an internal permisions key?

	ServeTemplateWithParams(res, req, "simpleImageUploader.html", req.FormValue("oid"))
}

func IMAGE_BrowserForm(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /image/browser
	ctx := appengine.NewContext(req)

	// ACTION: Give the user an internal permisions key?

	prefixQuery := storage.Query{}
	if req.FormValue("oid") == "" {
		prefixQuery.Prefix = "global"
	} else {
		prefixQuery.Prefix = req.FormValue("oid")
	}

	imgl, _ := getFileFromGCS(ctx, &prefixQuery) // get a list of files out of the CS

	imageBrowser := struct { // make a stuct on the fly for the page
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
func IMAGE_API_CKEDITOR_PlaceImageIntoCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// POST: url/api/ckeditor/create
	// Settings: if oid is set, will create image with bucket of oid, otherwise default to global

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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

	// image successfuly sent, let CK know the final url.
	fmt.Fprint(res, `<!DOCTYPE html><html><body><script type="text/javascript">window.parent.CKEDITOR.tools.callFunction('`+req.FormValue("CKEditorFuncNum")+`', "`+"/image?id="+fileName+`","");//window.close();</script></body></html>`)
	return
}

// ------------------------------------
// API - Post/Return/Delete Image to Cloud Storage
/////

func IMAGE_API_PlaceImageIntoCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// POST: /api/makeImage
	// Settings: if oid is set, will create image with bucket of oid, otherwise default to global
	// this is the normal part of the image upload. --not tied to ckeditor

	if validPerm, permErr := HasPermission(res, req, WritePermissions); !validPerm {
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

func IMAGE_API_GetImageFromCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /api/getImage
	// GET: /image
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

func IMAGE_API_RemoveImageFromCS(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// GET: /api/removeImage
	id := req.FormValue("id") // this is an id request only.
	if id == "" {             // if no id, exit with failure.
		fmt.Fprint(res, `{"result":"failure","reason":"missing image id","code":400}`)
		return
	}

	if validPerm, permErr := HasPermission(res, req, AdminPermissions); !validPerm {
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

func IMAGE_API_SendToCloudStorage(req *http.Request, mpf multipart.File, hdr *multipart.FileHeader, prefix string) (string, error) {
	ext, extErr := filterExtension(req, hdr) // ensure that file's extention is an image
	if extErr != nil {                       // if it is not, exit, returning error
		return "", extErr
	}

	uploadName := prefix + makeSHA(mpf) + "." + ext // build new filename based on the image data instead. this will keep us from making multiple files of the same data.
	mpf.Seek(0, 0)                                  // makeSHA moved the reader, move it back.

	ctx := appengine.NewContext(req)
	return uploadName, addFileToGCS(ctx, uploadName, mpf) // upload the file and name. if there is an error, our parent will catch it.}
}

func makeSHA(src multipart.File) string { // make a sha of the contents of the file. we do not want duplicate files.
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func filterExtension(req *http.Request, hdr *multipart.FileHeader) (string, error) {
	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:] // parse through the fileheader for it's extention.
	ext = strings.ToLower(ext)                                   // uppercase, lowercase. all the same here.

	for _, allowedExt := range Allowed_Filetypes { // for all allowed filetypes
		if allowedExt == ext { // found it? Excellent!
			return ext, nil
		}
	}
	// It was not a part of the allowed extentions, return an error.
	return ext, fmt.Errorf("Filetype %s is not allowed by server.", ext)
}

// ------------------------------------
// API - Internal Cloud Storage functions
// Local Only!
/////

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

func removeFileFromGCS(ctx context.Context, filename string) error {
	client, clientErr := storage.NewClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer client.Close()
	return client.Bucket(GCS_BucketID).Object(filename).Delete(ctx)
}

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
