package main

import (
	"fmt"
	"github.com/oklog/run"
	"justlend/internal/config"
	"justlend/internal/justlend"
	"justlend/internal/justlend/http"
	"justlend/internal/log"
	"justlend/internal/repos"
	"justlend/internal/tron"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Construct a new type to represent our application daemon.
	// This type lets us shared setup code with our end-to-end tests.
	d := newDaemon()

	d.StartHTTPServer()
	// This function just sits and waits for ctrl-C.
	w := make(chan struct{})
	d.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-w:
			return nil
		}
	}, func(error) {
		close(w)
		// Clean up program.
		if err := d.Close(); err != nil {
			log.ErrorW("fails to clean up program", "error", err)
		}
	})

	log.InfoW("shutting down...", "error", d.Run())
}

// daemon represents the admin daemon which contains all the
// dependencies required for the program to run.
type daemon struct {
	run.Group                   // Embed `run.Group` for running actors.
	Config     *config.Config   // Resolved config data.
	HTTPServer *http.Server     // HTTP server for handling HTTP communication.trxEnergy service is attached to it before running.
	Service    justlend.Service // application service.
	Endpoint   *tron.Endpoint
}

func newDaemon() *daemon {
	d := &daemon{}
	var err error
	d.Config = config.Resolve()

	//d.DB = repos.NewDB()
	//if err = d.DB.Open(d.Config.DBConnInfo()); err != nil {
	//	log.FatalW("cannot open db", "error", err)
	//}

	if d.Endpoint, err = tron.NewEndpoint(); err != nil {
		log.FatalW("cannot connect tron", "error", err)
	}

	d.Service = repos.NewService(d.Endpoint)

	return d
}

func (d *daemon) StartHTTPServer() {
	// Construct HTTP server.
	d.HTTPServer = http.NewServer(d.Service, d.Config)

	// Register server middlewares & routes before open.
	d.HTTPServer.RegisterRoutes()

	// Start the HTTP server.
	d.Add(func() error {
		log.InfoW("Running justlend HTTP server", "transport", "HTTP", "addr", d.Config.Addr)
		return d.HTTPServer.Open()
	}, func(err error) {
		d.HTTPServer.Close()
	})

	if d.HTTPServer.UseTLS() {
		// If TLS enabled, redirect non-TLS connections to TLS.
		d.Add(func() error {
			return http.ListenAndServeTLSRedirect(d.Config.Domain)
		}, func(error) {})
	}
}

func (d *daemon) Close() error {
	if d.HTTPServer != nil {
		if err := d.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if d.Endpoint != nil {
		if err := d.Endpoint.Close(); err != nil {
			return err
		}
	}
	return nil
}
