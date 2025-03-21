package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.html.tmpl
var tmplFS embed.FS

func main() {
	tmpl, err := template.ParseFS(tmplFS, "index.html.tmpl")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ipInfo, err := getIpAddressInformation(false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, ipInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", mux)
}
