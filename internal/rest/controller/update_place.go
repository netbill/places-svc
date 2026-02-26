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
		log.WithError(err).Info("invalid update Place  request")
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
		log.WithError(err).Info("Place  not found")
		render.ResponseError(w, problems.NotFound("Place  not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Info("not enough rights to update Place ")
		render.ResponseError(w, problems.Forbidden("not enough rights to update Place "))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Info("organization is suspended")
		render.ResponseError(w, problems.Conflict("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithError(err).Info("Place class is deprecated")
		render.ResponseError(w, problems.Conflict("Place class is deprecated"))
	case errors.Is(err, errx.ErrorPlaceIconKeyIsInvalid):
		log.WithError(err).Info("icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconFormatIsNotAllowed):
		log.WithError(err).Info("icon format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconContentIsExceedsMax):
		log.WithError(err).Info("icon content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceIconResolutionIsInvalid):
		log.WithError(err).Info("icon resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerKeyIsInvalid):
		log.WithError(err).Info("banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerFormatIsNotAllowed):
		log.WithError(err).Info("banner format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerContentIsExceedsMax):
		log.WithError(err).Info("banner content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerResolutionIsInvalid):
		log.WithError(err).Info("banner resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update Place ")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("Place  updated")
		render.Response(w, http.StatusOK, responses.Place(res))
	}
}
