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

const operationCreatePlaceClassUploadMediaLink = "create_place_class_upload_media_link"

func (c *Controller) CreatePlaceClassUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlaceClassUploadMediaLink)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Info("invalid Place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"place_class_id": err,
		})...)

		return
	}

	class, media, err := c.modules.Class.CreateUploadMediaLinks(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("Place class does not exist")
		render.ResponseError(w, problems.NotFound("Place class does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to create Place class upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadPlaceClassMediaLink(class, media))
	}
}
