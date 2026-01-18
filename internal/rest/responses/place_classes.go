package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
	"github.com/netbill/restkit/pagi"
)

func PlaceClass(model models.PlaceClass) resources.PlaceClassData {
	data := resources.PlaceClassData{
		Id:   model.ID,
		Type: "class",
		Attributes: resources.PlaceClassDataAttributes{
			Code:        model.Code,
			Name:        model.Name,
			Description: model.Description,
			Icon:        model.Icon,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
		},
	}
	if model.ParentID != nil {
		data.Relationships = &resources.PlaceClassDataRelationships{
			Parent: &resources.PlaceClassDataRelationshipsParent{
				Id:   *model.ParentID,
				Type: "class",
			},
		}
	}

	return data
}

func PlaceClasses(r *http.Request, page pagi.Page[[]models.PlaceClass]) resources.PlaceClassesCollection {
	data := make([]resources.PlaceClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = PlaceClass(mod)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.PlaceClassesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
		},
	}
}
