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

const operationDeleteUploadPlaceIcon = "delete_upload_place_icon"

func (c *Controller) DeletePlaceUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceIcon)

	req, err := requests.DeleteUploadPlaceIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload place  icon requests")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_id", req.Data.Id)

	err = c.modules.place.DeleteUploadPlaceIcon(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete place icon in upload session")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete place icon in upload session"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceIconKeyIsInvalid):
		log.WithError(err).Info("place icon key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("place does not exist")
		c.responser.RenderErr(w, problems.NotFound("place does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to delete place icon in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
