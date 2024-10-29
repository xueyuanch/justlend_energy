package http

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/acme/autocert"
	"justlend/internal/config"
	"justlend/internal/justlend"
	_ "justlend/internal/justlend/docs"
	"justlend/internal/log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"
)

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/energy) do not
// need to reference the "net/http" package at all.
type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router

	// Keys used for secure cookie encryption.
	hashKey  []byte
	blockKey []byte

	// Justlend service used by the various HTTP routes.
	service justlend.Service

	// Server option functions for requests.
	opts []kithttp.ServerOption

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocrat.
	addr   string
	domain string

	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	ShutdownTimeout time.Duration
}

// NewServer returns a new instance of Server.
func NewServer(service justlend.Service, c *config.Config) *Server {
	// Copy configuration settings to the new HTTP server and wraps
	// the net/http server & add a gorilla router.
	s := &Server{
		addr:    c.Addr,
		domain:  c.Domain,
		service: service,
		server: &http.Server{
			// Set timeouts to avoid Slow-loris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		router:          mux.NewRouter(),
		hashKey:         c.SCHashKey,
		blockKey:        c.SCBlockKey,
		ShutdownTimeout: c.GracefulTimeout,
		opts: []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.ErrorHandlerFunc(func(_ context.Context, err error) {
				log.Errorf("%v", err)
			})),
			kithttp.ServerErrorEncoder(encodeError),
		},
	}
	return s
}

// RegisterRoutes registers server middlewares & routes using
// the server mux router.
func (s *Server) RegisterRoutes() {
	s.router.Use(s.catchPanic)
	s.router.Use(s.timeout)

	// Setup endpoints to display deployed version.
	if !config.IsReleaseMode() {
		s.router.PathPrefix("/api-docs/").HandlerFunc(httpSwagger.Handler(
			httpSwagger.URL("justlend.org/ad/doc.json"), //The url pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods("GET")
	}

	// Set up a base router that excludes debug handling.
	router := s.router.PathPrefix("/").Subrouter()
	// Extract info's from request header and fill into request context.
	// Register unauthenticated routes.
	{
		r := router.PathPrefix("/").Subrouter()
		s.registerFeeRatioRouters(r)
		s.registerRentResourceRouters(r)
		s.registerReturnResourceRouters(r)
	}

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overriding methods.
	s.server.Handler = s.corsOptions(http.HandlerFunc(s.serveHTTP))
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Override method for forms passing "_method" value.
	if r.Method == http.MethodPost {
		switch v := r.PostFormValue("_method"); v {
		case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
			r.Method = v
		}
	}

	// Override content-type for certain extensions.
	// This allows us to easily cURL API endpoints with a ".json" or ".csv"
	// extension instead of having to explicitly set Content-type & Accept headers.
	// The extensions are removed, so they don't appear in the routes.
	switch ext := path.Ext(r.URL.Path); ext {
	case ".json":
		r.Header.Set("Accept", "application/json")
		r.Header.Set("Content-type", "application/json")
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)
	}

	// Delegate remaining HTTP handling to the gorilla router.
	s.router.ServeHTTP(w, r)
}

// UseTLS returns true if the cert & key file are specified.
func (s *Server) UseTLS() bool {
	return s.domain != ""
}

// Scheme returns the URL scheme for the server.
func (s *Server) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.domain != "" {
		domain = s.domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {

	// Open a listener on our bind address.
	if s.domain != "" {
		s.ln = autocert.NewListener(s.domain)
	} else {
		if s.ln, err = net.Listen("tcp", s.addr); err != nil {
			return err
		}
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	return s.server.Serve(s.ln)
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	return s.server.Shutdown(ctx)
}

// ListenAndServeTLSRedirect runs an HTTP server on port 80 to redirect users
// to the TLS-enabled port 443 server.
func ListenAndServeTLSRedirect(domain string) error {
	return http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domain, http.StatusFound)
	}))
}

// ListenAndServeDebug runs an HTTP server with /debug endpoints (e.g. pprof, vars).
func ListenAndServeDebug(addr string) error {
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, h)
}
