package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *OrgModule) authorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.OrgMember, error) {
	member, err := m.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorOrgMemberNotExists) {
		return models.OrgMember{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.OrgMember{}, err
	}

	if !member.Head {
		return models.OrgMember{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only organization head member can manage organization, but member %s is not head", member.ID,
			),
		)
	}

	return member, nil
}

func (m *OrgModule) authorizeOrgMember(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.OrgMember, error) {
	initiator, err := m.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorOrgMemberNotExists) {
		return models.OrgMember{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.OrgMember{}, err
	}

	return initiator, nil
}

func (m *OrgModule) validateOrg(ctx context.Context, organizationID uuid.UUID) (models.Organization, error) {
	org, err := m.org.Get(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	return org, nil
}
