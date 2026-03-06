package places

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/restkit/pagi"
	"github.com/paulmach/orb"
)

type orgAuth interface {
	AuthorizeOrgHead(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.OrgMember, error)

	AuthorizeOrgMember(
		ctx context.Context,
		accountID uuid.UUID,
		organizationID uuid.UUID,
	) (models.OrgMember, error)

	ValidateOrg(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)
}

type messenger interface {
	PublishCreatePlace(ctx context.Context, place models.Place) error
	PublishUpdatePlace(ctx context.Context, place models.Place) error
	PublishDeletePlace(ctx context.Context, placeID uuid.UUID) error
}

type checkerTerritory interface {
	ContainsLatLng(lat, lng float64) bool
}

type Service struct {
	place     placeRepo
	class     classRepo
	tombstone tombstoneRepo
	org       orgAuth
	tx        transaction
	media     media
	messenger messenger
	territory checkerTerritory
}

type PlaceModuleDeps struct {
	Repo      placeRepo
	Class     classRepo
	Tombstone tombstoneRepo
	Org       orgAuth
	Tx        transaction
	Media     media
	Messenger messenger
	Territory checkerTerritory
}

func NewModule(deps PlaceModuleDeps) *Service {
	return &Service{
		place:     deps.Repo,
		class:     deps.Class,
		tombstone: deps.Tombstone,
		org:       deps.Org,
		tx:        deps.Tx,
		media:     deps.Media,
		messenger: deps.Messenger,
		territory: deps.Territory,
	}
}

type CreateParams struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	ClassID        uuid.UUID `json:"class_id"`
	Point          orb.Point `json:"point"`
	Address        string    `json:"address"`
	Name           string    `json:"name"`

	Description *string `json:"description"`
	Website     *string `json:"website"`
	Phone       *string `json:"phone"`
}

func (s *Service) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (place models.Place, err error) {
	_, err = s.org.AuthorizeOrgHead(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	_, err = s.org.ValidateOrg(ctx, params.OrganizationID)
	if err != nil {
		return models.Place{}, err
	}

	if !s.territory.ContainsLatLng(params.Point[1], params.Point[0]) {
		return models.Place{}, errx.ErrorPlaceOutOfTerritory.Raise(
			fmt.Errorf("place point %v is out of allowed territory", params.Point),
		)
	}

	class, err := s.class.Get(ctx, params.ClassID)
	if err != nil {
		return models.Place{}, err
	}
	if class.DeprecatedAt != nil {
		return models.Place{}, errx.ErrorPlaceClassIsDeprecated.Raise(
			fmt.Errorf("place class %s is deprecated", params.ClassID),
		)
	}

	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		place, err = s.place.Create(ctx, params)
		if err != nil {
			return err
		}

		return s.messenger.PublishCreatePlace(ctx, place)
	})
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}

func (s *Service) Get(ctx context.Context, placeID uuid.UUID) (models.Place, error) {
	place, err := s.place.Get(ctx, placeID)
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

func (s *Service) GetList(
	ctx context.Context,
	params FilterPlaceParams,
	limit, offset uint,
) (pagi.Page[[]models.Place], error) {
	res, err := s.place.GetList(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Place]{}, err
	}

	return res, nil
}

func (s *Service) GetByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Place, error) {
	res, err := s.place.GetListByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) Delete(
	ctx context.Context,
	actor models.AccountActor,
	placeID uuid.UUID,
) error {
	place, err := s.place.Get(ctx, placeID)
	if err != nil {
		return err
	}

	_, err = s.org.AuthorizeOrgHead(ctx, actor, place.OrganizationID)
	if err != nil {
		return err
	}

	_, err = s.org.ValidateOrg(ctx, place.OrganizationID)
	if err != nil {
		return err
	}

	buried, err := s.tombstone.PlaceIsBuried(ctx, placeID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorPlaceDeleted.Raise(
			fmt.Errorf("place with id %s is already deleted", placeID),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = s.tombstone.BuryPlace(ctx, placeID); err != nil {
			return err
		}

		if err = s.place.Delete(ctx, placeID); err != nil {
			return err
		}

		return s.messenger.PublishDeletePlace(ctx, placeID)
	})
}
