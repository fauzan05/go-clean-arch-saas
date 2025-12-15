package route

import (
	"fmt"
	"go-clean-arch-saas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                    *fiber.App
	AuthController         *http.AuthController
	UserController         *http.UserController
	OrganizationController *http.OrganizationController
	SubscriptionController *http.SubscriptionController
	HealthController       *http.HealthController
	AuthMiddleware         fiber.Handler
	Config                 *viper.Viper
}

// getAPIBasePath returns the base path for API routes (e.g., "/api/v1")
func (c *RouteConfig) getAPIBasePath() string {
	prefix := c.Config.GetString("api.prefix")
	version := c.Config.GetString("api.version")
	return fmt.Sprintf("%s/%s", prefix, version)
}

func (c *RouteConfig) Setup() {
	c.SetupHealthRoutes()
	c.SetupGuestRoutes()
	c.SetupAuthRoutes()
}

func (c *RouteConfig) SetupHealthRoutes() {
	c.App.Get("/health", c.HealthController.Health)
	c.App.Get("/ready", c.HealthController.Ready)
}

func (c *RouteConfig) SetupGuestRoutes() {
	api := c.App.Group(c.getAPIBasePath())

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/refresh", c.AuthController.Refresh)
	auth.Post("/verify-email", c.AuthController.VerifyEmail)
	auth.Post("/resend-verification", c.AuthController.ResendVerification)
}

func (c *RouteConfig) SetupAuthRoutes() {
	api := c.App.Group(c.getAPIBasePath())
	api.Use(c.AuthMiddleware)

	// Auth routes (authenticated)
	auth := api.Group("/auth")
	auth.Delete("/logout", c.AuthController.Logout)

	// User routes
	users := api.Group("/users")
	users.Get("/current", c.UserController.Current)
	users.Patch("/current", c.UserController.Update)

	// Organization routes
	orgs := api.Group("/organizations")
	orgs.Get("/current", c.OrganizationController.GetCurrent)
	orgs.Patch("/current", c.OrganizationController.Update)
	orgs.Get("/members", c.OrganizationController.ListMembers)
	orgs.Delete("/members/:userId", c.OrganizationController.RemoveMember)

	// Subscription routes
	subs := api.Group("/subscriptions")
	subs.Get("/current", c.SubscriptionController.GetCurrent)
	subs.Post("/upgrade", c.SubscriptionController.Upgrade)
	subs.Post("/cancel", c.SubscriptionController.Cancel)
}
