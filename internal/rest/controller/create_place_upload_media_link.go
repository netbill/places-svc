package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreatePlaceUploadMediaLink = "create_place__upload_media_link"

func (c *Controller) CreatePlaceUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlaceUploadMediaLink)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)

		return
	}

	log = log.WithField("place_id", placeID)

	place, media, err := c.modules.Place.CreateUploadMediaLinks(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place does not exist")
		render.ResponseError(w, problems.NotFound("place does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to createplaceupload media link")
		render.ResponseError(w, problems.Forbidden("not enough rights to createplaceupload media link"))
	case err != nil:
		log.WithError(err).Error("failed to createplaceupload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadPlaceMediaLink(place, media))
	}
}
