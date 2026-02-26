package responses

import (
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/pkg/resources"
)

func UploadPlaceMediaLink(place models.Place, uploadLink models.UploadPlaceMediaLinks) resources.UploadPlaceMediaLinks {
	return resources.UploadPlaceMediaLinks{
		Data: resources.UploadPlaceMediaLinksData{
			Id:   place.ID,
			Type: "place_upload_links",
			Attributes: resources.UploadPlaceMediaLinksDataAttributes{
				Icon: resources.UploadResourcesLink{
					Key:        uploadLink.Icon.Key,
					UploadUrl:  uploadLink.Icon.UploadURL,
					PreloadUrl: uploadLink.Icon.PreloadUrl,
				},
				Banner: resources.UploadResourcesLink{
					Key:        uploadLink.Banner.Key,
					UploadUrl:  uploadLink.Banner.UploadURL,
					PreloadUrl: uploadLink.Banner.PreloadUrl,
				},
			},
			Relationships: resources.UploadPlaceMediaLinksDataRelationships{
				Place: &resources.UploadPlaceMediaLinksDataRelationshipsPlace{
					Data: resources.UploadPlaceMediaLinksDataRelationshipsPlaceData{
						Id:   place.ID,
						Type: "place",
					},
				},
			},
		},
		Included: []resources.PlaceData{
			Place(place).Data,
		},
	}
}
