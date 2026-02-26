package controller

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdatePlaceClass = "update_place_class"

func (c *Controller) UpdatePlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceClass)

	req, err := requests.UpdatePlaceClass(r)
	if err != nil {
		log.WithError(err).Info("invalid update Place class request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	res, err := c.modules.Class.Update(
		r.Context(),
		req.Data.Id,
		pclass.UpdateParams{
			ParentID:    req.Data.Attributes.ParentId,
			Name:        req.Data.Attributes.Name,
			Description: req.Data.Attributes.Description,
			IconKey:     req.Data.Attributes.IconKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Info("Place class not found")
		render.ResponseError(w, problems.NotFound("Place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassParentCycle):
		log.WithError(err).Info("setting parent would create a cycle")
		render.ResponseError(w, problems.Conflict("setting parent would create a cycle"))
	case errors.Is(err, errx.ErrorPlaceClassIconKeyIsInvalid):
		log.WithError(err).Info("icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconFormatIsNotAllowed):
		log.WithError(err).Info("icon format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconContentIsExceedsMax):
		log.WithError(err).Info("icon content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconResolutionIsInvalid):
		log.WithError(err).Info("icon resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update Place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("Place class updated")
		render.Response(w, http.StatusOK, responses.PlaceClass(res))
	}
}
