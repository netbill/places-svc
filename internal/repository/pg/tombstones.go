package pg

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/pgdbx"
	"github.com/netbill/places-svc/internal/repository"
)

type tombstones struct {
	db *pgdbx.DB
}

func NewTombstonesQ(db *pgdbx.DB) repository.TombstonesSql {
	return &tombstones{db: db}
}

func (t *tombstones) BuryPlace(ctx context.Context, placeID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
		INSERT INTO tombstones (entity_type, entity_id)
		VALUES ('place', $1)
		ON CONFLICT (entity_type, entity_id) DO NOTHING
	`, placeID)
	if err != nil {
		return fmt.Errorf("burying place: %w", err)
	}

	return nil
}

func (t *tombstones) PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tombstones
			WHERE entity_type = 'place' AND entity_id = $1
		)
	`, placeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking place is buried: %w", err)
	}

	return exists, nil
}

func (t *tombstones) BuryOrganization(ctx context.Context, orgID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
		INSERT INTO tombstones (entity_type, entity_id)
		SELECT 'organization', $1
		UNION ALL
		SELECT 'organization_member', om.id FROM organization_members om WHERE om.organization_id = $1
		ON CONFLICT (entity_type, entity_id) DO NOTHING
	`, orgID)
	if err != nil {
		return fmt.Errorf("burying organization: %w", err)
	}

	return nil
}

func (t *tombstones) OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tombstones
			WHERE entity_type = 'organization' AND entity_id = $1
		)
	`, orgID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking organization is buried: %w", err)
	}

	return exists, nil
}

func (t *tombstones) OrgMemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tombstones
			WHERE entity_type = 'organization_member' AND entity_id = $1
		)
	`, memberID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking org member is buried: %w", err)
	}

	return exists, nil
}

func (t *tombstones) BuryOrgMember(ctx context.Context, memberID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
		INSERT INTO tombstones (entity_type, entity_id)
		VALUES ('organization_member', $1)
		ON CONFLICT (entity_type, entity_id) DO NOTHING
	`, memberID)
	if err != nil {
		return fmt.Errorf("burying org member: %w", err)
	}

	return nil
}
