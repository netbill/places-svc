package controller

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetPlaceClasses = "get_place_classes"

func (c *Controller) GetPlaceClasses(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClasses)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.WithField("limit", limit).Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	params := pclass.FilterParams{}

	if text := r.URL.Query().Get("text"); text != "" {
		params.BestMatch = &text
	}

	if _, ok := r.URL.Query()["parent_id"]; ok {
		parentID, err := uuid.Parse(r.URL.Query().Get("parent_id"))
		if err != nil {
			log.WithError(err).Warn("invalid parent id")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("invalid parent id: %s", r.URL.Query().Get("parent_id")),
			})...)
			return
		}

		parent := &pclass.FilterPlaceClassParams{ID: parentID}

		if raw := r.URL.Query().Get("with_parents"); raw != "" {
			with, err := strconv.ParseBool(raw)
			if err != nil {
				log.WithError(err).Warn("invalid with_parents value")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query": fmt.Errorf("invalid with_parents value"),
				})...)
				return
			}
			parent.IncludedParents = with
		}

		if raw := r.URL.Query().Get("with_child"); raw != "" {
			with, err := strconv.ParseBool(raw)
			if err != nil {
				log.WithError(err).Warn("invalid with_child value")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query": fmt.Errorf("invalid with_child value"),
				})...)
				return
			}
			parent.IncludedChildren = with
		}

		params.Parent = parent
	}

	if deprecated := r.URL.Query().Get("deprecated"); deprecated != "" {
		with, err := strconv.ParseBool(deprecated)
		if err != nil {
			log.WithError(err).Warn("invalid deprecated value")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("invalid deprecated value"),
			})...)
			return
		}
		params.Deprecated = &with
	}

	classes, err := c.modules.Class.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get place classes")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.PlaceClassCollectionOption, 0)
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

	if slices.Contains(includes, "parents") {
		parentIDs := make([]uuid.UUID, 0, classes.Size)
		for _, p := range classes.Data {
			if p.ParentID != nil {
				parentIDs = append(parentIDs, *p.ParentID)
			}
		}

		parents, err := c.modules.Class.GetByIDs(r.Context(), parentIDs)
		if err != nil {
			log.WithError(err).Error("failed to get place classes")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionParentClass(parents))
	}

	render.Response(w, http.StatusOK, responses.PlaceClasses(r, classes, opts...))
}
