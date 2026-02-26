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

type placeClassCollectionResponse struct {
	data     []resources.PlaceClassData
	included []resources.PlaceClassData
}

type PlaceClassCollectionOption func(*placeClassCollectionResponse)

func WithCollectionParentClass(models []models.PlaceClass) PlaceClassCollectionOption {
	return func(r *placeClassCollectionResponse) {
		for _, model := range models {
			r.included = append(r.included, placeClassData(model))
		}
	}
}

func PlaceClasses(req *http.Request, page pagi.Page[[]models.PlaceClass], opts ...PlaceClassCollectionOption) resources.PlaceClassesCollection {
	data := make([]resources.PlaceClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = placeClassData(mod)
	}

	links := pagi.BuildPageLinks(req, page.Page, page.Size, page.Total)

	resp := &placeClassCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resources.PlaceClassesCollection{
		Data:     resp.data,
		Included: deduplicatePlaceClassIncluded(resp.included),
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
		},
	}
}

func deduplicatePlaceClassIncluded(items []resources.PlaceClassData) []resources.PlaceClassData {
	seen := make(map[string]struct{})
	result := make([]resources.PlaceClassData, 0, len(items))

	for _, item := range items {
		key := item.Id.String()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
