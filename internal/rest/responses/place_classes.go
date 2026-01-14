package responses

import (
	"net/http"

	"github.com/netbill/ape"
	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
)

func classData(model models.PlaceClass) resources.ClassData {
	data := resources.ClassData{
		Id:   model.ID,
		Type: "class",
		Attributes: resources.ClassDataAttributes{
			Code:        model.Code,
			Name:        model.Name,
			Description: model.Description,
			Icon:        model.Icon,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
		},
	}
	if model.ParentID != nil {
		data.Relationships = &resources.ClassDataRelationships{
			Parent: &resources.ClassDataRelationshipsParent{
				Id:   *model.ParentID,
				Type: "class",
			},
		}
	}

	return data
}

func PlaceClass(w http.ResponseWriter, status int, model models.PlaceClass) {
	ape.Render(w, status, resources.Class{
		Data: classData(model),
	})
}

func PlaceClasses(r *http.Request, w http.ResponseWriter, status int, page pagi.Page[[]models.PlaceClass]) {
	data := make([]resources.ClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = classData(mod)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	ape.Render(w, status, resources.ClassesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
		},
	})
}
