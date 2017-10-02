// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample bookshelf is a fully-featured app demonstrating several Google Cloud APIs, including Datastore, Cloud SQL, Cloud Storage.
// See https://cloud.google.com/go/getting-started/tutorial-app
package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	//"io"
	"log"
	"net/http"
	"os"
	//"path"
	"strconv"

	//"cloud.google.com/go/pubsub"
	//"cloud.google.com/go/storage"

	//"golang.org/x/net/context"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/satori/go.uuid"

	"google.golang.org/appengine"

	"github.com/felixwatts/lorefarm"	
	//"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
)

var (
	// See template.go
	pageTmpl   = parseTemplate("page.html")
)

func main() {
	registerHandlers()
	appengine.Main()
}

func registerHandlers() {
	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("/page/5649391675244544", http.StatusFound))

	r.Methods("GET").Path("/page/{id:[0-9]+}").
		Handler(appHandler(pageHandler))

	r.Methods("POST").Path("/new").
		Handler(appHandler(createHandler))

	// The following handlers are defined in auth.go and used in the
	// "Authenticating Users" part of the Getting Started guide.
	r.Methods("GET").Path("/login").
		Handler(appHandler(loginHandler))
	r.Methods("POST").Path("/logout").
		Handler(appHandler(logoutHandler))
	r.Methods("GET").Path("/oauth2callback").
		Handler(appHandler(oauthCallbackHandler))

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	// [START request_logging]
	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	// [END request_logging]
}

// bookFromRequest retrieves a book from the database given a book ID in the
// URL's path.
func pageFromRequest(r *http.Request) (*lorefarm.PageTemplateData, error) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad page id: %v", err)
	}
	page, err := lorefarm.DB.GetPageTemplateData(id)
	if err != nil {
		return nil, fmt.Errorf("could not find page: %v", err)
	}
	return page, nil
}

// detailHandler displays the details of a given book.
func pageHandler(w http.ResponseWriter, r *http.Request) *appError {
	page, err := pageFromRequest(r)
	if err != nil {
		return appErrorf(err, "%v", err)
	}

	return pageTmpl.Execute(w, r, page)
}

// pageFromForm populates the fields of a Page from form values
// (see templates/edit.html).
func pageFromForm(r *http.Request) (*lorefarm.Page, error) {

	i, err := strconv.Atoi(r.FormValue("parentId")) // todo check parent exists
	if(err != nil) {
		return nil, fmt.Errorf("Invalid parent id")
	}

	page := &lorefarm.Page{
	ParentId:  lorefarm.DB.PageId(int64(i)),
	Content:        r.FormValue("content"),
	}

	return page, nil
}

// createHandler adds a book to the database.
func createHandler(w http.ResponseWriter, r *http.Request) *appError {
	page, err := pageFromForm(r)
	if err != nil {
		return appErrorf(err, "could not parse page from form: %v", err)
	}
	id, err := lorefarm.DB.AddPage(page)
	if err != nil {
		return appErrorf(err, "could not save page: %v", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/page/%d", id), http.StatusFound)
	return nil
}

// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
