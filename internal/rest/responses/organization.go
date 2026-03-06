package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/places-svc/pkg/resources"
)

func organizationData(r *http.Request, organization models.Organization) resources.OrganizationData {
	res := resources.OrganizationData{
		Id:   organization.ID,
		Type: "organization",
		Attributes: resources.OrganizationDataAttributes{
			Status:    organization.Status,
			Name:      organization.Name,
			Version:   organization.Version,
			CreatedAt: organization.CreatedAt,
			UpdatedAt: organization.UpdatedAt,
		},
	}

	if organization.IconKey != nil {
		url := scope.ResolverURL(r, *organization.IconKey)
		res.Attributes.IconUrl = &url
	}
	if organization.BannerKey != nil {
		url := scope.ResolverURL(r, *organization.BannerKey)
		res.Attributes.BannerUrl = &url
	}

	return res
}
