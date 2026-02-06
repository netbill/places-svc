package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeletePlaceClassUploadIcon(w http.ResponseWriter, r *http.Request) {
	placeClassID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid placeClass id: %s", chi.URLParam(r, "place_class_id")),
		})...)

		return
	}

	uploadContentData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	if err = c.core.class.DeleteUpdateIconInSession(
		r.Context(),
		placeClassID,
		uploadContentData.GetUploadSessionID(),
	); err != nil {
		c.log.WithError(err).Errorf("failed to delete place class icon in upload session")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			c.responser.RenderErr(w, problems.NotFound("place class does not exist"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place class"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}
}
