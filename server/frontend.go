package server

import (
	"net/http"
)

// HandleFrontend serves the HTML frontend
func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
