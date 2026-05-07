package server

import (
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	"github.com/m0nsterrr/as212510.net/v3/internal/config"
	"github.com/m0nsterrr/as212510.net/v3/internal/mikrotik"
)

var (
	//go:embed all:templates/*.html
	tempFs embed.FS

	//go:embed static
	staticFiles embed.FS

	tmpl *template.Template
)

func viewHandler(config config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		router, err := mikrotik.NewRouter(config)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer router.Close()

		if err := router.Information(); err != nil {
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", router); err != nil {
			slog.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func StartServer(config config.Config) {
	staticFS := http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	tmpl = template.Must(template.ParseFS(tempFs, "templates/*.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler(config))
	mux.Handle("/static/", fs)

	slog.Info("Server is starting", "address", config.Server.Address)
	if err := http.ListenAndServe(config.Server.Address, mux); err != nil {
		log.Fatal(err.Error())
	}
}
