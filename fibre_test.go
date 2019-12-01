package fibre

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var (
	cwd_arg = flag.String("cwd", "", "set cwd")
)

func init() {
	flag.Parse()
	if *cwd_arg != "" {
		if err := os.Chdir(*cwd_arg); err != nil {
			fmt.Println("Chdir error:", err)
		}
	}
}

func TestNotFoundHandler(t *testing.T) {
	ws := new(WebService)

	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.NotFoundHandler)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("NotFoundHandler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestFavicoHandler(t *testing.T) {
	ws := new(WebService)

	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.FavicoHandler)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("FavicoHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if w.Header().Get("Content-Type") != "image/x-icon" {
		t.Errorf("FavicoHandler returned unexpected Content-Type header: got %v want %v", w.Header().Get("Content-Type"), "image/x-icon")
	}

	expected := "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n"
	if w.Body.String() != expected {
		t.Errorf("FavicoHandler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func TestHomeHandler(t *testing.T) {
	ws := new(WebService)
	ws.Instance = "test"

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.HomeHandler)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("HomeHandler (%v) returned wrong status code: got %v want %v", ws.Instance, status, http.StatusOK)
	}

	// any html response for now is fine.
	body := w.Body.String()
	if !strings.Contains(body, "html") {
		t.Errorf("HomeHandler returned unexpected body: %v", body)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	ws := new(WebService)

	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.HealthCheckHandler)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("TestHealthCheckHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("TestHealthCheckHandler returned unexpected Content-Type header: got %v want %v", w.Header().Get("Content-Type"), "application/json")
	}

	expected := `{"alive": true}`
	if w.Body.String() != expected {
		t.Errorf("TestHealthCheckHandler returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func TestPageHandler(t *testing.T) {
	ws := new(WebService)
	ws.Instance = "test"

	// always have an index in any instance, so;
	req, err := http.NewRequest("GET", "/page/index.html", nil)
	// manually set url vars for the router as so;
	req = mux.SetURLVars(req, map[string]string{"page": "index"})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(ws.PageHandler)
	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("TestPageHandler (%v) returned wrong status code: got %v want %v", ws.Instance, status, http.StatusOK)
	}

	// any html response for now is fine.
	body := w.Body.String()
	if !strings.Contains(body, "html") {
		t.Errorf("TestPageHandler returned unexpected body: %v", body)
	}
}

func TestNewWebService(t *testing.T) {
	ws := new(WebService)
	ws.Instance = "test"
	ws.Address = "127.0.0.1:7999"
}
