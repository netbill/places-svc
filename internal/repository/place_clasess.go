package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/classes"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

func (s Service) CreatePlaceClass(ctx context.Context, params classes.CreateParams) (models.PlaceClass, error) {
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

func (s Service) UpdatePlaceClass(ctx context.Context, classID uuid.UUID, params classes.UpdateParams) (models.PlaceClass, error) {
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
	res, err := s.placeClassesQ(ctx).FilterByParentIDTree(newParentID, 0).FilterByID(classID).Exists(ctx)
	if err != nil {
		return false, err
	}
	if res {
		return true, nil
	}

	res, err = s.placeClassesQ(ctx).FilterByParentIDTree(classID, 0).FilterByID(newParentID).Exists(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s Service) CheckClassHasChildren(ctx context.Context, classID uuid.UUID) (bool, error) {
	res, err := s.placeClassesQ(ctx).FilterByParentIDTree(classID, 0).Exists(ctx)
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
