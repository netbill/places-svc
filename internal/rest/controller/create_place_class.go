package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) CreatePlaceClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.CreatePlaceClass(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create place class request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.core.pclass.Create(r.Context(), pclass.CreateParams{
		ParentID:    req.Data.Attributes.ParentId,
		Code:        req.Data.Attributes.Code,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Icon:        req.Data.Attributes.Icon,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create place class")
		switch {
		case errors.Is(errx.ErrorPlaceClassNotExists, err):
			c.responser.RenderErr(w, problems.NotFound("parent place class not found"))
		case errors.Is(errx.ErrorPlaceClassParentCycle, err):
			c.responser.RenderErr(w, problems.Conflict("place class parent cycle detected"))
		case errors.Is(errx.ErrorPlaceClassCodeExists, err):
			c.responser.RenderErr(w, problems.Conflict("place class code already in use"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusCreated, responses.PlaceClass(res))
}
