package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type placeResponse struct {
	data     resources.PlaceData
	included []resources.PlaceIncludedInner
}

type PlaceOption func(*placeResponse)

func WithClass(r *http.Request, model models.PlaceClass) PlaceOption {
	return func(res *placeResponse) {
		inner := placeClassData(r, model)
		res.included = append(res.included, resources.PlaceIncludedInner{
			PlaceClassData: &inner,
		})
	}
}

func WithOrganization(r *http.Request, model models.Organization) PlaceOption {
	return func(res *placeResponse) {
		inner := organizationData(r, model)
		res.included = append(res.included, resources.PlaceIncludedInner{
			OrganizationData: &inner,
		})
	}
}

func Place(r *http.Request, model models.Place, opts ...PlaceOption) resources.Place {
	res := &placeResponse{
		data: placeData(r, model),
	}
	for _, opt := range opts {
		opt(res)
	}

	return resources.Place{
		Data:     res.data,
		Included: deduplicatePlaceIncluded(res.included),
	}
}

func placeData(r *http.Request, place models.Place) resources.PlaceData {
	res := resources.PlaceData{
		Id:   place.ID,
		Type: "place",
		Attributes: resources.PlaceDataAttributes{
			Status:   place.Status,
			Verified: place.Verified,
			Point: resources.Point{
				Latitude:  place.Point[1],
				Longitude: place.Point[0],
			},
			Address:     place.Address,
			Name:        place.Name,
			Description: place.Description,
			Website:     place.Website,
			Phone:       place.Phone,
			Version:     place.Version,
			CreatedAt:   place.CreatedAt,
			UpdatedAt:   place.UpdatedAt,
		},
		Relationships: resources.PlaceDataRelationships{
			PlaceClass: resources.PlaceDataRelationshipsPlaceClass{
				Id:   place.ClassID,
				Type: "place_class",
			},
			Organization: resources.PlaceDataRelationshipsOrganization{
				Id:   place.OrganizationID,
				Type: "organization",
			},
		},
	}
	if place.IconKey != nil {
		url := scope.ResolverURL(r, *place.IconKey)
		res.Attributes.IconUrl = &url
	}
	if place.BannerKey != nil {
		url := scope.ResolverURL(r, *place.BannerKey)
		res.Attributes.BannerUrl = &url
	}

	return res
}

type placeCollectionResponse struct {
	data     []resources.PlaceData
	included []resources.PlaceIncludedInner
}

type PlaceCollectionOption func(*placeCollectionResponse)

func WithCollectionClass(r *http.Request, models []models.PlaceClass) PlaceCollectionOption {
	return func(res *placeCollectionResponse) {
		for _, model := range models {
			inner := placeClassData(r, model)
			res.included = append(res.included, resources.PlaceIncludedInner{
				PlaceClassData: &inner,
			})
		}
	}
}

func WithCollectionOrganization(r *http.Request, models []models.Organization) PlaceCollectionOption {
	return func(res *placeCollectionResponse) {
		for _, model := range models {
			inner := organizationData(r, model)
			res.included = append(res.included, resources.PlaceIncludedInner{
				OrganizationData: &inner,
			})
		}
	}
}

func Places(r *http.Request, page pagi.Page[[]models.Place], opts ...PlaceCollectionOption) resources.PlacesCollection {
	data := make([]resources.PlaceData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = placeData(r, mod)
	}

	res := &placeCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(res)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.PlacesCollection{
		Data:     res.data,
		Included: deduplicatePlaceIncluded(res.included),
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func deduplicatePlaceIncluded(items []resources.PlaceIncludedInner) []resources.PlaceIncludedInner {
	seen := make(map[string]struct{})
	result := make([]resources.PlaceIncludedInner, 0, len(items))

	for _, item := range items {
		var key string
		switch {
		case item.PlaceClassData != nil:
			key = "place_class:" + item.PlaceClassData.Id.String()
		case item.OrganizationData != nil:
			key = "organization:" + item.OrganizationData.Id.String()
		default:
			continue
		}

		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func UploadPlaceMediaLink(
	r *http.Request,
	place models.Place,
	uploadLink models.UploadPlaceMediaLinks,
) resources.UploadPlaceMediaLinks {
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
			placeData(r, place),
		},
	}
}
