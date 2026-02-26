package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationCreatePlaceClass = "create_place_class"

func (c *Controller) CreatePlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlaceClass)

	req, err := requests.CreatePlaceClass(r)
	if err != nil {
		log.WithError(err).Info("invalid create Place class request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.modules.Class.Create(r.Context(), pclass.CreateParams{
		ParentID:    req.Data.Attributes.ParentId,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithField("place_class_id", req.Data.Attributes.ParentId).
			Info("parent Place class not found")
		c.responser.RenderErr(w, problems.NotFound("parent Place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithField("place_class_id", req.Data.Attributes.ParentId).
			Info("parent Place class is deprecated")
		c.responser.RenderErr(w, problems.Conflict("parent Place class is deprecated"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create Place class")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to create Place class"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
	case err != nil:
		log.WithError(err).Error("failed to create Place class")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.WithField("place_class_id", res.ID).Info("Place class created")
		c.responser.Render(w, http.StatusCreated, responses.PlaceClass(res))
	}
}
