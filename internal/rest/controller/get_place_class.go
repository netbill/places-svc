package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetPlaceClass(w http.ResponseWriter, r *http.Request) {
	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place class id"))...)

		return
	}

	res, err := c.core.pclass.Get(r.Context(), classID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, http.StatusOK, responses.PlaceClass(res))
}
