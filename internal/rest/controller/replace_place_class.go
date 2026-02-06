package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) ReplacePlaceClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.ReplacePlaceClass(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid replace place class request")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid replace place class request: %s", err))...)
		return
	}

	err = c.core.class.Replace(r.Context(), req.Data.Id, req.Data.Attributes.ClassReplaceId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to replace place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK)
}
