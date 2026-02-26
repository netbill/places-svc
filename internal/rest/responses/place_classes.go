package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type placeClassResponse struct {
	class    resources.PlaceClassData
	included []resources.PlaceClassData
}

type PlaceClassOption func(*placeClassResponse)

func WithParentClass(model models.PlaceClass) PlaceClassOption {
	return func(r *placeClassResponse) {
		r.included = append(r.included, placeClassData(model))
	}
}

func PlaceClass(model models.PlaceClass, opts ...PlaceClassOption) resources.PlaceClass {
	r := &placeClassResponse{
		class: placeClassData(model),
	}
	for _, opt := range opts {
		opt(r)
	}
	return resources.PlaceClass{
		Data:     r.class,
		Included: r.included,
	}
}

func PlaceClasses(req *http.Request, page pagi.Page[[]models.PlaceClass]) resources.PlaceClassesCollection {
	data := make([]resources.PlaceClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = PlaceClass(mod).Data
	}

	links := pagi.BuildPageLinks(req, page.Page, page.Size, page.Total)

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

func placeClassData(model models.PlaceClass) resources.PlaceClassData {
	res := resources.PlaceClassData{
		Id:   model.ID,
		Type: "place_class",
		Attributes: resources.PlaceClassDataAttributes{
			Name:         model.Name,
			Description:  model.Description,
			IconKey:      model.IconKey,
			Version:      model.Version,
			CreatedAt:    model.CreatedAt,
			UpdatedAt:    model.UpdatedAt,
			DeprecatedAt: model.DeprecatedAt,
		},
	}

	if model.ParentID != nil {
		res.Relationships = &resources.PlaceClassDataRelationships{
			Parent: &resources.PlaceClassDataRelationshipsParent{
				Id:   *model.ParentID,
				Type: "place_class",
			},
		}
	}

	return res
}
