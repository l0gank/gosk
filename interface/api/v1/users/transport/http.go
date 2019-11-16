package transport

import (
	"github.com/labstack/echo"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/domain/services/users"
	"net/http"
	"strconv"
)

// HTTP represents user http service
type HTTP struct {
	svc users.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc users.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/users")
	ur.POST("", h.create)
	ur.GET("", h.find)
	ur.GET("/id", h.view)
	ur.GET("/id/:id", h.view)
	ur.PATCH("", h.update)
	ur.DELETE("/id/:id", h.delete)
}

func (h *HTTP) create(c echo.Context) error {
	r := new(model.CreateUserReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	resp, _ := h.svc.Create(r)

	return c.JSON(resp.Code, resp)
}

func (h *HTTP) find(c echo.Context) error {
	var resp *entity.Response
	p := new(model.PaginationReq)
	req := new(model.FilterUserReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Respond(err, nil, http.StatusBadRequest))
	}

	p.Limit = req.Limit
	p.Page = req.Page

	result, total, err := h.svc.Find(p.Transform(), req)

	if err != nil {
		return c.JSON(http.StatusInsufficientStorage, entity.Respond(err, nil, http.StatusInsufficientStorage))
	}

	page := model.Page(p.Limit, p.Page, total)

	resp = entity.Respond(err, entity.ListRepond(result, page), http.StatusOK)

	return c.JSON(resp.Code, resp)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if id != 0 {
			return c.JSON(http.StatusBadRequest, entity.Respond(err, nil, http.StatusBadRequest))
		}
	}

	resp, _ := h.svc.View(&model.UserReq{Id: int64(id)})

	return c.JSON(resp.Code, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.Respond(err, nil, http.StatusBadRequest))
	}

	resp, _ := h.svc.Delete(&model.UserReq{Id: int64(id)})

	return c.JSON(resp.Code, resp)
}

func (h *HTTP) update(c echo.Context) error {
	req := new(model.UpdateUserReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Respond(err, nil, http.StatusBadRequest))
	}

	resp, _ := h.svc.Update(req)

	return c.JSON(resp.Code, resp)
}
