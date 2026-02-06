package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/restkit/pagi"
	"github.com/paulmach/orb"
)

type PlaceRow struct {
	ID             uuid.UUID  `json:"id"`
	ClassID        uuid.UUID  `json:"class_id"`
	OrganizationID *uuid.UUID `json:"organization_id"`

	Status   string `json:"status"`
	Verified bool   `json:"verified"`

	Point   orb.Point `json:"point"`
	Address string    `json:"address"`

	Name        string  `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Banner      *string `json:"banner"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r PlaceRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r PlaceRow) ToModel() models.Place {
	return models.Place{
		ID:             r.ID,
		ClassID:        r.ClassID,
		OrganizationID: r.OrganizationID,
		Status:         r.Status,
		Point:          r.Point,
		Verified:       r.Verified,
		Address:        r.Address,
		Name:           r.Name,
		Description:    r.Description,
		Icon:           r.Icon,
		Banner:         r.Banner,
		Website:        r.Website,
		Phone:          r.Phone,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type PlacesQ interface {
	New() PlacesQ
	Insert(ctx context.Context, input PlaceRow) (PlaceRow, error)

	Get(ctx context.Context) (PlaceRow, error)
	Select(ctx context.Context) ([]PlaceRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (PlaceRow, error)

	UpdateClassID(classID uuid.UUID) PlacesQ
	UpdateStatus(status string) PlacesQ
	UpdateVerified(verified bool) PlacesQ
	UpdateName(name string) PlacesQ
	UpdateAddress(address string) PlacesQ
	UpdateDescription(description *string) PlacesQ
	UpdateIcon(icon *string) PlacesQ
	UpdateBanner(banner *string) PlacesQ
	UpdateWebsite(website *string) PlacesQ
	UpdatePhone(phone *string) PlacesQ

	FilterByID(id uuid.UUID) PlacesQ
	FilterByClassID(children, parents bool, classIDs ...uuid.UUID) PlacesQ
	FilterByOrganizationID(organizationID *uuid.UUID) PlacesQ
	FilterByText(text string) PlacesQ
	FilterByParentID(parentID uuid.UUID) PlacesQ
	FilterByRadius(point orb.Point, radiusM uint) PlacesQ
	FilterLikeName(name string) PlacesQ
	FilterLikeDescription(description string) PlacesQ
	FilterByStatus(status ...string) PlacesQ
	FilterByVerified(verified bool) PlacesQ
	FilterLikeAddress(address string) PlacesQ

	Delete(ctx context.Context) error

	Page(limit, offset uint) PlacesQ
	Count(ctx context.Context) (uint, error)
}

func (r *Repository) CreatePlace(ctx context.Context, params place.CreateParams) (models.Place, error) {
	row, err := r.PlacesQ.New().Insert(ctx, PlaceRow{
		OrganizationID: params.OrganizationID,
		ClassID:        params.ClassID,
		Status:         models.PlaceStatusInactive,
		Verified:       false,
		Point:          params.Point,
		Address:        params.Address,
		Name:           params.Name,
		Description:    params.Description,
		Website:        params.Website,
		Phone:          params.Phone,
	})
	if err != nil {
		return models.Place{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error) {
	row, err := r.PlacesQ.New().FilterByID(id).Get(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) GetPlaces(
	ctx context.Context,
	params place.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Place], error) {
	q := r.PlacesQ.New()

	if params.OrganizationID != nil {
		q = q.FilterByOrganizationID(params.OrganizationID)
	}
	if params.Status != nil {
		q = q.FilterByStatus(params.Status...)
	}
	if params.BestMatch != nil {
		q = q.FilterByText(*params.BestMatch)
	}
	if params.Verified != nil {
		q = q.FilterByVerified(*params.Verified)
	}
	if params.Address != nil {
		q = q.FilterLikeAddress(*params.Address)
	}
	if params.Name != nil {
		q = q.FilterLikeName(*params.Name)
	}
	if params.Description != nil {
		q = q.FilterLikeDescription(*params.Description)
	}
	if params.Class != nil {
		q = q.FilterByClassID(params.Class.Children, params.Class.Parents, params.Class.ClassID...)
	}
	if params.Near != nil {
		q = q.FilterByRadius(params.Near.Point, params.Near.RadiusM)
	}

	res, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Place]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Place]{}, err
	}

	collection := make([]models.Place, 0, len(res))
	for _, row := range res {
		collection = append(collection, row.ToModel())
	}

	return pagi.Page[[]models.Place]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil

}

func (r *Repository) CheckPlaceClassExists(ctx context.Context, classID uuid.UUID) (bool, error) {
	exists, err := r.PlacesQ.New().FilterByClassID(false, false, classID).Exists(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) UpdatePlaceByID(
	ctx context.Context,
	placeID uuid.UUID,
	params place.UpdateParams,
) (models.Place, error) {
	row, err := r.PlacesQ.New().
		FilterByID(placeID).
		UpdateClassID(params.ClassID).
		UpdateAddress(params.Address).
		UpdateName(params.Name).
		UpdateDescription(params.Description).
		UpdateIcon(params.GetUpdatedIcon()).
		UpdateBanner(params.GetUpdatedBanner()).
		UpdateWebsite(params.Website).
		UpdatePhone(params.Phone).
		UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdatePlaceStatus(
	ctx context.Context,
	placeID uuid.UUID,
	status string,
) (models.Place, error) {
	row, err := r.PlacesQ.New().FilterByID(placeID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdatePlaceVerified(
	ctx context.Context,
	placeID uuid.UUID,
	verified bool,
) (models.Place, error) {
	row, err := r.PlacesQ.New().FilterByID(placeID).UpdateVerified(verified).UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdatePlaceStatusForOrg(
	ctx context.Context,
	organizationID uuid.UUID,
	status string,
) error {
	if _, err := r.PlacesQ.New().FilterByOrganizationID(&organizationID).UpdateStatus(status).UpdateOne(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReplacePlacesClassID(
	ctx context.Context,
	oldClassID, newClassID uuid.UUID,
) error {
	_, err := r.PlacesQ.New().
		FilterByClassID(false, false, oldClassID).
		UpdateClassID(newClassID).
		UpdateMany(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePlaceByID(ctx context.Context, placeID uuid.UUID) error {
	return r.PlacesQ.New().FilterByID(placeID).Delete(ctx)
}
