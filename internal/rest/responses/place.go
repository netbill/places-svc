package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
	"github.com/netbill/restkit/pagi"
)

func Place(model models.Place) resources.Place {
	res := resources.Place{
		Data: resources.PlaceData{
			Id:   model.ID,
			Type: "place",
			Attributes: resources.PlaceDataAttributes{
				Status:   model.Status,
				Verified: model.Verified,
				Point: resources.Point{
					Latitude:  model.Point[1],
					Longitude: model.Point[0],
				},

				Address:     model.Address,
				Name:        model.Name,
				Description: model.Description,
				Icon:        model.Icon,
				Banner:      model.Banner,
				Website:     model.Website,
				Phone:       model.Phone,
				CreatedAt:   model.CreatedAt,
				UpdatedAt:   model.UpdatedAt,
			},
			Relationships: resources.PlaceDataRelationships{
				Class: resources.PlaceDataRelationshipsClass{
					Id:   model.ClassID,
					Type: "class",
				},
			},
		},
	}

	if model.OrganizationID != nil {
		res.Data.Relationships.Organization = &resources.PlaceDataRelationshipsOrganization{
			Id:   *model.OrganizationID,
			Type: "organization",
		}
	}

	return res
}

func Places(r *http.Request, page pagi.Page[[]models.Place]) resources.PlacesCollection {
	data := make([]resources.PlaceData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = Place(mod).Data
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.PlacesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}

}

func OpenUpdatePlaceSession(
	place models.Place,
	uploadLinks models.UpdatePlaceMedia,
) resources.OpenUpdatePlaceMediaLinks {
	return resources.OpenUpdatePlaceMediaLinks{
		Data: resources.OpenUpdatePlaceMediaLinksData{
			Id:   uploadLinks.UploadSessionID,
			Type: "update_place_session",
			Attributes: resources.OpenUpdatePlaceMediaLinksDataAttributes{
				UploadToken:     uploadLinks.UploadToken,
				IconGetUrl:      uploadLinks.Links.IconGetURL,
				IconUploadUrl:   uploadLinks.Links.IconUploadURL,
				BannerGetUrl:    uploadLinks.Links.BannerGetURL,
				BannerUploadUrl: uploadLinks.Links.BannerUploadURL,
			},
			Relationships: resources.OpenUpdatePlaceMediaLinksDataRelationships{
				Place: &resources.OpenUpdatePlaceMediaLinksDataRelationshipsPlace{
					Data: resources.OpenUpdatePlaceMediaLinksDataRelationshipsPlaceData{
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
