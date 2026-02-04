package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c Controller) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := requests.UpdatePlaceStatus(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place status data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	res, err := c.core.UpdatePlaceStatus(r.Context(), initiator.ID, req.Data.Id, req.Data.Attributes.Status)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update status place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		case errors.Is(err, errx.ErrorPlaceStatusSuspended):
			ape.RenderErr(w, problems.Forbidden("place status is suspend"))
		case errors.Is(err, errx.ErrorCannotSetPlaceStatusSuspend):
			ape.RenderErr(w, problems.Forbidden("cannot set status suspended"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
