package httputil

import (
	"net/http"
)

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"errors": ["bad request"] }`))
}

func Success(w http.ResponseWriter) {
	w.Write([]byte(`{"status": "ok"}`))
}
