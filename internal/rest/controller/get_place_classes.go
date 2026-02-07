package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetPlaceClasses(w http.ResponseWriter, r *http.Request) {
	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)

		return
	}

	params := pclass.FilterParams{}

	if _, ok := r.URL.Query()["text"]; ok {
		text := r.URL.Query().Get("text")
		params.BestMatch = &text
	}

	if _, ok := r.URL.Query()["parent_id"]; ok {
		parentID, err := uuid.Parse(r.URL.Query().Get("parent_id"))
		if err != nil {
			c.log.WithError(err).Errorf("invalid parent id")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid parent id"))...)

			return
		}
		params.Parent.ID = parentID

		if _, ok := r.URL.Query()["with_parents"]; ok {
			with, err := strconv.ParseBool(r.URL.Query().Get("with_parents"))
			if err != nil {
				c.log.WithError(err).Errorf("invalid with_parents value")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid with_parents value"))...)

				return
			}

			params.Parent.IncludedParents = with
		}

		if _, ok := r.URL.Query()["with_child"]; ok {
			with, err := strconv.ParseBool(r.URL.Query().Get("with_child"))
			if err != nil {
				c.log.WithError(err).Errorf("invalid with_child value")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid with_child value"))...)

				return
			}

			params.Parent.IncludedChildren = with
		}
	}

	res, err := c.core.pclass.GetList(r.Context(), params, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get place classes")
		c.responser.RenderErr(w, problems.InternalError())

		return
	}

	c.responser.Render(w, http.StatusOK, responses.PlaceClasses(r, res))
}
