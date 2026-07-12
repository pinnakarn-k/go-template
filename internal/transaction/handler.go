package transaction

import (
	"context"
	"errors"
	"transaction-api/internal/httpresponse"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type handlerService interface {
	Search(ctx context.Context, req SearchRequest) (SearchOutput, error)
}

type Handler struct {
	service  handlerService
	validate *validator.Validate
}

func NewHandler(service handlerService) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Handler) Search(c *fiber.Ctx) error {
	var req SearchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "accountId, transactedAtStart and transactedAtEnd are required",
		})
	}

	output, err := h.service.Search(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidTimeRange) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(httpresponse.Response[[]SearchItem]{
		Data: output.Items,
		Meta: &output.Meta,
	})
}
