package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
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

func (s Service) GetPlaces(
	ctx context.Context,
	params place.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Place], error) {
	q := s.placeQ(ctx)

	if params.OrganizationID != nil {
		q = q.FilterByOrganizationID(params.OrganizationID)
	}
	if params.Status != nil {
		q = q.FilterByStatus(*params.Status)
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
		collection = append(collection, Place(row))
	}

	return pagi.Page[[]models.Place]{
		Data:  collection,
		Total: total,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
	}, nil

}

func (s Service) UpdatePlaceByID(ctx context.Context, id uuid.UUID, params place.UpdateParams) (models.Place, error) {
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
		if *params.Description == "" {
			upd = upd.UpdateDescription(sql.NullString{Valid: false, String: ""})
		} else {
			upd = upd.UpdateDescription(sql.NullString{Valid: true, String: *params.Description})
		}
	}
	if params.Icon != nil {
		if *params.Icon == "" {
			upd = upd.UpdateIcon(sql.NullString{Valid: false, String: ""})
		} else {
			upd = upd.UpdateIcon(sql.NullString{Valid: true, String: *params.Icon})
		}
	}
	if params.Banner != nil {
		if *params.Banner == "" {
			upd = upd.UpdateBanner(sql.NullString{Valid: false, String: ""})
		} else {
			upd = upd.UpdateBanner(sql.NullString{Valid: true, String: *params.Banner})
		}
	}
	if params.Website != nil {
		if *params.Website == "" {
			upd = upd.UpdateWebsite(sql.NullString{Valid: false, String: ""})
		} else {
			upd = upd.UpdateWebsite(sql.NullString{Valid: true, String: *params.Website})
		}
	}
	if params.Phone != nil {
		if *params.Phone == "" {
			upd = upd.UpdatePhone(sql.NullString{Valid: false, String: ""})
		} else {
			upd = upd.UpdatePhone(sql.NullString{Valid: true, String: *params.Phone})
		}
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

func (s Service) ReplacePlacesClassID(ctx context.Context, oldClassID, newClassID uuid.UUID) error {
	_, err := s.placeQ(ctx).
		FilterByClassID(false, false, oldClassID).
		UpdateClassID(newClassID).
		UpdateMany(ctx)
	if err != nil {
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
