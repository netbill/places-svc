package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/places"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
	"github.com/paulmach/orb"
)

const operationCreatePlace = "create_place"

type placeCore interface {
	Create(
		ctx context.Context,
		actor models.AccountActor,
		params places.CreateParams,
	) (place models.Place, err error)

	Delete(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
	) error

	Get(ctx context.Context, placeID uuid.UUID) (models.Place, error)

	GetByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.Place, error)

	GetList(
		ctx context.Context,
		params places.FilterPlaceParams,
		limit, offset uint,
	) (pagi.Page[[]models.Place], error)

	Update(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		params places.UpdateParams,
	) (place models.Place, err error)

	Activate(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
	) (place models.Place, err error)

	Deactivate(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
	) (place models.Place, err error)

	Verify(
		ctx context.Context,
		placeID uuid.UUID,
	) (place models.Place, err error)

	Unverify(
		ctx context.Context,
		placeID uuid.UUID,
	) (place models.Place, err error)

	CreateUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
	) (models.Place, models.UploadPlaceMediaLinks, error)

	DeleteUploadPlaceMedia(
		ctx context.Context,
		actor models.AccountActor,
		placeID uuid.UUID,
		params places.DeleteUploadPlaceMediaParams,
	) error
}

type organizationGetter interface {
	Get(ctx context.Context, id uuid.UUID) (models.Organization, error)

	GetByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]models.Organization, error)
}

type placeClassGetter interface {
	Get(ctx context.Context, id uuid.UUID) (models.PlaceClass, error)

	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.PlaceClass, error)
}

type PlaceController struct {
	place placeCore
	class placeClassGetter
	org   organizationGetter
}

type PlaceControllerDeps struct {
	Place placeCore
	Class placeClassGetter
	Org   organizationGetter
}

func NewPlaceController(deps PlaceControllerDeps) *PlaceController {
	return &PlaceController{
		place: deps.Place,
		class: deps.Class,
		org:   deps.Org,
	}
}

func (c *PlaceController) Create(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreatePlace)

	req, err := requests.CreatePlace(r)
	if err != nil {
		log.WithError(err).Warn("invalid create place request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId).
		WithField("class_id", req.Data.Attributes.ClassId)

	res, err := c.place.Create(r.Context(), scope.AccountActor(r), places.CreateParams{
		OrganizationID: req.Data.Attributes.OrganizationId,
		ClassID:        req.Data.Attributes.ClassId,
		Address:        req.Data.Attributes.Address,
		Name:           req.Data.Attributes.Name,
		Description:    req.Data.Attributes.Description,
		Website:        req.Data.Attributes.Website,
		Phone:          req.Data.Attributes.Phone,
		Point: orb.Point{
			req.Data.Attributes.Point.Longitude,
			req.Data.Attributes.Point.Latitude,
		},
	})
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("account is not organization head")
		render.ResponseError(w, problems.Forbidden("account is not organization head"))
	case errors.Is(err, errx.ErrorPlaceOutOfTerritory):
		log.WithError(err).Warn("place is out of organization's territory")
		render.ResponseError(w, problems.Forbidden("place is out of organization's territory"))
	case errors.Is(err, errx.ErrorPlaceClassNotExists):
		log.WithError(err).Warn("place class not found")
		render.ResponseError(w, problems.NotFound("place class not found"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithError(err).Warn("place class is deprecated")
		render.ResponseError(w, problems.Conflict("place class is deprecated"))
	case err != nil:
		log.WithError(err).Error("failed to create place")
		render.ResponseError(w, problems.InternalError())
	default:
		log.WithField("place_id", res.ID).Info("place created")
		render.Response(w, http.StatusCreated, responses.Place(r, res))
	}
}

const operationDeletePlace = "delete_place"

func (c *PlaceController) Delete(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeletePlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	err = c.place.Delete(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists),
		errors.Is(err, errx.ErrorPlaceDeleted):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("account is not organization head")
		render.ResponseError(w, problems.Forbidden("account is not organization head"))
	case err != nil:
		log.WithError(err).Error("failed to delete place")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, nil)
	}
}

const operationGetPlace = "get_place"

