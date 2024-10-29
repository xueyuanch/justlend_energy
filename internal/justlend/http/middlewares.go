package http

import (
	"context"
	"justlend/internal/config"
	"justlend/internal/derrors"
	"justlend/internal/log"
	"net/http"
	"time"
)

const (
	// defaultRequestTimeout defines the maximum time allowed before a request time out.
	defaultRequestTimeout = 300
	// defaultApikeyHeader is the default header name for apikey.
	defaultApikeyHeader = "X-APIKEY"
)

var eKeyEnabled = config.GetEnvBool("ENCRYPT_KEY_ENABLED", true)

// Timeout returns a new middleware that times out each request after the given duration.
func (s *Server) timeout(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			time.Duration(config.GetEnvInt("HTTP_REQUEST_TIMEOUT", defaultRequestTimeout))*time.Second,
		)
		defer cancel()
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// catchPanic is middleware for catching panics and logging them.
func (s *Server) catchPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				encodeError(r.Context(), derrors.Internal, w)
				log.Error(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Defines max length limit for request URI.
const maxURILength = 1000

// CORS strategy options goes below.
func (s *Server) corsOptions(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.String()) >= maxURILength {
			http.Error(w, http.StatusText(http.StatusRequestURITooLong), http.StatusRequestURITooLong)
			return
		}
		// Use request origin for debugging.
		origin := "*"
		if s.domain != "" {
			// Access-Control-Allow-Origin value should be the
			// base url if domain is set.
			origin = s.URL()
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		//w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-SF-Language, Content-Type, Accept, User-Agent, Authorization, X-Requested-With, x-request-passcode")
		//w.Header().Set("Access-Control-Max-Age", "86400")
		//w.Header().Set("X-Content-Type-Options", "nosniff") // Prevent MIME sniffing.
		//w.Header().Set("X-Frame-Options", "deny")           // Don't allow frame embedding.

		// Always set Vary headers
		// see https://github.com/rs/cors/issues/10,
		//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Vary", "Access-Control-Request-Method")
		w.Header().Set("Vary", "Access-Control-Request-Headers")

		// Return status 'OK' if it's a pre-flight OPTIONS request.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
