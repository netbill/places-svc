package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
	"github.com/paulmach/orb"
)

type placeRepo interface {
	Create(ctx context.Context, params CreatePlaceParams) (models.Place, error)

	Get(ctx context.Context, placeID uuid.UUID) (models.Place, error)
	GetListByIDs(ctx context.Context, placeIDs []uuid.UUID) ([]models.Place, error)
	GetList(ctx context.Context, params FilterPlaceParams, limit, offset uint) (pagi.Page[[]models.Place], error)

	Update(ctx context.Context, placeID uuid.UUID, params UpdatePlaceParams) (models.Place, error)
	UpdateStatus(ctx context.Context, placeID uuid.UUID, status string) (models.Place, error)
	UpdateVerified(ctx context.Context, placeID uuid.UUID, verified bool) (models.Place, error)
	UpdateClass(ctx context.Context, placeID uuid.UUID, classID uuid.UUID) (models.Place, error)

	Delete(ctx context.Context, id uuid.UUID) error
}

type placePClassRepo interface {
	Get(ctx context.Context, classID uuid.UUID) (models.PlaceClass, error)
	Exists(ctx context.Context, classID uuid.UUID) (bool, error)
}

type placeTombstone interface {
	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}

type orgAuth interface {
	authorizeOrgHead(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.OrgMember, error)

	authorizeOrgMember(
		ctx context.Context,
		accountID uuid.UUID,
		organizationID uuid.UUID,
	) (models.OrgMember, error)

	validateOrg(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)
}

type placeMessenger interface {
	PublishCreatePlace(ctx context.Context, place models.Place) error
	PublishUpdatePlace(ctx context.Context, place models.Place) error
	PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error
}

type placeMedia interface {
	CreatePlaceIconUploadMediaLinks(
		ctx context.Context,
		placeID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdatePlaceIcon(
		ctx context.Context,
		orgID uuid.UUID,
		tempKey string,
	) (newKey string, err error)

	DeletePlaceIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadPlaceIcon(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error

	CreatePlaceBannerUploadMediaLinks(
		ctx context.Context,
		placeID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdatePlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) (string, error)

	DeletePlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) error
	DeleteUploadPlaceBanner(
		ctx context.Context,
		placeID uuid.UUID,
		key string,
	) error
}

type checkerTerritory interface {
	ContainsLatLng(lat, lng float64) bool
}

type CreatePlaceParams struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	ClassID        uuid.UUID `json:"class_id"`
	Point          orb.Point `json:"point"`
	Address        string    `json:"address"`
	Name           string    `json:"name"`

	Description *string `json:"description"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`
}

type PlaceModule struct {
	repo      placeRepo
	class     placePClassRepo
	tombstone placeTombstone
	auth      orgAuth
	tx        transaction
	media     placeMedia
	messenger placeMessenger
	territory checkerTerritory
}

type PlaceModuleDeps struct {
	Repo      placeRepo
	Class     placePClassRepo
	Tombstone placeTombstone
	Auth      orgAuth
	Tx        transaction
	Media     placeMedia
	Messenger placeMessenger
	Territory checkerTerritory
}

func NewModule(deps PlaceModuleDeps) *PlaceModule {
	return &PlaceModule{
		repo:      deps.Repo,
		class:     deps.Class,
		tombstone: deps.Tombstone,
		auth:      deps.Auth,
		tx:        deps.Tx,
		media:     deps.Media,
		messenger: deps.Messenger,
		territory: deps.Territory,
	}
}

func (m *PlaceModule) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreatePlaceParams,
) (place models.Place, err error) {
	_, err = m.auth.authorizeOrgHead(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = m.auth.validateOrg(ctx, params.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	if !m.territory.ContainsLatLng(params.Point[1], params.Point[0]) {
		return models.Place{}, errx.ErrorPlaceOutOfTerritory.Raise(
			fmt.Errorf("place point %v is out of allowed territory", params.Point),
		)
	}

	class, err := m.class.Get(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, err
	}
	if class.DeprecatedAt != nil {
		return models.Place{}, errx.ErrorPlaceClassIsDeprecated.Raise(
			fmt.Errorf("place class %s is deprecated", params.ClassID),
		)
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		place, err = m.repo.Create(ctx, params)
		if err != nil {
			return err
		}

		return m.messenger.PublishCreatePlace(ctx, place)
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

func (m *PlaceModule) Get(ctx context.Context, placeID uuid.UUID) (models.Place, error) {
	place, err := m.repo.Get(ctx, placeID)
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

type FilterPlaceParams struct {
	Class          *FilterClassForPlaceParams
	Near           *FilterPlaceNearParams
	OrganizationID *uuid.UUID
	Status         []string
	OrgStatus      []string
	Verified       *bool

	BestMatch *string

	Website *string
	Phone   *string
}

type FilterClassForPlaceParams struct {
	ClassID  []uuid.UUID
	Parents  bool
	Children bool
}

type FilterPlaceNearParams struct {
	Point   orb.Point
	RadiusM uint
}

func (m *PlaceModule) GetList(
	ctx context.Context,
	params FilterPlaceParams,
	limit, offset uint,
) (pagi.Page[[]models.Place], error) {
	res, err := m.repo.GetList(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Place]{}, err
	}

	return res, nil
}

func (m *PlaceModule) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Place, error) {
	res, err := m.repo.GetListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *PlaceModule) Delete(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) error {
	place, err := m.repo.Get(ctx, placeID)
	if err != nil {
		return err
	}

	_, err = m.auth.authorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	_, err = m.auth.validateOrg(ctx, place.OrganizationID)
	if err != nil {
		return err
	}

	buried, err := m.tombstone.PlaceIsBuried(ctx, placeID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", placeID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstone.BuryPlace(ctx, placeID); err != nil {
			return err
		}

		if err = m.repo.Delete(ctx, placeID); err != nil {
			return err
		}

		return m.messenger.PublishDeletePlace(ctx, placeID)
	})
}
