package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/conjohnerton/go-tik-tak/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

// UserHandler is a handler that deals with user login and signup
type UserHandler struct {
	log  *log.Logger
	db   *sql.DB
	auth *jwtauth.JWTAuth
}

// NewUserHandler returns a pointer to a userAPI struct
func NewUserHandler(log *log.Logger, db *sql.DB, auth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{log: log, db: db, auth: auth}
}

// Routes exposes the route methods for the UserHandler
func (api UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/signup", api.SignUp)
	r.Post("/login", api.Login)
	return r
}

// SignUp signs a user up and responds with a jwt token
func (api UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user, err := models.ReadUser(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	print(user)
	return
}

// Login ensures a user exists and responds with an auth token
func (api UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := models.ReadUser(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if api.validateLogin(user) {
		api.log.Println("Creating auth token for user:", user.Email)

		tokenString, err := api.createToken(user)
		if err != nil {
			http.Error(w, "Could not create authentication token", http.StatusForbidden)
		}

		// You can add this struct and encode it to the writer or just construct the json manually
		// type token struct {
		// 	Token string `json:"token"`
		// }
		// json.NewEncoder(w).Encode(&token{Token: tokenString})

		w.Write([]byte("{\"token\":\"" + tokenString + "\"}\n"))
		return
	}

	http.Error(w, "Invalid login details", http.StatusForbidden)
}

func (api UserHandler) validateLogin(user *models.User) bool {
	if user.Email == "john" && user.Password == "pass" {
		return true
	}

	return false
}

func (api UserHandler) createToken(user *models.User) (string, error) {
	_, tokenString, err := api.auth.Encode(jwt.MapClaims{"user_email": user.Email})
	return tokenString, err
}
