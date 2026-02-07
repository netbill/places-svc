package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

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
			DeleteIcon:      req.Data.Attributes.DeleteIcon != nil && *req.Data.Attributes.DeleteIcon,
			DeleteBanner:    req.Data.Attributes.DeleteBanner != nil && *req.Data.Attributes.DeleteBanner,
		},
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update place")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		case errors.Is(err, errx.ErrorPlaceNotExists):
			c.responser.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorPlaceIconTooLarge),
			errors.Is(err, errx.ErrorPlaceIconContentFormatNotAllowed),
			errors.Is(err, errx.ErrorPlaceIconContentTypeNotAllowed),
			errors.Is(err, errx.ErrorPlaceIconResolutionNotAllowed):
			c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
				"icon": fmt.Errorf(err.Error()),
			})...)
		case errors.Is(err, errx.ErrorPlaceBannerTooLarge),
			errors.Is(err, errx.ErrorPlaceBannerContentFormatNotAllowed),
			errors.Is(err, errx.ErrorPlaceBannerContentTypeNotAllowed),
			errors.Is(err, errx.ErrorPlaceBannerResolutionNotAllowed):
			c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
				"banner": fmt.Errorf(err.Error()),
			})...)
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.Place(res))
}
