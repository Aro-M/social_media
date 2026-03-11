package handler

import (
	"social_media/internal/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/health", h.Health)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/ping", h.Ping)
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(200, "OK")
}

func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "pong"})
}
