package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeleteUploadPlaceBanner = "delete_upload_place_banner"

func (c *Controller) DeletePlaceUploadBanner(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceBanner)

	req, err := requests.DeleteUploadPlaceBanner(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload Place  banner requests")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

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
		log.Info("not enough rights to delete Place banner in upload session")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete Place banner in upload session"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceBannerKeyIsInvalid):
		log.WithError(err).Info("Place banner key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place does not exist")
		c.responser.RenderErr(w, problems.NotFound("Place does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to delete Place banner in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
