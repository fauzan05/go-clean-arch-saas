package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey              string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

func NewJWTService(secret string, accessExpireMinutes, refreshExpireDays int) *JWTService {
	return &JWTService{
		secretKey:              secret,
		accessTokenExpiration:  time.Duration(accessExpireMinutes) * time.Minute,
		refreshTokenExpiration: time.Duration(refreshExpireDays) * 24 * time.Hour,
	}
}

func (s *JWTService) GenerateAccessToken(userID, email, orgID string) (string, error) {
	claims := &Claims{
		UserID:         userID,
		Email:          email,
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (s *JWTService) GetAccessTokenExpiration() time.Duration {
	return s.accessTokenExpiration
}

func (s *JWTService) GetRefreshTokenExpiration() time.Duration {
	return s.refreshTokenExpiration
}
