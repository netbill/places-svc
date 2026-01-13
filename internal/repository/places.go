package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/place"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

func (s Service) CreatePlace(ctx context.Context, params place.CreateParams) (models.Place, error) {
	res, err := s.placeQ(ctx).Insert(ctx, pgdb.PlacesInsertInput{
		OrganizationID: params.OrganizationID,
		ClassID:        params.ClassID,
		Status:         models.PlaceStatusInactive,
		Verified:       false,
		Point:          params.Point,
		Address:        params.Address,
		Name:           params.Name,
		Description:    params.Description,
		Icon:           params.Icon,
		Banner:         params.Banner,
		Website:        params.Website,
		Phone:          params.Phone,
	})
	if err != nil {
		return models.Place{}, err
	}

	return Place(res), nil
}

func (s Service) GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error) {
	row, err := s.placeQ(ctx).FilterByID(id).Get(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return Place(row), nil
}

func (s Service) UpdatePlaceByID(ctx context.Context, id uuid.UUID, params place.UpdatePlaceParams) (models.Place, error) {
	upd := s.placeQ(ctx).FilterByID(id)

	if params.ClassID != nil {
		upd = upd.UpdateClassID(*params.ClassID)
	}
	if params.Address != nil {
		upd = upd.UpdateAddress(*params.Address)
	}
	if params.Name != nil {
		upd = upd.UpdateName(*params.Name)
	}
	if params.Description != nil {
		upd = upd.UpdateDescription(params.Description)
	}
	if params.Icon != nil {
		upd = upd.UpdateIcon(params.Icon)
	}
	if params.Banner != nil {
		upd = upd.UpdateBanner(params.Banner)
	}
	if params.Website != nil {
		upd = upd.UpdateWebsite(params.Website)
	}
	if params.Phone != nil {
		upd = upd.UpdatePhone(params.Phone)
	}

	row, err := upd.UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return Place(row), nil
}

func (s Service) UpdatePlaceStatus(ctx context.Context, placeID uuid.UUID, status string) (models.Place, error) {
	res, err := s.placeQ(ctx).FilterByID(placeID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return Place(res), nil
}

func (s Service) UpdatePlaceVerified(ctx context.Context, placeID uuid.UUID, verified bool) (models.Place, error) {
	res, err := s.placeQ(ctx).FilterByID(placeID).UpdateVerified(verified).UpdateOne(ctx)
	if err != nil {
		return models.Place{}, err
	}

	return Place(res), nil
}

func (s Service) UpdatePlaceStatusForOrg(ctx context.Context, organizationID uuid.UUID, status string) error {
	if _, err := s.placeQ(ctx).FilterByOrganizationID(&organizationID).UpdateStatus(status).UpdateOne(ctx); err != nil {
		return err
	}

	return nil
}

func (s Service) DeletePlaceByID(ctx context.Context, id uuid.UUID) error {
	return s.placeQ(ctx).FilterByID(id).Delete(ctx)
}

func Place(row pgdb.Place) models.Place {
	res := models.Place{
		ID:        row.ID,
		ClassID:   row.ClassID,
		Status:    row.Status,
		Point:     row.Point,
		Verified:  row.Verified,
		Address:   row.Address,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.OrganizationID.Valid {
		orgID := row.OrganizationID.UUID
		res.OrganizationID = &orgID
	}
	if row.Description.Valid {
		res.Description = &row.Description.String
	}
	if row.Icon.Valid {
		res.Icon = &row.Icon.String
	}
	if row.Banner.Valid {
		res.Banner = &row.Banner.String
	}
	if row.Website.Valid {
		res.Website = &row.Website.String
	}
	if row.Phone.Valid {
		res.Phone = &row.Phone.String
	}

	return res
}
