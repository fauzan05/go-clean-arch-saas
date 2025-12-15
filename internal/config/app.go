package config

import (
	"go-clean-arch-saas/internal/delivery/http"
	"go-clean-arch-saas/internal/delivery/http/middleware"
	"go-clean-arch-saas/internal/delivery/http/route"
	"go-clean-arch-saas/internal/repository"
	"go-clean-arch-saas/internal/usecase"
	"go-clean-arch-saas/pkg/email"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup JWT service
	jwtService := NewJWT(config.Config)

	// setup Email service
	emailService := email.NewEmailService(
		config.Config.GetString("email.host"),
		config.Config.GetInt("email.port"),
		config.Config.GetString("email.username"),
		config.Config.GetString("email.password"),
		config.Config.GetString("email.from"),
		config.Log,
	)

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	organizationRepository := repository.NewOrganizationRepository(config.Log)
	organizationMemberRepository := repository.NewOrganizationMemberRepository(config.Log)
	planRepository := repository.NewPlanRepository(config.Log)
	subscriptionRepository := repository.NewSubscriptionRepository(config.Log)

	// setup use cases
	authUseCase := usecase.NewAuthUseCase(
		config.DB,
		config.Log,
		config.Validate,
		userRepository,
		organizationRepository,
		organizationMemberRepository,
		planRepository,
		subscriptionRepository,
		jwtService,
		emailService,
		config.Config.GetString("base_url"),
	)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	organizationUseCase := usecase.NewOrganizationUseCase(
		config.DB,
		config.Log,
		config.Validate,
		organizationRepository,
		organizationMemberRepository,
	)
	subscriptionUseCase := usecase.NewSubscriptionUseCase(
		config.DB,
		config.Log,
		config.Validate,
		subscriptionRepository,
		planRepository,
	)

	// setup controllers
	authController := http.NewAuthController(authUseCase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	organizationController := http.NewOrganizationController(organizationUseCase, config.Log)
	subscriptionController := http.NewSubscriptionController(subscriptionUseCase, config.Log)
	healthController := http.NewHealthController(config.DB, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(authUseCase)

	routeConfig := route.RouteConfig{
		App:                    config.App,
		AuthController:         authController,
		UserController:         userController,
		OrganizationController: organizationController,
		SubscriptionController: subscriptionController,
		HealthController:       healthController,
		AuthMiddleware:         authMiddleware,
		Config:                 config.Config,
	}
	routeConfig.Setup()
}
