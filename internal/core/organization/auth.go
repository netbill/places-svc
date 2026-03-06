package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/errx"
	"github.com/netbill/places-svc/internal/models"
)

func (s *Service) AuthorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.OrgMember, error) {
	member, err := s.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorOrgMemberNotExists) {
		return models.OrgMember{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a memberRepo of orgRepo %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.OrgMember{}, err
	}

	if !member.Head {
		return models.OrgMember{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only orgRepo head memberRepo can manage orgRepo, but memberRepo %s is not head", member.ID,
			),
		)
	}

	return member, nil
}

func (s *Service) AuthorizeOrgMember(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.OrgMember, error) {
	initiator, err := s.member.GetForAccountAndOrg(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorOrgMemberNotExists) {
		return models.OrgMember{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a memberRepo of orgRepo %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.OrgMember{}, err
	}

	return initiator, nil
}

func (s *Service) ValidateOrg(ctx context.Context, organizationID uuid.UUID) (models.Organization, error) {
	org, err := s.Get(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("orgRepo with id %s is suspended", organizationID),
		)
	}

	return org, nil
}
