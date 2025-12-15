package route

import (
	"go-clean-arch-saas/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                    *fiber.App
	AuthController         *http.AuthController
	UserController         *http.UserController
	OrganizationController *http.OrganizationController
	SubscriptionController *http.SubscriptionController
	HealthController       *http.HealthController
	AuthMiddleware         fiber.Handler
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
	api := c.App.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/refresh", c.AuthController.Refresh)
}

func (c *RouteConfig) SetupAuthRoutes() {
	api := c.App.Group("/api/v1")
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
