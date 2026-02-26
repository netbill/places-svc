package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetPlaceClass = "get_place_class"

func (c *Controller) GetPlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Info("invalid Place class id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid Place class id"))...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.modules.Class.Get(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("Place class not found")
		c.responser.RenderErr(w, problems.NotFound("Place class not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get Place class")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.PlaceClassOption, 0)

	if slices.Contains(includes, "parent") && class.ParentID != nil {
		parent, err := c.modules.Class.Get(r.Context(), *class.ParentID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("parent_class_id", class.ParentID).Info("parent Place class not found")
			c.responser.RenderErr(w, problems.NotFound("parent Place class not found"))
			return
		case err != nil:
			log.WithError(err).Error("failed to get parent Place class")
			c.responser.RenderErr(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithParentClass(parent))
		}
	}

	c.responser.Render(w, http.StatusOK, responses.PlaceClass(class, opts...))
}
