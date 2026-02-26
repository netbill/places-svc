package responses

import (
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/pkg/resources"
)

func UploadPlaceClassMediaLink(class models.PlaceClass, uploadLink models.UploadPlaceClassMediaLinks) resources.UploadPlaceClassMediaLinks {
	return resources.UploadPlaceClassMediaLinks{
		Data: resources.UploadPlaceClassMediaLinksData{
			Id:   class.ID,
			Type: "place_class_upload_links",
			Attributes: resources.UploadPlaceClassMediaLinksDataAttributes{
				Icon: resources.UploadResourcesLink{
					Key:        uploadLink.Icon.Key,
					UploadUrl:  uploadLink.Icon.UploadURL,
					PreloadUrl: uploadLink.Icon.PreloadUrl,
				},
			},
			Relationships: resources.UploadPlaceClassMediaLinksDataRelationships{
				PlaceClass: &resources.UploadPlaceClassMediaLinksDataRelationshipsPlaceClass{
					Data: resources.UploadPlaceClassMediaLinksDataRelationshipsPlaceClassData{
						Id:   class.ID,
						Type: "place_class",
					},
				},
			},
		},
		Included: []resources.PlaceClassData{
			PlaceClass(class).Data,
		},
	}
}
