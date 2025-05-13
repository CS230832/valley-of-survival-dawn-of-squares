package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" || len(path) == 0 {
		http.ServeFile(w, r, "static/index.html")
		return
	}

	if !strings.HasSuffix(path, ".html") {
		htmlPath := filepath.Join("static/", path+".html")
		if file, err := os.Stat(htmlPath); err == nil && !file.IsDir() {
			http.ServeFile(w, r, htmlPath)
			return
		}
	}

	fullPath := filepath.Join("static/", path)
	if file, err := os.Stat(fullPath); err == nil && !file.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
