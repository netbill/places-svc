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

const operationUpdatePlaceVerify = "update_place_verify"

func (c *Controller) UpdatePlaceVerify(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceVerify)

	req, err := requests.UpdatePlaceVerify(r)
	if err != nil {
		log.WithError(err).Info("invalid update Place verify request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_id", req.Data.Id)

	res, err := c.modules.Place.UpdateVerified(r.Context(), req.Data.Id, req.Data.Attributes.Verify)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place not found")
		c.responser.RenderErr(w, problems.NotFound("Place not found"))
	case err != nil:
		log.WithError(err).Error("failed to update Place verify")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Place(res))
	}
}
