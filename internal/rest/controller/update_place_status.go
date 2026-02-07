package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))

		return
	}

	req, err := requests.UpdatePlaceStatus(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place status data")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}
	res, err := c.core.place.UpdateStatus(r.Context(), initiator, req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update status place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotExists):
			c.responser.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.Place(res))
}
