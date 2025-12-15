package http

import (
	"go-clean-arch-saas/internal/delivery/http/middleware"
	"go-clean-arch-saas/internal/model"
	"go-clean-arch-saas/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrganizationController struct {
	Log     *logrus.Logger
	UseCase *usecase.OrganizationUseCase
}

func NewOrganizationController(useCase *usecase.OrganizationUseCase, logger *logrus.Logger) *OrganizationController {
	return &OrganizationController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *OrganizationController) GetCurrent(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	response, err := c.UseCase.GetByID(ctx.UserContext(), orgID)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current organization")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OrganizationResponse]{Data: response})
}

func (c *OrganizationController) Update(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	request := new(model.UpdateOrganizationRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	request.ID = orgID
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to update organization")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.OrganizationResponse]{Data: response})
}

func (c *OrganizationController) ListMembers(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)

	response, err := c.UseCase.ListMembers(ctx.UserContext(), orgID)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to list organization members")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.OrganizationMemberResponse]{Data: response})
}

func (c *OrganizationController) RemoveMember(ctx *fiber.Ctx) error {
	orgID := middleware.GetOrganizationID(ctx)
	userID := ctx.Params("userId")

	request := &model.RemoveOrganizationMemberRequest{
		OrganizationID: orgID,
		UserID:         userID,
	}

	err := c.UseCase.RemoveMember(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to remove organization member")
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: "Member removed successfully"})
}
