package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeletePlace = "delete_place"

func (c *Controller) DeletePlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeletePlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalidplaceid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalidplaceid: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	err = c.modules.Place.Delete(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists) || errors.Is(err, errx.ErrorPlaceDeleted):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to delete Place")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete Place"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to delete Place")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, nil)
	}
}
