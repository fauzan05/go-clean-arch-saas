package config

import (
	jwtPkg "go-clean-arch-saas/pkg/jwt"

	"github.com/spf13/viper"
)

func NewJWT(config *viper.Viper) *jwtPkg.JWTService {
	secret := config.GetString("jwt.secret")
	accessExpireMinutes := config.GetInt("jwt.access_expire_minutes")
	refreshExpireDays := config.GetInt("jwt.refresh_expire_days")

	return jwtPkg.NewJWTService(secret, accessExpireMinutes, refreshExpireDays)
}
