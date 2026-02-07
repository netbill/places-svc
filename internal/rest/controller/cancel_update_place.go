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

func (c *Controller) CancelUpdatePlace(w http.ResponseWriter, r *http.Request) {
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

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)

		return
	}

	err = c.core.place.CancelUpdate(r.Context(), initiator, placeID, uploadFilesData.GetUploadSessionID())
	if err != nil {
		c.log.WithError(err).Errorf("failed to cancel update place session")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to cancel update place"))
		case errors.Is(err, errx.ErrorPlaceNotExists):
			c.responser.RenderErr(w, problems.NotFound("place not found"))
		}

		return
	}

	c.responser.Render(w, http.StatusNoContent, nil)
}
