package controller

import (
	"errors"
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

	ID, err := uuid.Parse(chi.URLParam(r, "place__id"))
	if err != nil {
		log.WithError(err).Info("invalid Place  id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"place__id": err,
		})...)

		return
	}

	place, media, err := c.modules.Place.CreateUploadMediaLinks(r.Context(), ID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place does not exist")
		render.ResponseError(w, problems.NotFound("Place does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create Place upload media link")
		render.ResponseError(w, problems.Forbidden("not enough rights to create Place upload media link"))
	case err != nil:
		log.WithError(err).Error("failed to create Place upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadPlaceMediaLink(place, media))
	}
}
