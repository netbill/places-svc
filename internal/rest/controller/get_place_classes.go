package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

const operationGetPlaceClasses = "get_place_classes"

func (c *Controller) GetPlaceClasses(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClasses)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.WithField("limit", limit).Info("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	params := pclass.FilterParams{}

	if text := r.URL.Query().Get("text"); text != "" {
		params.BestMatch = &text
	}

	if _, ok := r.URL.Query()["parent_id"]; ok {
		parentID, err := uuid.Parse(r.URL.Query().Get("parent_id"))
		if err != nil {
			log.WithError(err).Info("invalid parent id")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid parent id"))...)
			return
		}

		parent := &pclass.FilterPlaceClassParams{ID: parentID}

		if raw := r.URL.Query().Get("with_parents"); raw != "" {
			with, err := strconv.ParseBool(raw)
			if err != nil {
				log.WithError(err).Info("invalid with_parents value")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid with_parents value"))...)
				return
			}
			parent.IncludedParents = with
		}

		if raw := r.URL.Query().Get("with_child"); raw != "" {
			with, err := strconv.ParseBool(raw)
			if err != nil {
				log.WithError(err).Info("invalid with_child value")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid with_child value"))...)
				return
			}
			parent.IncludedChildren = with
		}

		params.Parent = parent
	}

	if deprecated := r.URL.Query().Get("deprecated"); deprecated != "" {
		with, err := strconv.ParseBool(deprecated)
		if err != nil {
			log.WithError(err).Info("invalid deprecated value")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid deprecated value"))...)
			return
		}
		params.Deprecated = &with
	}

	res, err := c.modules.pclass.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get place classes")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.PlaceClasses(r, res))
	}
}
