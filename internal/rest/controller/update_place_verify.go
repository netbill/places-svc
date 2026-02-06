package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdatePlaceVerify(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlaceVerify(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place verify data")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.place.UpdateVerified(r.Context(), req.Data.Id, req.Data.Attributes.Verify)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place verify data")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			c.responser.RenderErr(w, problems.NotFound("place not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Place(res))
}
