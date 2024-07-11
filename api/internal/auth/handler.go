package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authService Service
}

func NewHandler(authService Service) Handler {
	return Handler{
		authService: authService,
	}
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(*req); err != nil {
		return err
	}

	res, err := h.authService.Login(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
