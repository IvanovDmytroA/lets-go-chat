package handler

import (
	"net/http"
)

const indexPage string = "web/static/index.html"

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, indexPage)
}
