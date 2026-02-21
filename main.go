// Package main provides a web utility to identify the public exit IP address of the server.
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

//go:embed index.html.tmpl
var tmplFS embed.FS

func setupRoutes(tmpl *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		ipInfo, err := getIpAddressInformation(r.Context(), false)
		if err != nil {
			slog.Error("failed to get IP information", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(ipInfo); err != nil {
				slog.Error("failed to encode JSON", "error", err)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, ipInfo)
		if err != nil {
			slog.Error("failed to execute template", "error", err)
		}
	})
	return mux
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	tmpl, err := template.ParseFS(tmplFS, "index.html.tmpl")
	if err != nil {
		slog.Error("failed to parse template", "error", err)
		os.Exit(1)
	}

	mux := setupRoutes(tmpl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
