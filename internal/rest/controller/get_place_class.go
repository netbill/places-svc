package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetPlaceClass = "get_place_class"

func (c *Controller) GetPlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid place class id: %s", chi.URLParam(r, "place_class_id")),
		})...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.modules.Class.Get(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get place class")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.PlaceClassOption, 0)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 1)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	if slices.Contains(includes, "parent") && class.ParentID != nil {
		parent, err := c.modules.Class.Get(r.Context(), *class.ParentID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("parent_class_id", class.ParentID).Warn("place class not found")
			render.ResponseError(w, problems.NotFound("place class not found"))
			return
		case err != nil:
			log.WithError(err).Error("failed to get place class")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithParentClass(parent))
		}
	}

	render.Response(w, http.StatusOK, responses.PlaceClass(class, opts...))
}
