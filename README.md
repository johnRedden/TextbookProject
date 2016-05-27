# Textbook Project
### [EduNetSystems](http://www.edunetsystems.com/)

In this section we've used Google Cloud Storage: 
*Warning*, If you have not installed and setup GCS this project will not deploy/serve for you.




### Instructions for GCloud Storage installation
  1. *Dependency* - Google App Engine

    If you have not already done so, install Google App Engine. [Quickstart Guide](https://cloud.google.com/appengine/docs/go/)
  2. [Install GCS SDK](https://cloud.google.com/sdk/downloads)
  3. Create a project on Google Cloud Platform

    If you have not already made an account on [Google Cloud Console](https://console.cloud.google.com/project) you should do that now.
  4. Create a Cloud Storage Bucket

    Find, within the console, the section for cloud storage and create a bucket.
  5. Install required packages

    ```
    go get -u golang.org/x/oauth2
    go get -u google.golang.org/cloud/storage
    go get -u google.golang.org/appengine/...
    ```
  6. Configure your project's App.yaml

    ```
    application: <Project ID>
    version: 0
    runtime: go
    api_version: go1

    handlers:
    - url: /.*
    script: _go_app
    ```
  7. Congratulations!

You're all installed and ready to deploy your Appengine project including GCloud storage.

**Final Note:**
  Google Cloud Storage will only work when *deployed*. You can not use the storage features offline.

