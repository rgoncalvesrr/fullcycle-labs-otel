package infra

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebServer struct {
	Router   chi.Router
	Handlers map[string]http.HandlerFunc
	Port     string
}

func NewWebServer(listenPort string) *WebServer {
	return &WebServer{
		Router:   chi.NewRouter(),
		Handlers: make(map[string]http.HandlerFunc),
		Port:     listenPort,
	}
}

func (w *WebServer) RegisterHandler(route string, handler http.HandlerFunc) {
	w.Handlers[route] = handler
}

func (w *WebServer) Start() {
	w.Router.Use(middleware.Logger)
	w.Router.Use(middleware.RequestID)
	w.Router.Use(middleware.RealIP)
	w.Router.Use(middleware.Recoverer)
	w.Router.Use(middleware.Timeout(60 * time.Second))
	w.Router.Handle("/metrics", promhttp.Handler())

	for path, handler := range w.Handlers {
		sr := strings.Split(path, " ")
		if len(sr) == 2 {
			w.Router.Method(sr[0], sr[1], handler)
		} else {
			w.Router.Handle(sr[0], handler)
		}
	}

	err := http.ListenAndServe(":"+w.Port, w.Router)

	if err != nil {
		log.Fatal(err)
	}
}
