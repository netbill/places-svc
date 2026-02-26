package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdatePlaceStatus = "update_place_status"

func (c *Controller) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceStatus)

	req, err := requests.UpdatePlaceStatus(r)
	if err != nil {
		log.WithError(err).Info("invalid update Place status request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_id", req.Data.Id)

	res, err := c.modules.Place.UpdateStatus(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.Status,
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place not found")
		render.ResponseError(w, problems.NotFound("Place not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update Place status")
		render.ResponseError(w, problems.Forbidden("not enough rights to update Place status"))
	case errors.Is(err, errx.ErrorCannotSetStatusSuspend):
		log.Info("cannot set status suspended")
		render.ResponseError(w, problems.Forbidden("cannot set status suspended"))
	case errors.Is(err, errx.ErrorPlaceStatusIsInvalid):
		log.Info("place status is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"data/attributes/status": fmt.Errorf("place status is invalid"),
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update Place status")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Place(res))
	}
}
