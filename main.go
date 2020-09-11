package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/conjohnerton/go-tik-tak/routes"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

var auth *jwtauth.JWTAuth
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)

	key, found := os.LookupEnv("JWT_SECRET_KEY")
	if !found {
		log.Fatal("Could not find argument: JWT_SECRET_KEY")
	}

	auth = jwtauth.New("HS256", []byte(key), nil)
}

func main() {

	port := ":8080"
	server := http.Server{
		Addr:         port,                         // configure the bind address
		Handler:      newRouter(logger, nil, auth), // set the default handler
		ErrorLog:     logger,                       // set the logger for the server
		ReadTimeout:  5 * time.Second,              // max time to read request from the client
		WriteTimeout: 10 * time.Second,             // max time to write response to the client
		IdleTimeout:  120 * time.Second,            // max time for connections using TCP Keep-Alive
	}

	go startServer(&server)
	setServerTermination(&server)
}

func newRouter(log *log.Logger, db *sql.DB, auth *jwtauth.JWTAuth) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// The routes with auth middlewares
	r.Group(routesWithAuth(log))
	r.Mount("/api/users", routes.NewUserHandler(log, db, auth).Routes())

	return r
}

func routesWithAuth(l *log.Logger) func(r chi.Router) {
	return func(r chi.Router) {
		r.Mount("/api/yaks", nil)
	}
}

func startServer(server *http.Server) {
	log.Println("Starting server on port", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func setServerTermination(server *http.Server) {
	// Set termination signals to be sent to c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a termination signal is recieved
	sig := <-c
	log.Println("Got signal:", sig)

	// Try to finish up any currently running tasks before shutting down the server... Very graceful :3
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
