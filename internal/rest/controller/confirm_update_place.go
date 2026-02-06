package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

// TOOD remade
func (c *Controller) ConfirmUpdatePlace(w http.ResponseWriter, r *http.Request) {
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

	req, err := requests.ConfirmUpdatePlace(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place request data")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.place.ConfirmUpdateSession(r.Context(), initiator, req.Data.Id, place.UpdateParams{
		Address:     req.Data.Attributes.Address,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Website:     req.Data.Attributes.Website,
		Phone:       req.Data.Attributes.Phone,
		Media: place.UpdateMediaParams{
			UploadSessionID: uploadFilesData.GetUploadSessionID(),
			DeleteIcon:      req.Data.Attributes.DeleteIcon,
			DeleteBanner:    req.Data.Attributes.DeleteBanner,
		},
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			c.responser.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Place(res))
}
