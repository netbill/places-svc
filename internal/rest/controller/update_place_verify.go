package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
)

func (c Controller) UpdatePlaceVerify(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlaceVerify(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place verify data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.UpdatePlaceVerified(r.Context(), req.Data.Id, req.Data.Attributes.Verify)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place verify data")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
