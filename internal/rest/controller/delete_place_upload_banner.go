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
		log.WithError(err).Info("invalid delete upload place  banner requests")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_id", req.Data.Id)

	err = c.modules.place.DeleteUploadPlaceBanner(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.BannerKey,
	)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete place banner in upload session")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete place banner in upload session"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceBannerKeyIsInvalid):
		log.WithError(err).Info("place banner key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("place does not exist")
		c.responser.RenderErr(w, problems.NotFound("place does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to delete place banner in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
