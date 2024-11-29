package web

import (
	"net/http"
)

type HellotHandler struct {
}

func NewHelloHandler() *HellotHandler {
	return &HellotHandler{}
}

func (h *HellotHandler) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
