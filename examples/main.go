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
