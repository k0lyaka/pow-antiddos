package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/k0lyaka/pow-antiddos/internal/proxy"
	"github.com/k0lyaka/pow-antiddos/internal/redis"
)

var (
	templates *template.Template
)

func main() {
	redis.InitRedis()
	defer redis.Client.Close()

	config.LoadConfig()
	templates = template.Must(template.ParseGlob("templates/*.html"))

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	handler := proxy.ProxyHandlerWithConfig{Config: config.Config, Templates: templates}

	http.HandleFunc("/", handler.ServeHTTP)

	log.Printf("Starting server on %s", config.Config.ListenAddr)
	log.Fatal(http.ListenAndServe(config.Config.ListenAddr, nil))
}
