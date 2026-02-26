package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetPlace = "get_place"

func (c *Controller) GetPlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid Place id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid Place id"))...)
		return
	}

	log = log.WithField("place_id", placeID)

	place, err := c.modules.Place.Get(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place not found")
		render.ResponseError(w, problems.NotFound(fmt.Sprintf("Place with id %s not found", placeID)))
		return
	case err != nil:
		log.WithError(err).Error("failed to get Place")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.PlaceOption, 0)

	if slices.Contains(includes, "place_class") {
		class, err := c.modules.Class.Get(r.Context(), place.ClassID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("place_class_id", place.ClassID).Info("Place class not found")
			render.ResponseError(w, problems.NotFound(fmt.Sprintf("Place class with id %s not found", place.ClassID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get Place class")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithClass(class))
		}
	}

	if slices.Contains(includes, "organization") {
		org, err := c.modules.Org.Get(r.Context(), place.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists):
			log.WithField("organization_id", place.OrganizationID).Info("organization not found")
			render.ResponseError(w, problems.NotFound(fmt.Sprintf("organization with id %s not found", place.OrganizationID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithOrganization(org))
		}
	}

	render.Response(w, http.StatusOK, responses.Place(place, opts...))
}
