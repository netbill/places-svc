package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationDeprecatePlaceClass = "deprecate_place_class"

func (c *Controller) DeprecatePlaceClass(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeprecatePlaceClass)

	req, err := requests.DeprecatedPlaceClass(r)
	if err != nil {
		log.WithError(err).Info("invalid request body")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	class, err := c.modules.Class.Deprecate(r.Context(), req.Data.Id, req.Data.Attributes.Deprecated)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.Info("Place class not found")
		render.ResponseError(w, problems.NotFound("Place class not found"))
	case err != nil:
		log.WithError(err).Error("failed to deprecate Place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("Place class deprecated successfully")
		render.Response(w, http.StatusOK, responses.PlaceClass(class))
	}
}
