package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// userHandler is a handler that deals with user login and signup
type userHandler struct {
	log *log.Logger
}

// NewuserHandler returns a pointer to a userAPI struct
func NewUserHandler(log *log.Logger, db *sql.DB) *userHandler {
	return &userHandler{log: log}
}

// Routes exposes the route methods for the userHandler
func (api userHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", api.SignUp)
	return r
}

// SignUp signs a user up and responds with a jwt token
func (api userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signing you up."))
}
