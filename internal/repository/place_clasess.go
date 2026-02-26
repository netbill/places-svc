package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/restkit/pagi"
)

type PlaceClassRow struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`

	Name        string  `json:"name"`
	Description string  `json:"description"`
	IconKey     *string `json:"icon_key,omitempty"`

	Version      int32      `json:"version"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeprecatedAt *time.Time `json:"deprecated_at,omitempty"`
}

func (r PlaceClassRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r PlaceClassRow) ToModel() models.PlaceClass {
	return models.PlaceClass{
		ID:           r.ID,
		ParentID:     r.ParentID,
		Name:         r.Name,
		Description:  r.Description,
		IconKey:      r.IconKey,
		Version:      r.Version,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		DeprecatedAt: r.DeprecatedAt,
	}
}

type PlaceClassesQ interface {
	New() PlaceClassesQ
	Insert(ctx context.Context, input PlaceClassRow) (PlaceClassRow, error)

	Get(ctx context.Context) (PlaceClassRow, error)
	Select(ctx context.Context) ([]PlaceClassRow, error)

	Exists(ctx context.Context) (bool, error)
	Page(limit, offset uint) PlaceClassesQ
	Count(ctx context.Context) (uint, error)

	UpdateOne(ctx context.Context) (PlaceClassRow, error)

	UpdateParent(parentID *uuid.UUID) PlaceClassesQ
	UpdateName(name string) PlaceClassesQ
	UpdateDescription(description string) PlaceClassesQ
	UpdateIconKey(key *string) PlaceClassesQ
	UpdateDeprecatedAt(time *time.Time) PlaceClassesQ

	FilterByID(id ...uuid.UUID) PlaceClassesQ
	FilterByParentID(parentID ...uuid.UUID) PlaceClassesQ
	FilterBestMatch(text string) PlaceClassesQ
	FilterByClassID(classID uuid.UUID, includeChild, includeParent bool) PlaceClassesQ

	OrderName(asc bool) PlaceClassesQ
	OrderRoot(asc bool) PlaceClassesQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreatePlaceClass(ctx context.Context, params pclass.CreateParams) (models.PlaceClass, error) {
	row, err := r.PlaceClassesSql.New().Insert(ctx, PlaceClassRow{
		ParentID:    params.ParentID,
		Name:        params.Name,
		Description: params.Description,
		IconKey:     params.IconKey,
	})
	if err != nil {
		return models.PlaceClass{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) PlaceClassExists(ctx context.Context, id uuid.UUID) (bool, error) {
	res, err := r.PlaceClassesSql.New().FilterByID(id).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (r *Repository) GetPlaceClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	row, err := r.PlaceClassesSql.New().FilterByID(id).Get(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) GetPlaceClassesByIDs(ctx context.Context, ids []uuid.UUID) ([]models.PlaceClass, error) {
	rows, err := r.PlaceClassesSql.New().FilterByID(ids...).Select(ctx)
	if err != nil {
		return nil, err
	}

	collection := make([]models.PlaceClass, len(rows))
	for i, row := range rows {
		collection[i] = row.ToModel()
	}

	return collection, nil
}

func (r *Repository) GetPlaceClasses(
	ctx context.Context,
	params pclass.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.PlaceClass], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}

	q := r.PlaceClassesSql.New()

	if params.BestMatch != nil {
		q = q.FilterBestMatch(*params.BestMatch)
	}
	if params.Parent != nil {
		q = q.FilterByClassID(params.Parent.ID, params.Parent.IncludedChildren, params.Parent.IncludedParents)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.PlaceClass]{}, err
	}

	collection := make([]models.PlaceClass, len(rows))
	for i, row := range rows {
		collection[i] = row.ToModel()
	}

	return pagi.Page[[]models.PlaceClass]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

func (r *Repository) UpdatePlaceClass(ctx context.Context, classID uuid.UUID, params pclass.UpdateParams) (models.PlaceClass, error) {
	row, err := r.PlaceClassesSql.New().
		FilterByID(classID).
		UpdateParent(params.ParentID).
		UpdateName(params.Name).
		UpdateDescription(params.Description).
		UpdateIconKey(params.IconKey).
		UpdateOne(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return row.ToModel(), nil
}

func (r *Repository) CheckParentCycle(ctx context.Context, classID, newParentID uuid.UUID) (bool, error) {
	res, err := r.PlaceClassesSql.New().FilterByClassID(newParentID, true, true).FilterByID(classID).Exists(ctx)
	if err != nil {
		return false, err
	}
	if res {
		return true, nil
	}

	return false, nil
}

func (r *Repository) CheckPlaceClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error) {
	res, err := r.PlaceClassesSql.New().FilterByClassID(classID, true, false).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (r *Repository) DeletePlaceClass(ctx context.Context, classID uuid.UUID) error {
	return r.PlaceClassesSql.New().FilterByID(classID).Delete(ctx)
}

func (r *Repository) DeprecatedPlaceClass(
	ctx context.Context,
	classID uuid.UUID,
	deprecate bool,
) (models.PlaceClass, error) {
	var deprecateAt *time.Time

	if deprecate {
		now := time.Now().UTC()
		deprecateAt = &now
	}

	row, err := r.PlaceClassesSql.New().
		FilterByID(classID).
		UpdateDeprecatedAt(deprecateAt).
		UpdateOne(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return row.ToModel(), nil
}
