package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type placeResponse struct {
	place    resources.PlaceData
	included []resources.PlaceIncludedInner
}

type PlaceOption func(*placeResponse)

func WithClass(model models.PlaceClass) PlaceOption {
	return func(r *placeResponse) {
		inner := placeClassData(model)
		r.included = append(r.included, resources.PlaceIncludedInner{
			PlaceClassData: &inner,
		})
	}
}

func WithOrganization(model models.Organization) PlaceOption {
	return func(r *placeResponse) {
		r.included = append(r.included, resources.PlaceIncludedInner{
			OrganizationData: organizationData(model),
		})
	}
}

func Place(model models.Place, opts ...PlaceOption) resources.Place {
	r := &placeResponse{
		place: placeData(model),
	}
	for _, opt := range opts {
		opt(r)
	}
	return resources.Place{
		Data:     r.place,
		Included: r.included,
	}
}

func placeData(model models.Place) resources.PlaceData {
	return resources.PlaceData{
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
			IconKey:     model.IconKey,
			BannerKey:   model.BannerKey,
			Website:     model.Website,
			Phone:       model.Phone,
			Version:     model.Version,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
		},
		Relationships: resources.PlaceDataRelationships{
			PlaceClass: resources.PlaceDataRelationshipsPlaceClass{
				Id:   model.ClassID,
				Type: "place_class",
			},
			Organization: resources.PlaceDataRelationshipsOrganization{
				Id:   model.OrganizationID,
				Type: "organization",
			},
		},
	}
}

type placeCollectionResponse struct {
	data     []resources.PlaceData
	included []resources.PlaceIncludedInner
}

type PlaceCollectionOption func(*placeCollectionResponse)

func WithCollectionClass(models []models.PlaceClass) PlaceCollectionOption {
	return func(r *placeCollectionResponse) {
		for _, model := range models {
			inner := placeClassData(model)
			r.included = append(r.included, resources.PlaceIncludedInner{
				PlaceClassData: &inner,
			})
		}
	}
}

func WithCollectionOrganization(models []models.Organization) PlaceCollectionOption {
	return func(r *placeCollectionResponse) {
		for _, model := range models {
			r.included = append(r.included, resources.PlaceIncludedInner{
				OrganizationData: organizationData(model),
			})
		}
	}
}

func Places(req *http.Request, page pagi.Page[[]models.Place], opts ...PlaceCollectionOption) resources.PlacesCollection {
	data := make([]resources.PlaceData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = placeData(mod)
	}

	links := pagi.BuildPageLinks(req, page.Page, page.Size, page.Total)

	resp := &placeCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resources.PlacesCollection{
		Data:     resp.data,
		Included: deduplicatePlaceIncluded(resp.included),
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
