package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/classification"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type placeClassCore interface {
	Create(
		ctx context.Context,
		params classification.CreateParams,
	) (class models.PlaceClass, err error)

	Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error)

	GetByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.PlaceClass, error)

	GetList(
		ctx context.Context,
		params classification.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.PlaceClass], error)

	Deprecate(
		ctx context.Context,
		classID uuid.UUID,
	) (models.PlaceClass, error)

	Undeprecate(
		ctx context.Context,
		classID uuid.UUID,
	) (models.PlaceClass, error)

	Update(
		ctx context.Context,
		classID uuid.UUID,
		params classification.UpdateParams,
	) (class models.PlaceClass, err error)

	CreateUploadMediaLinks(
		ctx context.Context,
		placeClassID uuid.UUID,
	) (models.PlaceClass, models.UploadPlaceClassMediaLinks, error)

	DeleteUploadMedia(
		ctx context.Context,
		classID uuid.UUID,
		params classification.DeleteUploadPlaceClassMediaParams,
	) error
}

type PlaceClassController struct {
	class placeClassCore
}

type PlaceClassControllerDeps struct {
	Class placeClassCore
}

func NewPlaceClassController(deps PlaceClassControllerDeps) *PlaceClassController {
	return &PlaceClassController{
		class: deps.Class,
	}
}

const operationCreatePlaceClass = "create_place_class"

func (c *PlaceClassController) Create(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlaceClass)

	req, err := requests.CreatePlaceClass(r)
	if err != nil {
		log.WithError(err).Warn("invalid create place class request")
		render.ResponseError(w, problems.BadRequest(err)...)

		return
	}

	res, err := c.class.Create(r.Context(), classification.CreateParams{
		ParentID:    req.Data.Attributes.ParentId,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
	})
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithField("place_class_id", req.Data.Attributes.ParentId).
			Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithField("place_class_id", req.Data.Attributes.ParentId).Warn("place class is deprecated")
		render.ResponseError(w, problems.Conflict("place class is deprecated"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to create place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.WithField("place_class_id", res.ID).Info("place class created")
		render.Response(w, http.StatusCreated, responses.PlaceClass(r, res))
	}
}

const operationDeprecatePlaceClass = "deprecate_place_class"

func (c *PlaceClassController) Deprecate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeprecatePlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_class_id": fmt.Errorf(
				"invalid place class id: %s", chi.URLParam(r, "place_class_id"),
			),
		})...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.class.Deprecate(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
	case err != nil:
		log.WithError(err).Error("failed to deprecate place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place class deprecated successfully")
		render.Response(w, http.StatusOK, responses.PlaceClass(r, class))
	}
}

const operationUndeprecatePlaceClass = "undeprecate_place_class"

func (c *PlaceClassController) Undeprecate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUndeprecatePlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_class_id": fmt.Errorf(
				"invalid place class id: %s", chi.URLParam(r, "place_class_id"),
			),
		})...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.class.Undeprecate(r.Context(), classID)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
	case err != nil:
		log.WithError(err).Error("failed to undeprecate place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place class undeprecated successfully")
		render.Response(w, http.StatusOK, responses.PlaceClass(r, class))
	}
}

const operationGetPlaceClass = "get_place_class"

func (c *PlaceClassController) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClass)

	classID, err := uuid.Parse(chi.URLParam(r, "place_class_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place class id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_class_id": fmt.Errorf(
				"invalid place class id: %s", chi.URLParam(r, "place_class_id"),
			),
		})...)
		return
	}

	log = log.WithField("place_class_id", classID)

	class, err := c.class.Get(r.Context(), classID)
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
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "parent") && class.ParentID != nil {
		parent, err := c.class.Get(r.Context(), *class.ParentID)
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
			opts = append(opts, responses.WithParentClass(r, parent))
		}
	}

	render.Response(w, http.StatusOK, responses.PlaceClass(r, class, opts...))
}

const operationGetPlaceClasses = "get_place_classes"

func (c *PlaceClassController) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaceClasses)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.WithField("limit", limit).Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/size": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	params := classification.FilterParams{}

	if text := r.URL.Query().Get("text"); text != "" {
		params.BestMatch = &text
	}

	if _, ok := r.URL.Query()["parent_id"]; ok {
		parentID, err := uuid.Parse(r.URL.Query().Get("parent_id"))
		if err != nil {
			log.WithError(err).Warn("invalid parent id")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/parent_id": fmt.Errorf("invalid parent id: %s", r.URL.Query().Get("parent_id")),
			})...)
			return
		}

		parent := &classification.FilterParentParams{ID: parentID}

		if raw := r.URL.Query().Get("with_parents"); raw != "" {
			with, err := strconv.ParseBool(raw)
			if err != nil {
				log.WithError(err).Warn("invalid with_parents value")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query/with_parents": fmt.Errorf("invalid with_parents value"),
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
					"query/with_child": fmt.Errorf("invalid with_child value"),
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
				"query/deprecated": fmt.Errorf("invalid deprecated value"),
			})...)
			return
		}
		params.Deprecated = &with
	}

	classes, err := c.class.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get place classes")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.PlaceClassCollectionOption, 0)
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "parents") {
		parentIDs := make([]uuid.UUID, 0, classes.Size)
		for _, p := range classes.Data {
			if p.ParentID != nil {
				parentIDs = append(parentIDs, *p.ParentID)
			}
		}

		parents, err := c.class.GetByIDs(r.Context(), parentIDs)
		if err != nil {
			log.WithError(err).Error("failed to get place classes")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionParentClass(r, parents))
	}

	render.Response(w, http.StatusOK, responses.PlaceClasses(r, classes, opts...))
}

const operationUpdatePlaceClass = "update_place_class"

func (c *PlaceClassController) Update(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceClass)

	req, err := requests.UpdatePlaceClass(r)
	if err != nil {
		log.WithError(err).Warn("invalid update place class request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_class_id", req.Data.Id)

	res, err := c.class.Update(
		r.Context(),
		req.Data.Id,
		classification.UpdateParams{
			ParentID:    req.Data.Attributes.ParentId,
			Name:        req.Data.Attributes.Name,
			Description: req.Data.Attributes.Description,
			IconKey:     req.Data.Attributes.IconKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassParentCycle):
		log.WithError(err).Warn("setting parent would create a cycle")
		render.ResponseError(w, problems.Conflict("setting parent would create a cycle"))
	case errors.Is(err, errx.ErrorPlaceClassIconIsInvalid):
		log.WithError(err).Warn("upload icon is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update place class")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place class updated")
		render.Response(w, http.StatusOK, responses.PlaceClass(r, res))
	}
}
