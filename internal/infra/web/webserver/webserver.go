package webserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type methodPatternHandler struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

type WebServer struct {
	Router     chi.Router
	Handlers   []methodPatternHandler
	Server     *http.Server
	middleware func(http.Handler) http.Handler
}

func NewWebServer(serverPort string, middleware func(http.Handler) http.Handler) *WebServer {
	return &WebServer{
		Router:     chi.NewRouter(),
		Handlers:   []methodPatternHandler{},
		Server:     &http.Server{Addr: serverPort},
		middleware: middleware,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc, method string) {
	s.Handlers = append(s.Handlers, methodPatternHandler{
		method:  method,
		pattern: path,
		handler: handler,
	})
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	s.Router.Use(s.middleware)
	for _, handler := range s.Handlers {
		s.Router.MethodFunc(handler.method, handler.pattern, handler.handler)
	}
	s.Server.Handler = s.Router
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *WebServer) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
