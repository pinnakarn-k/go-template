package transaction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *Handler
}

func NewModule(db *pgxpool.Pool) *Module {
	repository := NewPostgresRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	transactions := router.Group("/transactions")

	transactions.Post("/search", m.handler.Search)
}
