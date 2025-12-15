package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}
