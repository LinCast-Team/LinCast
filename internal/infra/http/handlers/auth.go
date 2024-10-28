package handlers

import (
	"encoding/json"
	"lincast/internal/app/dto/request"
	"lincast/internal/app/dto/response"
	"lincast/internal/app/service"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/go-chi/render"
)

type AuthHandlers struct {
	userService *service.UserApplicationService
}

func NewAuthHandlers(userService *service.UserApplicationService) *AuthHandlers {
	return &AuthHandlers{
		userService: userService,
	}
}

func (m *AuthHandlers) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	data := &request.SignUpDto{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	_, jwt, err := m.userService.RegisterUser(data.Username, data.Email, data.Name, data.Password)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	err = json.NewEncoder(w).Encode(jwt)
	if err != nil {
		log.WithError(err).Panicln("failed to encode response")
	}
}
