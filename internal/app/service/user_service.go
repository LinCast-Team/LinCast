package service

import (
	"lincast/internal/app/dto/response"
	"lincast/internal/domain/entities"
	"lincast/internal/domain/repositories"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserApplicationService struct {
	userRepo  repositories.UserRepository
	tokenAuth *jwtauth.JWTAuth
}

func NewUserApplicationService(userRepo repositories.UserRepository, tokenAuth *jwtauth.JWTAuth) *UserApplicationService {
	return &UserApplicationService{
		userRepo:  userRepo,
		tokenAuth: tokenAuth,
	}
}

func (uas *UserApplicationService) RegisterUser(username, email, name, password string) (*entities.User, *response.CompleteJWTDto, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user, err := entities.NewUser(username, hashedPassword, email, name)
	if err != nil {
		return nil, nil, err
	}

	// TODO Verify if the user already exists

	err = uas.userRepo.Create(user)
	if err != nil {
		return nil, nil, err
	}

	tokenClaims := map[string]interface{}{"user_id": user.ID.String()}
	tokenExpiration := time.Now().Add(time.Minute * 15)
	jwtauth.SetExpiry(tokenClaims, tokenExpiration)

	refreshTokenClaims := map[string]interface{}{"user_id": user.ID.String()}
	refreshTokenExpiration := time.Now().Add(time.Hour * 24 * 7)
	jwtauth.SetExpiry(tokenClaims, refreshTokenExpiration)

	_, tokenString, err := uas.tokenAuth.Encode(tokenClaims)
	if err != nil {
		return nil, nil, err
	}

	_, refreshTokenString, err := uas.tokenAuth.Encode(refreshTokenClaims)
	if err != nil {
		return nil, nil, err
	}

	user.RefreshToken = refreshTokenString
	user.RefreshTokenExpiry = refreshTokenExpiration

	err = uas.userRepo.Update(user)
	if err != nil {
		return nil, nil, err
	}

	return user, response.NewCompleteJWTDto(tokenString, refreshTokenString, tokenExpiration), nil
}

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
