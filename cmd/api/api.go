package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)


type application struct {
	config config
}


type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	
	// setup endpoint
	r.Route("/v1", func(r chi.Router){
		r.Get("/health", app.healthCheckHandler)

	})
	return r
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second*30,
		ReadTimeout: time.Second*10,
		IdleTimeout: time.Minute,
	}

	log.Println("server has started at", app.config.addr)
	return server.ListenAndServe()
}