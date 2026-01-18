package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
)

func (c Controller) GetPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place id"))...)
		return
	}

	res, err := c.core.GetPlace(r.Context(), placeID)
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	ape.Render(w, http.StatusOK, res)
}
