package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/restkit/problems"
)

func (c Controller) DeletePlaceClass(w http.ResponseWriter, r *http.Request) {
	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		c.log.WithError(err).Errorf("invalid place class id")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place class id"))...)
		return
	}

	err = c.core.DeletePlaceClass(r.Context(), classID)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassHaveChildren):
			ape.RenderErr(w, problems.Forbidden("cannot delete class because it has child classes"))
		case errors.Is(err, errx.ErrorPlacesExitsWithThisClass):
			ape.RenderErr(w, problems.Forbidden("cannot delete class when places exist with this class"))
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			ape.RenderErr(w, problems.NotFound("place class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK)
}
