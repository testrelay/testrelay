package httputil

import (
	"net/http"
)

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"errors": ["unauthorized"] }`))
}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"errors": ["bad request"] }`))
}

func Success(w http.ResponseWriter) {
	w.Write([]byte(`{"status": "ok"}`))
}
