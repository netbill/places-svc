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
)

const operationUpdatePlaceClass = "update_place_class"

func (c *Controller) UpdatePlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceClass)

	req, err := requests.UpdatePlaceClass(r)
	if err != nil {
		log.WithError(err).Info("invalid update place class request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	res, err := c.modules.pclass.Update(
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
		log.WithError(err).Info("place class not found")
		c.responser.RenderErr(w, problems.NotFound("place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassParentCycle):
		log.WithError(err).Info("setting parent would create a cycle")
		c.responser.RenderErr(w, problems.Conflict("setting parent would create a cycle"))
	case errors.Is(err, errx.ErrorPlaceClassIconKeyIsInvalid):
		log.WithError(err).Info("icon key is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconFormatIsNotAllowed):
		log.WithError(err).Info("icon format is not allowed")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconContentIsExceedsMax):
		log.WithError(err).Info("icon content is exceeds max")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceClassIconResolutionIsInvalid):
		log.WithError(err).Info("icon resolution is invalid")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update place class")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Info("place class updated")
		c.responser.Render(w, http.StatusOK, responses.PlaceClass(res))
	}
}
