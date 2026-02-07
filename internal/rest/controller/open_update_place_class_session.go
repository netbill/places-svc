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

func (c *Controller) OpenUpdatePlaceClassSession(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get user from context")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	placeClassID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class id")
		c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place class id: %s", chi.URLParam(r, "place_class_id")),
		})...)

		return
	}

	placeClass, media, err := c.core.pclass.OpenUpdateSession(
		r.Context(),
		initiator,
		placeClassID,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get preload link for update place_class")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			c.responser.RenderErr(w, problems.NotFound("place_class does not exist"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}

		return
	}

	c.responser.Render(w, 200, responses.OpenUpdatePlaceClassSession(placeClass, media))
}
