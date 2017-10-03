package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"

	"github.com/felixwatts/lorefarm"
)

var (
	pageTmpl = parseTemplate("page.html")
)

func main() {
	registerHandlers()
	appengine.Main()
}

func registerHandlers() {
	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := mux.NewRouter()

	// root redirects to first page
	r.Handle("/", http.RedirectHandler(fmt.Sprintf("/page/%d", lorefarm.ROOT_PAGE_ID), http.StatusFound))

	r.Methods("GET").Path("/page/{id:[0-9]+}").
		Handler(appHandler(pageHandler))

	r.Methods("POST").Path("/new").
		Handler(appHandler(createHandler))

	// The following handlers are defined in auth.go
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

	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}

// pageFromRequest retrieves a page from the database given a page Id in the
// URL's path.
func pageFromRequest(r *http.Request) (*lorefarm.Page, error) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad page id: %v", err)
	}
	if(id == 0) {
		id = lorefarm.ROOT_PAGE_ID
	}
	page, err := lorefarm.DB.GetPage(id)
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

// pageFromForm populates the fields of a PageData from form values
func pageFromForm(r *http.Request) (*lorefarm.PageData, error) {

	i, err := strconv.Atoi(r.FormValue("parentId")) // todo check parent exists
	if(err != nil) {
		return nil, fmt.Errorf("Invalid parent id")
	}

	page := &lorefarm.PageData{
	ParentId:  int64(i),
	Content:        r.FormValue("content"),
	}

	return page, nil
}

// createHandler adds a book to the database based on POST values
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
