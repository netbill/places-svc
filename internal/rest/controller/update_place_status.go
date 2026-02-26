package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdatePlaceStatus = "update_place_status"

func (c *Controller) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceStatus)

	req, err := requests.UpdatePlaceStatus(r)
	if err != nil {
		log.WithError(err).Info("invalid update place status request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_id", req.Data.Id)

	res, err := c.modules.place.UpdateStatus(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.Status,
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("place not found")
		c.responser.RenderErr(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update place status")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place status"))
	case err != nil:
		log.WithError(err).Error("failed to update place status")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Place(res))
	}
}
