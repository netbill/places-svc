package responses

import (
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
)

func placeData(model models.Place) resources.PlaceData {
	data := resources.PlaceData{
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
	}

	if model.OrganizationID != nil {
		data.Relationships.Organization = &resources.PlaceDataRelationshipsOrganization{
			Id:   *model.OrganizationID,
			Type: "organization",
		}
	}

	return data
}
func Place(w http.ResponseWriter, status int, model models.Place) {
	ape.Render(w, status, resources.Place{
		Data: placeData(model),
	})
}

func Places(r *http.Request, w http.ResponseWriter, status int, page pagi.Page[[]models.Place]) {
	data := make([]resources.PlaceData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = placeData(mod)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	response := resources.PlacesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}

	ape.Render(w, status, response)
}
