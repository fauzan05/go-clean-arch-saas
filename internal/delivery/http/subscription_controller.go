package http

import (
	"go-clean-arch-saas/internal/delivery/http/middleware"
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SubscriptionController struct {
	Log     *logrus.Logger
	UseCase *usecase.SubscriptionUseCase
}

func NewSubscriptionController(useCase *usecase.SubscriptionUseCase, logger *logrus.Logger) *SubscriptionController {
	return &SubscriptionController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *SubscriptionController) GetCurrent(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	response, err := c.UseCase.GetCurrentSubscription(ctx.UserContext(), orgID)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current subscription")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SubscriptionResponse]{Data: response})
}

func (c *SubscriptionController) Upgrade(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	request := new(model.UpgradeSubscriptionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	request.OrganizationID = orgID
	response, err := c.UseCase.Upgrade(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to upgrade subscription")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SubscriptionResponse]{Data: response})
}

func (c *SubscriptionController) Cancel(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	request := &model.CancelSubscriptionRequest{
		OrganizationID: orgID,
	}

	err := c.UseCase.Cancel(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to cancel subscription")
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: "Subscription cancelled successfully"})
}
