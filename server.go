package web

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func NewServer(conf config.Config, logger log.Logger) *webServer {
	port := conf.Get("hdb.server.port", config.AsStringPtr("8080"))
	return &webServer{
		port:   *port,
		conf:   conf,
		logger: logger,
	}
}

// Run starts a HTTP server to listen for rendering requests.
func (server *webServer) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	defer server.logger.Flush()

	router := mux.NewRouter()
	router.Use(server.jsonContentTypeMiddleware)
	router.Use(server.logMiddleware)

	router.HandleFunc("/health", server.handleHealthCheckRequest).Methods("GET")

	//router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/public/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/index.htm"
		}
		fs.ServeHTTP(w, r)
	})).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(notFound)

	server.logger.Infof("Listen [%s]", server.port)
	server.logger.Flush()
	server.httpServer = &http.Server{Addr: ":" + server.port, Handler: router}

	endChan := make(chan error, 1)
	go func() {
		endChan <- server.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		server.stopHttpServer()
	case err := <-endChan:
		return err
	}
	return nil
}

// StopHttpServer will try to sop running HTTP server graceful. Timeout is 3s.
func (server *webServer) stopHttpServer() {
	server.logger.Info("Stopping HTTP server.")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("Unable to stop HTTP server, reason: ", err)
	}
}

// JsonContentTypeMiddleware adds JSON content-type header
func (server *webServer) jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// LogMiddleware adds a logger for all requests. Used log level if debug.
func (server *webServer) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.logger.Debugf("Method: %s, URL: %+v, Header: %+v, URI: %s", r.Method, r.URL, r.Header, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// HandleHealthCheckRequest always returns a 204 status code.
func (server *webServer) handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// NotFound write 404 header to response.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
