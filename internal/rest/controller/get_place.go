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
)

const operationGetPlace = "get_place"

func (c *Controller) GetPlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place id"))...)
		return
	}

	log = log.WithField("place_id", placeID)

	place, err := c.modules.place.Get(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("place not found")
		c.responser.RenderErr(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		return
	case err != nil:
		log.WithError(err).Error("failed to get place")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.PlaceOption, 0)

	if slices.Contains(includes, "class") {
		class, err := c.modules.pclass.Get(r.Context(), place.ClassID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("place_class_id", place.ClassID).Info("place class not found")
			c.responser.RenderErr(w, problems.NotFound(fmt.Sprintf("place class with id %s not found", place.ClassID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get place class")
			c.responser.RenderErr(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithClass(class))
	}

	if slices.Contains(includes, "organization") {
		org, err := c.modules.organization.Get(r.Context(), place.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists):
			log.WithField("organization_id", place.OrganizationID).Info("organization not found")
			c.responser.RenderErr(w, problems.NotFound(fmt.Sprintf("organization with id %s not found", place.OrganizationID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			c.responser.RenderErr(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithOrganization(org))
	}

	c.responser.Render(w, http.StatusOK, responses.Place(place, opts...))
}
