package middleware

import (
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(authUseCase *usecase.AuthUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			authUseCase.Log.Warn("Missing authorization header")
			return fiber.ErrUnauthorized
		}

		// Remove "Bearer " prefix
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		authUseCase.Log.Debugf("Validating token: %s", token[:10]+"...")

		auth, err := authUseCase.VerifyToken(ctx.UserContext(), token)
		if err != nil {
			authUseCase.Log.Warnf("Failed to verify token: %+v", err)
			return fiber.ErrUnauthorized
		}

		authUseCase.Log.Debugf("Authenticated user: %s, org: %s", auth.UserID, auth.OrganizationID)

		// Set auth context
		ctx.Locals("auth", auth)
		ctx.Locals("user_id", auth.UserID)
		ctx.Locals("email", auth.Email)
		ctx.Locals("organization_id", auth.OrganizationID)

		return ctx.Next()
	}
}

func GetAuth(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}

func GetUserID(ctx *fiber.Ctx) string {
	return ctx.Locals("user_id").(string)
}

func GetOrganizationID(ctx *fiber.Ctx) string {
	return ctx.Locals("organization_id").(string)
}
