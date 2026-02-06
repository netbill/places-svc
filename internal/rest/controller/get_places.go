package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetPlaces(w http.ResponseWriter, r *http.Request) {
	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	params := place.FilterParams{}

	if orgIDStr := r.URL.Query().Get("organization_id"); orgIDStr != "" {
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
			return
		}

		params.OrganizationID = &orgID
	}

	if statuses, ok := r.URL.Query()["statuses"]; ok {
		params.Status = statuses
	}

	if verified := r.URL.Query().Get("verified"); verified != "" {
		value, err := strconv.ParseBool(verified)
		if err != nil {
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid verified flag"))...)
			return
		}

		params.Verified = &value
	}

	if _, ok := r.URL.Query()["text"]; ok {
		text := r.URL.Query().Get("text")
		params.BestMatch = &text
	}

	if classIDs, ok := r.URL.Query()["class_ids"]; ok {
		ids := make([]uuid.UUID, 0, len(classIDs))

		for _, classID := range classIDs {
			id, err := uuid.Parse(classID)
			if err != nil {
				c.log.WithError(err).Errorf("invalid class_id %s", classID)
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid class_id"))...)
				return
			}
			ids = append(ids, id)
		}

		if incParentStr := r.URL.Query().Get("include_parent"); incParentStr != "" {
			value, err := strconv.ParseBool(incParentStr)
			if err != nil {
				c.log.WithError(err).Errorf("invalid include_parent flag")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid include_parent"))...)
				return
			}

			params.Class.Parents = value
		}

		if includeChildStr := r.URL.Query().Get("include_children"); includeChildStr != "" {
			value, err := strconv.ParseBool(includeChildStr)
			if err != nil {
				c.log.WithError(err).Errorf("invalid include_children")
				c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid include_children"))...)
				return
			}

			params.Class.Children = value
		}

		params.Class.ClassID = ids
	}

	if lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64); err == nil {
		params.Near.Point[0] = lon

		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid lat"))...)
			return
		}
		params.Near.Point[1] = lat

		radius, err := strconv.ParseUint(r.URL.Query().Get("radius"), 10, 64)
		if err != nil {
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid radius"))...)
			return
		}
		params.Near.RadiusM = uint(radius)
	}

	res, err := c.core.place.GetList(r.Context(), params, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("error getting places")
		c.responser.RenderErr(w, problems.InternalError())
	}

	c.responser.Render(w, http.StatusOK, responses.Places(r, res))
}
