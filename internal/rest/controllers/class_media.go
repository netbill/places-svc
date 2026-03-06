package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/classification"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreatePlaceClassUploadMediaLink = "create_place_class_upload_media_link"

func (c *PlaceClassController) CreateUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlaceClassUploadMediaLink)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/size": fmt.Errorf("invalid place class id: %s", chi.URLParam(r, "place_class_id")),
		})...)

		return
	}

	log = log.WithField("place_class_id", classID)

	class, media, err := c.class.CreateUploadMediaLinks(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class does not exist")
		render.ResponseError(w, problems.NotFound("place class does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to create place class upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadPlaceClassMediaLink(r, class, media))
	}
}

const operationDeleteUploadPlaceClassMedia = "delete_upload_place_class_nedia"

func (c *PlaceClassController) DeleteUploadMedia(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceClassMedia)

	req, err := requests.DeleteUploadPlaceClassMedia(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload place class icon requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	err = c.class.DeleteUploadMedia(
		r.Context(),
		req.Data.Id,
		classification.DeleteUploadPlaceClassMediaParams{
			Icon: req.Data.Attributes.IconKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class does not exist")
		render.ResponseError(w, problems.NotFound("place class does not exist"))
	case errors.Is(err, errx.ErrorPlaceClassIconIsInvalid):
		log.WithError(err).Warn("place class icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to delete place class icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
