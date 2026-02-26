package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationDeprecatePlaceClass = "deprecate_place_class"

func (c *Controller) DeprecatePlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeprecatePlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Info("invalid place class id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid place class id"))...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.modules.pclass.Deprecate(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("place class not found")
		c.responser.RenderErr(w, problems.NotFound("place class not found"))
	case err != nil:
		log.WithError(err).Error("failed to deprecate place class")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		log.Info("place class deprecated successfully")
		c.responser.Render(w, http.StatusOK, responses.PlaceClass(class))
	}
}
