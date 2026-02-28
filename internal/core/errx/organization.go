package errx

import "github.com/netbill/ape"

var (
	ErrorOrganizationNotExists     = ape.DeclareError("ORGANIZATION_NOT_EXISTS")
	ErrorOrganizationAlreadyExists = ape.DeclareError("ORGANIZATION_ALREADY_EXISTS")
	ErrorOrganizationDeleted       = ape.DeclareError("ORGANIZATION_DELETED")

	ErrorOrganizationIsSuspended = ape.DeclareError("ORGANIZATION_SUSPENDED")

	ErrorOrgMemberNotExists     = ape.DeclareError("ORG_MEMBER_NOT_EXISTS")
	ErrorOrgMemberAlreadyExists = ape.DeclareError("ORG_MEMBER_ALREADY_EXISTS")
	ErrorOrgMemberDeleted       = ape.DeclareError("ORG_MEMBER_DELETED")
)
