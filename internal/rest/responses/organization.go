package responses

import (
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/pkg/resources"
)

func organizationData(model models.Organization) *resources.OrganizationData {
	return &resources.OrganizationData{
		Id:   model.ID,
		Type: "organization",
		Attributes: resources.OrganizationDataAttributes{
			Status:    model.Status,
			Name:      model.Name,
			IconKey:   model.IconKey,
			BannerKey: model.BannerKey,
			Version:   model.Version,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
	}
}
