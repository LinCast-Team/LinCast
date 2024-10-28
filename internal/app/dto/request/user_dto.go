package request

import (
	"errors"
	"net/http"
	"strings"
)

type SignUpDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func (su *SignUpDto) Bind(r *http.Request) error {
	if (su.Email == "") || (su.Password == "") || (su.Username == "") || (su.Name == "") {
		return errors.New("missing required email, password, username, or name fields")
	}

	if strings.Contains(su.Username, " ") {
		return errors.New("username cannot contain spaces")
	}

	su.Email = strings.TrimSpace(strings.ToLower(su.Email))

	return nil
}
