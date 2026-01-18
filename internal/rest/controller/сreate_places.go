package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/paulmach/orb"
)

func (c Controller) CreatePlace(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := requests.CreatePlace(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create places request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.CreatePlace(r.Context(), initiator.ID, place.CreateParams{
		OrganizationID: req.Data.Attributes.OrganizationId,
		ClassID:        req.Data.Attributes.ClassId,
		Address:        req.Data.Attributes.Address,
		Name:           req.Data.Attributes.Name,
		Description:    req.Data.Attributes.Description,
		Icon:           req.Data.Attributes.Icon,
		Banner:         req.Data.Attributes.Banner,
		Website:        req.Data.Attributes.Website,
		Phone:          req.Data.Attributes.Phone,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to create place")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to create place"))
		case errors.Is(err, errx.ErrorPlaceOutOfTerritory):
			ape.RenderErr(w, problems.Forbidden("place is out of organization's territory"))
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			ape.RenderErr(w, problems.NotFound("place class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Place(res))
}
