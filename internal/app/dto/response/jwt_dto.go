package response

import "time"

type CompleteJWTDto struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	Expiration   time.Time `json:"expiration"`
}

type PartialJWTDto struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

func NewCompleteJWTDto(token, refreshToken string, expiration time.Time) *CompleteJWTDto {
	return &CompleteJWTDto{
		Token:        token,
		RefreshToken: refreshToken,
		Expiration:   expiration,
	}
}

func NewPartialJWTDto(token string, expiration time.Time) *PartialJWTDto {
	return &PartialJWTDto{
		Token:      token,
		Expiration: expiration,
	}
}
