package handler

import (
	"net/http"
)

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}
