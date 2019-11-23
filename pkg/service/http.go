// service includes generic handlers, net/http and mux code for instances of
// servers with API endpoints further defined within their respective packages.
package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type WebService struct {
	Router *mux.Router

	Instance string
	Address  string
	Apikey   string
}

type ProxyOverride struct {
	Match string
	Host  string
	Path  string
}

type ProxyConfig struct {
	Path     string
	Host     string
	Override ProxyOverride
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

// LogMiddleware simply prints request URIs.
func (ws *WebService) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got request URI: %s\n", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// APIKeyMiddleware provides a built in check for api key, for json api services
func (ws *WebService) APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apik := r.Header.Get("api_key")
		if len(apik) == 0 || apik != ws.Apikey {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid api_key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JsonStatusResponse takees a response writer, response string and status,
// and writes the status and encodes the response string.
func (ws *WebService) JsonStatusResponse(w http.ResponseWriter, response string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
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

func (ws *WebService) SetupProxy(config ProxyConfig) http.Handler {
	// referenced https://www.integralist.co.uk/posts/golang-reverse-proxy/#3
	purl, _ := url.Parse(config.Host)

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", purl.Host)
			req.Host = purl.Host
			req.URL.Host = purl.Host
			req.URL.Scheme = purl.Scheme

			if config.Override.Path != "" && config.Override.Match != "" {
				if strings.HasPrefix(req.URL.Path, config.Override.Match) {
					req.URL.Path = trimLeftChars(req.URL.Path, len(config.Override.Match)) + config.Override.Path
				}
			}
		},

		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return proxy
}

func (ws *WebService) Proxy(config []ProxyConfig) {
	for _, pc := range config {
		proxy := ws.SetupProxy(pc)

		ws.Router.HandleFunc(pc.Path, func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
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
