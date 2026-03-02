package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
		log.WithError(err).Warn("invalid place_id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place_id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	place, err := c.modules.Place.Get(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		return
	case err != nil:
		log.WithError(err).Error("failed to get Place")
		render.ResponseError(w, problems.InternalError())
		return
	}

	slog.Debug(fmt.Sprintf("place found: %s", place.Name))

	opts := make([]responses.PlaceOption, 0)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 2)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	if slices.Contains(includes, "place_class") {
		class, err := c.modules.Class.Get(r.Context(), place.ClassID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("place_class_id", place.ClassID).Warn("place class not found")
			render.ResponseError(w, problems.NotFound(fmt.Sprintf("place class with id %s not found", place.ClassID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get place class")
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
			log.WithField("organization_id", place.OrganizationID).Warn("organization not found")
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
