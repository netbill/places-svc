package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/class"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

func (s Service) CreatePlaceClass(ctx context.Context, params class.CreateParams) (models.PlaceClass, error) {
	row, err := s.placeClassesQ(ctx).Insert(ctx, pgdb.PlaceClassesInsertInput{
		ParentID:    params.ParentID,
		Code:        params.Code,
		Name:        params.Name,
		Description: params.Description,
		Icon:        params.Icon,
	})
	if err != nil {
		return models.PlaceClass{}, err
	}

	return toModel(row), nil
}

func (s Service) GetClass(ctx context.Context, id uuid.UUID) (models.PlaceClass, error) {
	row, err := s.placeClassesQ(ctx).FilterByID(id).Get(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return toModel(row), nil
}

func (s Service) GetClasses(ctx context.Context, params class.FilterParams, limit, offset uint) (pagi.Page[[]models.PlaceClass], error) {
	q := s.placeClassesQ(ctx)

	if params.Name != nil {
		q = q.FilterNameLike(*params.Name)
	}
	if params.Description != nil {
		q = q.FilterDescriptionLike(*params.Description)
	}
	if params.Text != nil {
		q = q.FilterByText(*params.Text)
	}
	if params.Parent != nil {
		q = q.FilterByClassID(params.Parent.ParentID, params.Parent.Children, params.Parent.Parents)
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
		collection[i] = toModel(row)
	}

	return pagi.Page[[]models.PlaceClass]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil
}

func (s Service) UpdatePlaceClass(ctx context.Context, classID uuid.UUID, params class.UpdateParams) (models.PlaceClass, error) {
	q := s.placeClassesQ(ctx).FilterByID(classID)
	if params.Code != nil {
		q = q.UpdateCode(*params.Code)
	}
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Description != nil {
		q = q.UpdateDescription(*params.Description)
	}
	if params.Icon != nil {
		q = q.UpdateIcon(params.Icon)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return toModel(row), nil
}

func (s Service) CheckParentCycle(ctx context.Context, classID, newParentID uuid.UUID) (bool, error) {
	res, err := s.placeClassesQ(ctx).FilterByClassID(newParentID, true, true).FilterByID(classID).Exists(ctx)
	if err != nil {
		return false, err
	}
	if res {
		return true, nil
	}

	return false, nil
}

func (s Service) CheckClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error) {
	res, err := s.placeClassesQ(ctx).FilterByClassID(classID, true, false).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s Service) CheckPlaceExistForClass(ctx context.Context, classID uuid.UUID) (bool, error) {
	res, err := s.placeQ(ctx).FilterByParentID(classID).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s Service) CheckPlaceClassExists(ctx context.Context, classID uuid.UUID) (bool, error) {
	res, err := s.placeClassesQ(ctx).FilterByID(classID).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s Service) DeletePlaceClass(ctx context.Context, classID uuid.UUID) error {
	return s.placeClassesQ(ctx).FilterByID(classID).Delete(ctx)
}

func (s Service) UpdatePlaceClassParent(ctx context.Context, classID uuid.UUID, parentID *uuid.UUID) (models.PlaceClass, error) {
	row, err := s.placeClassesQ(ctx).
		FilterByID(classID).
		UpdateParent(parentID).
		UpdateOne(ctx)
	if err != nil {
		return models.PlaceClass{}, err
	}

	return toModel(row), nil
}

func toModel(row pgdb.PlaceClass) models.PlaceClass {
	res := models.PlaceClass{
		ID:          row.ID,
		Code:        row.Code,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if row.ParentID.Valid {
		res.ParentID = &row.ParentID.UUID
	}
	if row.Icon.Valid {
		res.Icon = &row.Icon.String
	}
	return res
}
