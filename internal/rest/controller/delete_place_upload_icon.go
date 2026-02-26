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

const operationDeleteUploadPlaceIcon = "delete_upload_place_icon"

func (c *Controller) DeletePlaceUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceIcon)

	req, err := requests.DeleteUploadPlaceIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload Place  icon requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_id", req.Data.Id)

	err = c.modules.Place.DeleteUploadPlaceIcon(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to delete Place icon in upload session")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete Place icon in upload session"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceIconKeyIsInvalid):
		log.WithError(err).Info("Place icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.Info("Place does not exist")
		render.ResponseError(w, problems.NotFound("Place does not exist"))
	case err != nil:
		log.WithError(err).Error("failed to delete Place icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
