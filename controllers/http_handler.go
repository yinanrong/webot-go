package controllers

import "net/http"

type HttpHandler struct {
	Mux map[string]func(http.ResponseWriter, *http.Request)
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := h.Mux[r.URL.Path]; ok {
		h(w, r)
		return
	}
	http.ServeFile(w, r, "../views/notfound.html")
}
