package repository

import (
	"context"

	"github.com/google/uuid"
)

type TombstonesSql interface {
	BuryOrganization(ctx context.Context, orgID uuid.UUID) error
	BuryOrgMember(ctx context.Context, memberID uuid.UUID) error
	BuryPlace(ctx context.Context, placeID uuid.UUID) error

	OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error)
	OrgMemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error)
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)
}
