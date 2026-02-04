package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PlaceClassRow struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        *string    `json:"icon"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (r PlaceClassRow) IsNil() bool {
	return r.ID == uuid.Nil
}

type PlaceClassesQ interface {
	New() PlaceClassesQ
	Insert(ctx context.Context, input PlaceClassRow) (PlaceClassRow, error)

	Get(ctx context.Context) (PlaceClassRow, error)
	Select(ctx context.Context) ([]PlaceClassRow, error)

	Exists(ctx context.Context) (bool, error)
	Page(limit, offset uint) PlaceClassesQ
	Count(ctx context.Context) (uint, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (PlaceClassRow, error)

	UpdateParent(parentID *uuid.UUID) PlaceClassesQ
	UpdateCode(code string) PlaceClassesQ
	UpdateName(name string) PlaceClassesQ
	UpdateDescription(description string) PlaceClassesQ
	UpdateIcon(icon *string) PlaceClassesQ

	FilterByID(id ...uuid.UUID) PlaceClassesQ
	FilterByParentID(parentID ...uuid.UUID) PlaceClassesQ
	FilterRoots() PlaceClassesQ
	FilterByCode(code string) PlaceClassesQ
	FilterNameLike(name string) PlaceClassesQ
	FilterDescriptionLike(description string) PlaceClassesQ
	FilterByText(text string) PlaceClassesQ
	FilterByClassID(classID uuid.UUID, includeChild, includeParent bool) PlaceClassesQ

	OrderName(asc bool) PlaceClassesQ

	Delete(ctx context.Context) error
}

//func (r *Repository) CreatePlaceClass(ctx context.Context, params pclass.CreateParams) (models.PlaceClass, error) {
//	row, err := r.placeClassesQ(ctx).Insert(ctx, pg.PlaceClassesInsertInput{
//		ParentID:    params.ParentID,
//		Code:        params.Code,
//		Name:        params.Name,
//		Description: params.Description,
//		Icon:        params.Icon,
//	})
//	if err != nil {
//		return models.PlaceClass{}, err
//	}
//
//	return toModel(row), nil
//}
//
//func (r *Repository) GetPlaceClassByCode(ctx context.Context, code string) (models.PlaceClass, error) {
//	row, err := r.placeClassesQ(ctx).FilterByCode(code).Get(ctx)
//	if err != nil {
//		return models.PlaceClass{}, err
//	}
//
//	return toModel(row), nil
//}
//
//func (r *Repository) PlaceClassExists(ctx context.Context, id uuid.UUID) (bool, error) {
//	res, err := r.placeClassesQ(ctx).FilterByID(id).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) PlaceClassExistsByCode(ctx context.Context, code string) (bool, error) {
//	res, err := r.placeClassesQ(ctx).FilterByCode(code).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) GetPlaceClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
//	row, err := r.placeClassesQ(ctx).FilterByID(id).Get(ctx)
//	if err != nil {
//		return models.PlaceClass{}, err
//	}
//
//	return toModel(row), nil
//}
//
//func (r *Repository) GetPlaceClasses(ctx context.Context, params pclass.FilterParams, limit, offset uint) (pagi.Page[[]models.PlaceClass], error) {
//	q := r.placeClassesQ(ctx)
//
//	if params.Name != nil {
//		q = q.FilterNameLike(*params.Name)
//	}
//	if params.Description != nil {
//		q = q.FilterDescriptionLike(*params.Description)
//	}
//	if params.BestMatch != nil {
//		q = q.FilterByText(*params.BestMatch)
//	}
//	if params.Parent != nil {
//		q = q.FilterByClassID(params.Parent.ID, params.Parent.IncludedChildren, params.Parent.IncludedParents)
//	}
//
//	total, err := q.Count(ctx)
//	if err != nil {
//		return pagi.Page[[]models.PlaceClass]{}, err
//	}
//
//	rows, err := q.Page(limit, offset).Select(ctx)
//	if err != nil {
//		return pagi.Page[[]models.PlaceClass]{}, err
//	}
//
//	collection := make([]models.PlaceClass, len(rows))
//	for i, row := range rows {
//		collection[i] = toModel(row)
//	}
//
//	return pagi.Page[[]models.PlaceClass]{
//		Data:  collection,
//		Total: total,
//		Page:  uint(offset/limit) + 1,
//		Size:  uint(len(collection)),
//	}, nil
//}
//
//func (r *Repository) PlaceExists(ctx context.Context, id uuid.UUID) (bool, error) {
//	res, err := r.placeQ(ctx).FilterByID(id).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) UpdatePlaceClass(ctx context.Context, classID uuid.UUID, params pclass.UpdateParams) (models.PlaceClass, error) {
//	q := r.placeClassesQ(ctx).FilterByID(classID)
//	if params.Code != nil {
//		q = q.UpdateCode(*params.Code)
//	}
//	if params.Name != nil {
//		q = q.UpdateName(*params.Name)
//	}
//	if params.Description != nil {
//		q = q.UpdateDescription(*params.Description)
//	}
//	if params.Icon != nil {
//		if *params.Icon == "" {
//			q = q.UpdateIcon(sql.NullString{Valid: false, String: ""})
//		} else {
//			q = q.UpdateIcon(sql.NullString{Valid: true, String: *params.Icon})
//		}
//	}
//	if params.ParentID != nil {
//		if *params.ParentID == uuid.Nil {
//			q = q.UpdateParent(uuid.NullUUID{Valid: false, UUID: uuid.Nil})
//		} else {
//			q = q.UpdateParent(uuid.NullUUID{Valid: true, UUID: *params.ParentID})
//		}
//	}
//
//	row, err := q.UpdateOne(ctx)
//	if err != nil {
//		return models.PlaceClass{}, err
//	}
//
//	return toModel(row), nil
//}
//
//func (r *Repository) CheckParentCycle(ctx context.Context, classID, newParentID uuid.UUID) (bool, error) {
//	res, err := r.placeClassesQ(ctx).FilterByClassID(newParentID, true, true).FilterByID(classID).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	if res {
//		return true, nil
//	}
//
//	return false, nil
//}
//
//func (r *Repository) CheckPlaceClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error) {
//	res, err := r.placeClassesQ(ctx).FilterByClassID(classID, true, false).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) CheckPlaceExistForClass(ctx context.Context, classID uuid.UUID) (bool, error) {
//	res, err := r.placeQ(ctx).FilterByParentID(classID).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) CheckPlaceClassExists(ctx context.Context, classID uuid.UUID) (bool, error) {
//	res, err := r.placeClassesQ(ctx).FilterByID(classID).Exists(ctx)
//	if err != nil {
//		return false, err
//	}
//	return res, nil
//}
//
//func (r *Repository) DeletePlaceClass(ctx context.Context, classID uuid.UUID) error {
//	return r.placeClassesQ(ctx).FilterByID(classID).Delete(ctx)
//}
//
//func toModel(row pg.PlaceClass) models.PlaceClass {
//	res := models.PlaceClass{
//		ID:          row.ID,
//		Code:        row.Code,
//		Name:        row.Name,
//		Description: row.Description,
//		CreatedAt:   row.CreatedAt,
//		UpdatedAt:   row.UpdatedAt,
//	}
//	if row.ParentID.Valid {
//		res.ParentID = &row.ParentID.UUID
//	}
//	if row.Icon.Valid {
//		res.Icon = &row.Icon.String
//	}
//	return res
//}
