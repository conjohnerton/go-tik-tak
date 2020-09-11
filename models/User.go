package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// User holds the data for a user object
type User struct {
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisterDate time.Time `json:"register_date,omitempty"`
}

// ToJSON serializes a User object json using the User model
func (user *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(user)
}

// FromJSON deserializes a User object, referencing the User model
func (user *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(user)
}

// ReadUser creates a new user from a request body if they exist
func ReadUser(r io.Reader) (*User, error) {
	user := &User{}
	err := user.FromJSON(r)
	if err != nil {
		return nil, errors.New("Could not decode a user from the JSON request")
	}

	return user, nil
}
