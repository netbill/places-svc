package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/geoukr"
	"github.com/netbill/restkit/pagi"
)

type Service struct {
	repo      repo
	messanger messanger
	territory checkerTerritory
}

func New(repo repo, messanger messanger) Service {
	ter, err := geoukr.New()
	if err != nil {
		panic(fmt.Sprintf("failed to create territory checker: %v", err))
	}

	return Service{
		repo:      repo,
		messanger: messanger,
		territory: ter,
	}
}

type repo interface {
	CreatePlace(ctx context.Context, params CreateParams) (models.Place, error)

	GetPlaceByID(ctx context.Context, id uuid.UUID) (models.Place, error)
	GetPlaces(ctx context.Context, params FilterParams, limit, offset uint) (pagi.Page[[]models.Place], error)

	UpdatePlaceByID(ctx context.Context, id uuid.UUID, params UpdateParams) (models.Place, error)
	UpdatePlaceStatus(ctx context.Context, placeID uuid.UUID, status string) (models.Place, error)
	UpdatePlaceVerified(ctx context.Context, placeID uuid.UUID, verified bool) (models.Place, error)

	DeletePlaceByID(ctx context.Context, id uuid.UUID) error

	CheckPlaceClassExists(ctx context.Context, classID uuid.UUID) (bool, error)

	GetOrgMemberByAccountID(ctx context.Context, organizationID, accountID uuid.UUID) (models.Member, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	PublishCreatePlace(ctx context.Context, place models.Place) error

	PublishUpdatePlace(ctx context.Context, place models.Place) error
	PublishUpdatePlaceStatus(ctx context.Context, place models.Place) error
	PublishUpdatePlaceVerified(ctx context.Context, place models.Place) error
	PublishUpdatePlaceClassID(ctx context.Context, place models.Place) error

	PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error
}

type checkerTerritory interface {
	ContainsLatLng(lat, lng float64) bool
}

func (s Service) chekPermissionForManagePlace(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
) error {
	member, err := s.repo.GetOrgMemberByAccountID(ctx, organizationID, accountID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get initiator member: %w", err))
	}
	if member.IsNil() {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator is not a member of organization %s", organizationID),
		)
	}

	access, err := s.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionManageOrganization,
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}
