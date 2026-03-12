package controller

import (
	"social_media/internal/service"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	svc service.ServiceInterface
}

func New(svc service.ServiceInterface) *Controller {
	return &Controller{svc: svc}
}

func (ctrl *Controller) Health(c echo.Context) error {
	return c.String(200, "OK")
}

func (ctrl *Controller) Ping(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "pong"})
}