func (c *PlaceController) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlace)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Warn("invalid place_id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place_id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	place, err := c.place.Get(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists),
		errors.Is(err, errx.ErrorPlaceDeleted):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		return
	case err != nil:
		log.WithError(err).Error("failed to get place")
		render.ResponseError(w, problems.InternalError())
		return
	}

	slog.Debug(fmt.Sprintf("place found: %s", place.Name))

	opts := make([]responses.PlaceOption, 0)
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "place_class") {
		class, err := c.class.Get(r.Context(), place.ClassID)
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotExists):
			log.WithField("place_class_id", place.ClassID).Warn("place class not found")
			render.ResponseError(w, problems.NotFound(fmt.Sprintf("place class with id %s not found", place.ClassID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get place class")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithClass(r, class))
		}
	}

	if slices.Contains(includes, "organization") {
		org, err := c.org.Get(r.Context(), place.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists),
			errors.Is(err, errx.ErrorOrganizationDeleted):
			log.WithField("organization_id", place.OrganizationID).Warn("organization not found")
			render.ResponseError(w, problems.NotFound(fmt.Sprintf("organization with id %s not found", place.OrganizationID)))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithOrganization(r, org))
		}
	}

	render.Response(w, http.StatusOK, responses.Place(r, place, opts...))
}

const operationGetPlaces = "get_places"

func (c *PlaceController) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetPlaces)

	limit, offset := pagi.GetPagination(r)
	if limit > 1000 {
		log.WithField("limit", limit).Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/size": fmt.Errorf("pagination limit cannot be greater than 1000"),
		})...)
		return
	}

	params := places.FilterPlaceParams{}

	if orgIDStr := r.URL.Query().Get("organization_id"); orgIDStr != "" {
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			log.WithError(err).Warn("invalid organization id")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/organization_id": fmt.Errorf("invalid organization id: %s", orgIDStr),
			})...)
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
			log.WithError(err).Warn("invalid verified flag")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/verified": fmt.Errorf("invalid verified flag: %s", verified),
			})...)
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
				log.WithError(err).WithField("class_id", classID).Warn("invalid class_id")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query/class_ids": fmt.Errorf("invalid class_id: %s", classID),
				})...)
				return
			}
			ids = append(ids, id)
		}

		class := &places.FilterClassForPlaceParams{ClassID: ids}

		if incParentStr := r.URL.Query().Get("include_parent"); incParentStr != "" {
			value, err := strconv.ParseBool(incParentStr)
			if err != nil {
				log.WithError(err).Warn("invalid include_parent flag")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query/include_parent": fmt.Errorf("invalid include_parent flag: %s", incParentStr),
				})...)
				return
			}
			class.Parents = value
		}

		if includeChildStr := r.URL.Query().Get("include_children"); includeChildStr != "" {
			value, err := strconv.ParseBool(includeChildStr)
			if err != nil {
				log.WithError(err).Warn("invalid include_children flag")
				render.ResponseError(w, problems.BadRequest(validation.Errors{
					"query/include_children": fmt.Errorf("invalid include_children flag: %s", includeChildStr),
				})...)
				return
			}
			class.Children = value
		}

		params.Class = class
	}

	if lonStr := r.URL.Query().Get("lon"); lonStr != "" {
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			log.WithError(err).Warn("invalid lon")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/lon": fmt.Errorf("invalid lon"),
			})...)
			return
		}

		lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			log.WithError(err).Warn("invalid lat")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/lat": fmt.Errorf("invalid lat"),
			})...)
			return
		}

		radius, err := strconv.ParseUint(r.URL.Query().Get("radius"), 10, 64)
		if err != nil {
			log.WithError(err).Warn("invalid radius")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/radius": fmt.Errorf("invalid radius"),
			})...)
			return
		}

		params.Near = &places.FilterPlaceNearParams{
			Point:   orb.Point{lon, lat},
			RadiusM: uint(radius),
		}
	}

	places, err := c.place.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get places")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.PlaceCollectionOption, 0)
	includes := restkit.ParseIncludes(r)

	log.WithField("includes", includes).Info("parsed includes")

	if slices.Contains(includes, "place_classes") {
		classIDs := make([]uuid.UUID, 0, places.Size)
		for _, p := range places.Data {
			classIDs = append(classIDs, p.ClassID)
		}

		classes, err := c.class.GetByIDs(r.Context(), classIDs)
		if err != nil {
			log.WithError(err).Error("failed to get place classes")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionClass(r, classes))
	}

	if slices.Contains(includes, "organizations") {
		orgIDs := make([]uuid.UUID, 0, places.Size)
		for _, p := range places.Data {
			orgIDs = append(orgIDs, p.OrganizationID)
		}

		orgs, err := c.org.GetByIDs(r.Context(), orgIDs)
		if err != nil {
			log.WithError(err).Error("failed to get organizations")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionOrganization(r, orgs))
	}

	render.Response(w, http.StatusOK, responses.Places(r, places, opts...))
}
