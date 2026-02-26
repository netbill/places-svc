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

const operationDeleteUploadPlaceClassIcon = "delete_upload_place_class_icon"

func (c *Controller) DeletePlaceClassUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadPlaceClassIcon)

	req, err := requests.DeleteUploadPlaceClassIcon(r)
	if err != nil {
		log.WithError(err).Info("invalid delete upload place class icon requests")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	err = c.modules.pclass.DeleteUploadPlaceClassIcon(
		r.Context(),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("place class does not exist")
		c.responser.RenderErr(w, problems.NotFound("place class does not exist"))
	case errors.Is(err, errx.ErrorPlaceClassIconKeyIsInvalid):
		log.WithError(err).Info("place class icon key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to delete place class icon in upload session")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
