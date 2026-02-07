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

func (c *Controller) CancelUpdatePlaceClass(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))

		return
	}

	uploadFilesData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	placeClassID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_class_id")),
		})...)

		return
	}

	err = c.core.place.CancelUpdate(r.Context(), initiator, placeClassID, uploadFilesData.GetUploadSessionID())
	if err != nil {
		c.log.WithError(err).Errorf("failed to cancel update place class session")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusNoContent, nil)
}
