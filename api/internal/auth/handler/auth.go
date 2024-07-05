package handler

import (
	"github.com/labstack/echo/v4"
	"goechotemplate/api/internal/auth/dto"
	"goechotemplate/api/internal/auth/service"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.LoginRequest)
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
