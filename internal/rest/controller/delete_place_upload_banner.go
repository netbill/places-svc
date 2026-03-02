package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeleteUploadPlaceBanner = "delete_upload_place_banner"

func (c *Controller) DeletePlaceUploadBanner(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceBanner)

	req, err := requests.DeleteUploadPlaceBanner(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload place banner requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_id", req.Data.Id)

	err = c.modules.Place.DeleteUploadPlaceBanner(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.BannerKey,
	)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to delete place banner in upload session")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete place banner in upload session"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceBannerKeyIsInvalid):
		log.WithError(err).Warn("place banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place does not exist")
		render.ResponseError(w, problems.NotFound("place does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to delete place banner in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Debug("place banner in upload session is deleted successfully")
		w.WriteHeader(http.StatusNoContent)
	}
}
