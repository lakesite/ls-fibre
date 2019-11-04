// service includes generic handlers, net/http and mux code for instances of
// servers with API endpoints further defined within their respective packages.
package service

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type WebService struct {
	Router *mux.Router

	Instance string
	Address  string
}

// NotFoundHandler provides a default not found handler for the instance.
func (ws *WebService) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `404 page not found`)
}

func (ws *WebService) FavicoHandler(w http.ResponseWriter, r *http.Request) {
	// blank favico default handler.
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	io.WriteString(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n")
}

// Home handler provides a default index handler for the instance.
func (ws *WebService) HomeHandler(w http.ResponseWriter, r *http.Request) {
	templateLocation := "web/" + ws.Instance + "/page/index.html"
	baseTemplateLocation := "web/" + ws.Instance + "/templates/base.html"
	tmpl, err := template.ParseFiles(templateLocation, baseTemplateLocation)
	if err != nil {
		ws.NotFoundHandler(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(w, "base", struct{ Data string }{Data: "data"})
	}
}

// HealthCheckHandler provides a default health check response (in JSON) for the
// instance.
func (ws *WebService) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}

// Generic handler for /page/<page>.html requests, which reads from the
// root/web/<instance>/templates/<page>.html template.
func (ws *WebService) PageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateLocation := "web/" + ws.Instance + "/page/" + vars["page"] + ".html"
	baseTemplateLocation := "web/" + ws.Instance + "/templates/base.html"
	tmpl, err := template.ParseFiles(templateLocation, baseTemplateLocation)
	if err != nil {
		ws.NotFoundHandler(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(w, "base", struct{ Data string }{Data: "data"})
	}
}

// Create a web service with appropriate handlers.
// instance is a key that will be used in loading templates, static files, etc.
// address is the host and port to listen on
func NewWebService(instance string, address string) *WebService {
	r := mux.NewRouter()
	ws := &WebService{
		Instance: instance,
		Address:  address,
		Router:   r,
	}

	r.NotFoundHandler = http.HandlerFunc(ws.NotFoundHandler)
	r.HandleFunc("/favicon.ico", ws.FavicoHandler)
	r.HandleFunc("/", ws.HomeHandler)
	r.HandleFunc("/healthcheck", ws.HealthCheckHandler)
	r.HandleFunc("/page/{page}.html", ws.PageHandler)

	return ws
}

// Creates a new net/http service with a WebService configuration,
// then run the http.Server
func (ws *WebService) RunWebServer() {
	server := &http.Server{
		Handler:      ws.Router,
		Addr:         ws.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("%v serving on: %v.\n", ws.Instance, ws.Address)
	log.Fatal(server.ListenAndServe())
}
