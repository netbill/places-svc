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

func Places(req *http.Request, page pagi.Page[[]models.Place]) resources.PlacesCollection {
	data := make([]resources.PlaceData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = Place(mod).Data
	}

	links := pagi.BuildPageLinks(req, page.Page, page.Size, page.Total)

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
