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

const operationDeleteUploadPlaceClassIcon = "delete_upload_place_class_icon"

func (c *Controller) DeletePlaceClassUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceClassIcon)

	req, err := requests.DeleteUploadPlaceClassIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload Place class icon requests")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	err = c.modules.Class.DeleteUploadPlaceClassIcon(
		r.Context(),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("Place class does not exist")
		render.ResponseError(w, problems.NotFound("Place class does not exist"))
	case errors.Is(err, errx.ErrorPlaceClassIconKeyIsInvalid):
		log.WithError(err).Info("Place class icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to delete Place class icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
