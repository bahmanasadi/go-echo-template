package person

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	jwtmodel "goechotemplate/api/internal/auth"
	"net/http"
)

type Handler struct {
	authService   jwtmodel.Service
	personService Service
}

func NewPersonHandler(authService jwtmodel.Service, personService Service) *Handler {
	return &Handler{
		authService:   authService,
		personService: personService,
	}
}

func (h *Handler) GetPerson(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	authenticatedPersonID, err := c.Get(jwtmodel.DefaultJWTConfig.ContextKey).(*jwt.Token).Claims.GetSubject()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err = h.authService.VerifyAuthorisation(ctx, jwtmodel.VerifyAuthorisationParams{
		AuthenticatedPersonID: &authenticatedPersonID,
		TargetPersonID:        &id,
	})
	if err != nil {
		return echo.ErrForbidden
	}

	user, err := h.personService.GetByExternalID(ctx, id)
	if err != nil {
		c.Logger().Errorf("GetPerson: %v", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreatePerson(c echo.Context) error {
	ctx := c.Request().Context()

	person := new(Person)
	if err := c.Bind(person); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(person); err != nil {
		return err
	}

	createdPerson, err := h.personService.Create(ctx, person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdPerson)
}
