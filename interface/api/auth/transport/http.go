package transport

import (
	"github.com/labstack/echo"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/services/auth"
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

// NewHTTP creates new auth http service
func NewHTTP(svc auth.Service, e *echo.Echo, mw echo.MiddlewareFunc) {
	h := HTTP{svc}
	e.POST("/login", h.login)
	e.GET("/me", h.me, mw)
}

func (h *HTTP) login(c echo.Context) error {
	cred := new(entity.Credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}
	r, err := h.svc.Authenticate(cred.Username, cred.Password)
	if err != nil {
		return err
	}
	return c.JSON(r.Code, r)
}

func (h *HTTP) me(c echo.Context) error {
	res, _ := h.svc.Me(c.Get("email").(string))
	return c.JSON(res.Code, res)
}
