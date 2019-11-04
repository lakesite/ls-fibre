# ls-fibre #

🌐 generic default web service with basic handlers 🌐

## usage ##

fibre provides a generic set of handlers for favicon, pages, not found, health
checks, etc as a wrapper around gorilla/mux:

```
  r.HandleFunc("/favicon.ico", ws.FavicoHandler)
  r.HandleFunc("/", ws.HomeHandler)
  r.HandleFunc("/healthcheck", ws.HealthCheckHandler)
  r.HandleFunc("/page/{page}.html", ws.PageHandler)
```

To use, make sure your app has a bin folder layout such as:

```
/
  main.go
  /bin
    main
    /main
      /page
        index.html
      /static
      /templates
        base.html
```
  $ go build -o bin/main main.go
  $ ./bin/main

### index.html ###

```
{{define "content"}}
<p>main service.</p>
{{end}}

```

### base.html ###

```
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
  <link href="data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=" rel="icon" type="image/x-icon">
</head>
<body>
  {{template "content" .}}
</body>
</html>
{{end}}
```

### main.go ###

```
package main

import (
	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

func main() {
	address := config.Getenv("MAIN_HOST", "127.0.0.1") + ":" + config.Getenv("MAIN_PORT", "8080")
	ws := service.NewWebService("main", address)
	ws.RunWebServer()
}

```

## license ##

MIT
