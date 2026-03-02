package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdatePlace = "update_place"

func (c *Controller) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlace)

	req, err := requests.UpdatePlace(r)
	if err != nil {
		log.WithError(err).Warn("invalid updateplace request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_id", req.Data.Id)

	res, err := c.modules.Place.Update(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		place.UpdateParams{
			ClassID:     req.Data.Attributes.ClassId,
			Name:        req.Data.Attributes.Name,
			Address:     req.Data.Attributes.Address,
			Description: req.Data.Attributes.Description,
			Website:     req.Data.Attributes.Website,
			Phone:       req.Data.Attributes.Phone,
			IconKey:     req.Data.Attributes.IconKey,
			BannerKey:   req.Data.Attributes.BannerKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place  not found")
		render.ResponseError(w, problems.NotFound("place  not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to update place")
		render.ResponseError(w, problems.Forbidden("not enough rights to update place"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Conflict("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithError(err).Warn("place class is deprecated")
		render.ResponseError(w, problems.Conflict("place class is deprecated"))
	case errors.Is(err, errx.ErrorPlaceIconKeyIsInvalid):
		log.WithError(err).Warn("icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconFormatIsNotAllowed):
		log.WithError(err).Warn("icon format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconContentIsExceedsMax):
		log.WithError(err).Warn("icon content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconResolutionIsInvalid):
		log.WithError(err).Warn("icon resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerKeyIsInvalid):
		log.WithError(err).Warn("banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerFormatIsNotAllowed):
		log.WithError(err).Warn("banner format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerContentIsExceedsMax):
		log.WithError(err).Warn("banner content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerResolutionIsInvalid):
		log.WithError(err).Warn("banner resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update place")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place  updated")
		render.Response(w, http.StatusOK, responses.Place(res))
	}
}
