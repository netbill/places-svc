package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/paulmach/orb"
)

const operationCreatePlace = "create_place"

func (c *Controller) CreatePlace(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlace)

	req, err := requests.CreatePlace(r)
	if err != nil {
		log.WithError(err).Info("invalid create place request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId).
		WithField("class_id", req.Data.Attributes.ClassId)

	res, err := c.modules.place.Create(r.Context(), scope.AccountActor(r), place.CreateParams{
		OrganizationID: req.Data.Attributes.OrganizationId,
		ClassID:        req.Data.Attributes.ClassId,
		Address:        req.Data.Attributes.Address,
		Name:           req.Data.Attributes.Name,
		Description:    req.Data.Attributes.Description,
		Website:        req.Data.Attributes.Website,
		Phone:          req.Data.Attributes.Phone,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
	})
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to create place")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to create place"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.Info("organization is suspended")
		c.responser.RenderErr(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceOutOfTerritory):
		log.Info("place is out of organization's territory")
		c.responser.RenderErr(w, problems.Forbidden("place is out of organization's territory"))
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("place class not found")
		c.responser.RenderErr(w, problems.NotFound("place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.Info("place class is deprecated")
		c.responser.RenderErr(w, problems.Conflict("place class is deprecated"))
	case err != nil:
		log.WithError(err).Error("failed to create place")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.WithField("place_id", res.ID).Info("place created")
		c.responser.Render(w, http.StatusCreated, responses.Place(res))
	}
}
