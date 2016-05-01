package main

/*
Globals.go by Allen J. Mills
    mm.d.yy

    This file is meant to be used as an exchangable .settings
    file for datastore. Any variables expected to be regularly
    changed out could/should be placed here.
*/

import ()

const GCS_BucketID = "edueditorimages"

var Allowed_Filetypes = []string{"png", "jpg", "jpeg", "gif"} // Ensure that all types here are also supported as image/jpeg for MIME types.

// Place any other variables that could regularly be changed out here.
