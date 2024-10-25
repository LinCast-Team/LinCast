package handlers

import (
	"lincast/api/dtos/requests"
	"lincast/api/dtos/responses"
	"lincast/models"
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

// Función para hashear la contraseña
func HashPassword(password string) (string, error) {
	// Genera el hash usando bcrypt con un costo de 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Función para verificar la contraseña
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (m *Manager) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	data := &requests.SignUpDto{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	// TODO Generate the JWT and register the user
	user := &models.User{
		Username: data.Username,
		Email:    data.Email,
		Name:     data.Name,
		PasswordHash: hashedPassword,
	}

	// m.userRepository.
}
