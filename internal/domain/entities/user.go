package entities

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID         `json:"id" gorm:"type:char(36);primary_key"`
	Username        string            `json:"username" gorm:"unique"`
	PasswordHash    string            `json:"-"`
	Email           string            `json:"email" gorm:"unique"`
	Name            string            `json:"name"`
	PlayerID        uuid.UUID         `json:"playerID"`
	Player          PlaybackInfo      `json:"player"`
	Queue           []QueueEpisode    `json:"queue"`
	EpisodeProgress []EpisodeProgress `json:"episodeProgress"`
	CreatedAt       time.Time         `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt       time.Time         `json:"deletedAt" gorm:"autoDeleteTime"`
}

func NewUser(username, passwordHash, email, name string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePassword(passwordHash); err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		Name:         name,
	}, nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return ErrInvalidUsernameLength
	}

	return nil
}

// Validation related regex
var (
	emailRegex      = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	hasNumberRegex  = regexp.MustCompile(`[0-9]`)
	hasUpperRegex   = regexp.MustCompile(`[A-Z]`)
	hasSpecialRegex = regexp.MustCompile(`[!@#$%^&*]`)
)

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasNumber := hasNumberRegex.MatchString(password)
	hasUpper := hasUpperRegex.MatchString(password)
	hasSpecial := hasSpecialRegex.MatchString(password)

	if !hasNumber || !hasUpper || !hasSpecial {
		return ErrPasswordRequirements
	}

	return nil
}
