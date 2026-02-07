package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place id"))...)

		return
	}

	res, err := c.core.place.Get(r.Context(), placeID)
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotExists):
			c.responser.RenderErr(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, res)
}
