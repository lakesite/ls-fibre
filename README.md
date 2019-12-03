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

fibre also provides generic api key middleware and logging middleware (for 
debugging):

```
  ...
	address := config.Getenv("MAIN_HOST", "127.0.0.1") + ":" + config.Getenv("MAIN_PORT", "8080")
	ws := fibre.NewWebService("main", address)

  ws.Apikey = config.Getenv("MAIN_API_KEY", "default_key")

  ws.Router.Use(ws.LogMiddleware)
  ws.Router.Use(ws.APIKeyMiddleware)

```

fibre also provides a simple method for proxying requests:

```
  cfg := []service.ProxyConfig{
    service.ProxyConfig{
      Path: "/",
      Host: "redirecthost.co"
    },
    service.ProxyConfig{
      Path: "/whatever/path",
      Host: "google.com",
      Override: service.ProxyOverride{
        Match: "/api/v2",
        Path: "/api/v3",
      },
    },
  }

  ws.Proxy(cfg)
```

To use pages with templates, make sure your app has a bin folder layout 
matching the service name (in this case, main) such as:

```
/
  main.go
  /bin
    main
    /web
      /main
        /page
          index.html
        /static
        /templates
          base.html
```

        $ go build -o bin/main main.go
        $ cd bin
        $ ./main

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
	"github.com/lakesite/ls-config"
	"github.com/lakesite/ls-fibre"
)

func main() {
	address := config.Getenv("MAIN_HOST", "127.0.0.1") + ":" + config.Getenv("MAIN_PORT", "8080")
	ws := fibre.NewWebService("main", address)
	ws.RunWebServer()
}
```

## testing ##

  $ go test

## running ##

  $ cd examples
  $ go run main.go

  Visit http://localhost:8080/

## license ##

MIT
