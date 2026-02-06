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
	"github.com/paulmach/orb"
)

func (c Controller) CreatePlace(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := requests.CreatePlace(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create places request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.place.CreatePlace(r.Context(), initiator, place.CreateParams{
		OrganizationID: req.Data.Attributes.OrganizationId,
		ClassID:        req.Data.Attributes.ClassId,
		Address:        req.Data.Attributes.Address,
		Name:           req.Data.Attributes.Name,
		Description:    req.Data.Attributes.Description,
		//Icon:           req.Data.Attributes.Icon,
		//Banner:         req.Data.Attributes.Banner,
		Website: req.Data.Attributes.Website,
		Phone:   req.Data.Attributes.Phone,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create place")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to create place"))
		case errors.Is(err, errx.ErrorPlaceOutOfTerritory):
			c.responser.RenderErr(w, problems.Forbidden("place is out of organization's territory"))
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusCreated, responses.Place(res))
}
