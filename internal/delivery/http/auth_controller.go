package http

import (
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(authUseCase *usecase.AuthUseCase, logger *logrus.Logger) *AuthController {
	return &AuthController{
		Log:         logger,
		AuthUseCase: authUseCase,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RegisterResponse]{Data: response})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Refresh(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to refresh token: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RefreshTokenResponse]{Data: response})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)

	err := c.AuthUseCase.Logout(ctx.UserContext(), userID)
	if err != nil {
		c.Log.Warnf("Failed to logout: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: "Successfully logged out"})
}

func (c *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	request := new(model.VerifyEmailRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.VerifyEmail(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to verify email: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.VerifyEmailResponse]{Data: response})
}

func (c *AuthController) ResendVerification(ctx *fiber.Ctx) error {
	request := new(model.ResendVerificationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.ResendVerification(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to resend verification: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ResendVerificationResponse]{Data: response})
}
