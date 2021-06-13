package proxy

import (
	"net/http"
	"net/http/httputil"
)

type lbHandler struct {
	hosts map[string]*httputil.ReverseProxy
}

func (h *lbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.hosts[r.Host].ServeHTTP(w, r)
}
