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
)

func (c Controller) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	req, err := requests.UpdatePlace(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place request data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.UpdatePlace(r.Context(), initiator.ID, req.Data.Id, place.UpdateParams{
		ClassID:     req.Data.Attributes.ClassId,
		Address:     req.Data.Attributes.Address,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Icon:        req.Data.Attributes.Icon,
		Banner:      req.Data.Attributes.Banner,
		Website:     req.Data.Attributes.Website,
		Phone:       req.Data.Attributes.Phone,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			ape.RenderErr(w, problems.NotFound("place class not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
