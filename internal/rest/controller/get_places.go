package controller

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
	"github.com/paulmach/orb"
)

const operationGetPlaces = "get_places"

func (c *Controller) GetPlaces(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaces)

	limit, offset := pagi.GetPagination(r)
	if limit > 1000 {
		log.WithField("limit", limit).Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	params := place.FilterParams{}

	if orgIDStr := r.URL.Query().Get("organization_id"); orgIDStr != "" {
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			log.WithError(err).Info("invalid organization id")
			render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
			return
		}
		params.OrganizationID = &orgID
	}

	if statuses, ok := r.URL.Query()["statuses"]; ok {
		params.Status = statuses
	}

	if statuses, ok := r.URL.Query()["org_status"]; ok {
		params.OrgStatus = statuses
	}

	if verified := r.URL.Query().Get("verified"); verified != "" {
		value, err := strconv.ParseBool(verified)
		if err != nil {
			log.WithError(err).Info("invalid verified flag")
			render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid verified flag"))...)
			return
		}
		params.Verified = &value
	}

	if text := r.URL.Query().Get("text"); text != "" {
		params.BestMatch = &text
	}

	if classIDs, ok := r.URL.Query()["class_ids"]; ok {
		ids := make([]uuid.UUID, 0, len(classIDs))
		for _, classID := range classIDs {
			id, err := uuid.Parse(classID)
			if err != nil {
				log.WithError(err).WithField("class_id", classID).Info("invalid class_id")
				render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid class_id"))...)
				return
			}
			ids = append(ids, id)
		}

		class := &place.FilterClassParams{ClassID: ids}

		if incParentStr := r.URL.Query().Get("include_parent"); incParentStr != "" {
			value, err := strconv.ParseBool(incParentStr)
			if err != nil {
				log.WithError(err).Info("invalid include_parent flag")
				render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid include_parent"))...)
				return
			}
			class.Parents = value
		}

		if includeChildStr := r.URL.Query().Get("include_children"); includeChildStr != "" {
			value, err := strconv.ParseBool(includeChildStr)
			if err != nil {
				log.WithError(err).Info("invalid include_children flag")
				render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid include_children"))...)
				return
			}
			class.Children = value
		}

		params.Class = class
	}

	if lonStr := r.URL.Query().Get("lon"); lonStr != "" {
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			log.WithError(err).Info("invalid lon")
			render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid lon"))...)
			return
		}

		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			log.WithError(err).Info("invalid lat")
			render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid lat"))...)
			return
		}

		radius, err := strconv.ParseUint(r.URL.Query().Get("radius"), 10, 64)
		if err != nil {
			log.WithError(err).Info("invalid radius")
			render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid radius"))...)
			return
		}

		params.Near = &place.FilterNearParams{
			Point:   orb.Point{lon, lat},
			RadiusM: uint(radius),
		}
	}

	places, err := c.modules.Place.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get places")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.PlaceCollectionOption, 0)

	if slices.Contains(includes, "place_class") {
		classIDs := make([]uuid.UUID, 0, places.Size)
		for _, p := range places.Data {
			classIDs = append(classIDs, p.ClassID)
		}

		classes, err := c.modules.Class.GetByIDs(r.Context(), classIDs)
		if err != nil {
			log.WithError(err).Error("failed to get place classes")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionClass(classes))
	}

	if slices.Contains(includes, "organization") {
		orgIDs := make([]uuid.UUID, 0, places.Size)
		for _, p := range places.Data {
			orgIDs = append(orgIDs, p.OrganizationID)
		}

		orgs, err := c.modules.Org.GetByIDs(r.Context(), orgIDs)
		if err != nil {
			log.WithError(err).Error("failed to get organizations")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionOrganization(orgs))
	}

	render.Response(w, http.StatusOK, responses.Places(r, places, opts...))
}
