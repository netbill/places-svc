package place

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	bucket    bucket
	messenger messenger
	territory checkerTerritory
	token     token
}

func New(
	repo repo,
	bucket bucket,
	messenger messenger,
	territory checkerTerritory,
	token token,
) *Module {
	return &Module{
		repo:      repo,
		bucket:    bucket,
		messenger: messenger,
		territory: territory,
		token:     token,
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

	GetOrgMemberByAccountID(ctx context.Context, organizationID, accountID uuid.UUID) (models.OrgMember, error)

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	PublishCreatePlace(ctx context.Context, place models.Place) error

	PublishUpdatePlace(ctx context.Context, place models.Place) error
	PublishUpdatePlaceStatus(ctx context.Context, place models.Place) error
	PublishUpdatePlaceVerified(ctx context.Context, place models.Place) error
	PublishUpdatePlaceClassID(ctx context.Context, place models.Place) error

	PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error
}

type bucket interface {
	GeneratePreloadLinkForPlaceMedia(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) (models.PlaceUploadMediaLinks, error)

	AcceptUpdatePlaceMedia(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) (models.PlaceMedia, error)

	CancelUpdatePlaceIcon(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error

	CancelUpdatePlaceBanner(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error

	DeletePlaceIcon(
		ctx context.Context,
		placeID uuid.UUID,
	) error

	DeletePlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
	) error

	CleanPlaceMediaSession(
		ctx context.Context,
		placeID, sessionID uuid.UUID,
	) error
}

type checkerTerritory interface {
	ContainsLatLng(lat, lng float64) bool
}

type token interface {
	NewUploadPlaceMediaToken(
		OwnerAccountID uuid.UUID,
		ResourceID uuid.UUID,
		UploadSessionID uuid.UUID,
	) (string, error)
}

func (m *Module) chekPermissionForCreatePlace(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
) error {
	member, err := m.repo.GetOrgMemberByAccountID(ctx, organizationID, initiator.GetAccountID())
	if err != nil {
		return err
	}

	if member.Head {
		return nil
	}

	access, err := m.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionPlaceCreate,
	)
	if err != nil {
		return err
	}

	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}

func (m *Module) chekPermissionForDeletePlace(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
) error {
	member, err := m.repo.GetOrgMemberByAccountID(ctx, organizationID, initiator.GetAccountID())
	if err != nil {
		return err
	}

	if member.Head {
		return nil
	}

	access, err := m.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionPlaceDelete,
	)
	if err != nil {
		return err
	}

	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}

func (m *Module) chekPermissionForUpdatePlace(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
) error {
	member, err := m.repo.GetOrgMemberByAccountID(ctx, organizationID, initiator.GetAccountID())
	if err != nil {
		return err
	}

	if member.Head {
		return nil
	}

	access, err := m.repo.CheckMemberHavePermission(
		ctx,
		member.ID,
		models.RolePermissionPlaceUpdate,
	)
	if err != nil {
		return err
	}

	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}
