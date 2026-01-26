package identity

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type SecretKey []byte

func (s *SecretKey) UnmarshalText(text []byte) error {
	var err error
	if *s, err = base64.StdEncoding.DecodeString(string(text)); err != nil {
		return fmt.Errorf("failed to decode secret key: %w", err)
	}
	return nil
}

type TokenHandler struct {
	secret SecretKey
}

func NewTokenHandler(secret []byte) TokenHandler {
	return TokenHandler{secret}
}

type UserClaims struct {
	jwt.RegisteredClaims

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ProfileImageURL string `json:"profile_image"`
}

func (c *UserClaims) ID() int64 {
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return id
}

func (h *TokenHandler) Sign(claims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte(h.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func (h *TokenHandler) Verify(token string) (*UserClaims, error) {
	tok, err := jwt.ParseWithClaims(token, new(UserClaims), func(t *jwt.Token) (any, error) {
		return []byte(h.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	claims, ok := tok.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return claims, nil
}
