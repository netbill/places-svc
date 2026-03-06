package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/places"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationCreatePlaceUploadMediaLink = "create_place_upload_media_link"

func (c *PlaceController) CreateUploadMediaLink(w http.ResponseWriter, r *http.Request) {
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

	place, media, err := c.place.CreateUploadMediaLinks(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists),
		errors.Is(err, errx.ErrorPlaceDeleted):
		log.WithError(err).Warn("place does not exist")
		render.ResponseError(w, problems.NotFound("place does not exist"))
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("initiator is not a member of organization")
		render.ResponseError(w, problems.Forbidden("initiator is not a member of organization"))
	case err != nil:
		log.WithError(err).Error("failed to create place upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadPlaceMediaLink(r, place, media))
	}
}

const operationDeleteUploadPlaceMedia = "delete_upload_place_media"

func (c *PlaceController) DeleteUploadMedia(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceMedia)

	req, err := requests.DeleteUploadPlaceMedia(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload place banner requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_id", req.Data.Id)

	err = c.place.DeleteUploadPlaceMedia(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		places.DeleteUploadPlaceMediaParams{
			Banner: req.Data.Attributes.BannerKey,
			Icon:   req.Data.Attributes.IconKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists),
		errors.Is(err, errx.ErrorPlaceDeleted):
		log.WithError(err).Warn("place does not exist")
		render.ResponseError(w, problems.NotFound("place does not exist"))
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("initiator is not a member of organization")
		render.ResponseError(w, problems.Forbidden("initiator is not a member of organization"))
	case errors.Is(err, errx.ErrorPlaceBannerIsInvalid):
		log.WithError(err).Warn("place banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconIsInvalid):
		log.WithError(err).Warn("place icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to delete place banner in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Debug("place banner in upload session is deleted successfully")
		render.Response(w, http.StatusOK, nil)
	}
}
