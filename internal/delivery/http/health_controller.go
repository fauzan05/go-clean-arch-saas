package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthController struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewHealthController(db *gorm.DB, logger *logrus.Logger) *HealthController {
	return &HealthController{
		Log: logger,
		DB:  db,
	}
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (c *HealthController) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(HealthResponse{
		Status: "ok",
	})
}

func (c *HealthController) Ready(ctx *fiber.Ctx) error {
	response := ReadinessResponse{
		Status:   "ok",
		Database: "disconnected",
	}

	// Check database connection
	sqlDB, err := c.DB.DB()
	if err == nil {
		if err := sqlDB.Ping(); err == nil {
			response.Database = "connected"
		}
	}

	return ctx.JSON(response)
}
