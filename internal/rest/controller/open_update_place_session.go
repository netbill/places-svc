package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) OpenUpdatePlaceSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)

		return
	}

	place, media, err := c.core.place.OpenUpdateSession(
		r.Context(),
		initiator,
		placeID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get preload link for update place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotExists):
			c.responser.RenderErr(w, problems.NotFound("place does not exist"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, 200, responses.OpenUpdatePlaceSession(place, media))
}
