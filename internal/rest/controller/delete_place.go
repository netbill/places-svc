package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeletePlace = "delete_place"

func (c *Controller) DeletePlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeletePlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid Place id")
		c.responser.RenderErr(w, problems.BadRequest(
			fmt.Errorf("invalid place_id: %s", chi.URLParam(r, "place_id")))...,
		)
		return
	}

	log = log.WithField("place_id", placeID)

	err = c.modules.Place.Delete(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place not found")
		c.responser.RenderErr(w, problems.NotFound("Place not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete Place")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete Place"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to delete Place")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK)
	}
}
