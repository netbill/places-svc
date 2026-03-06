package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/models"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type placeClassResponse struct {
	data     resources.PlaceClassData
	included []resources.PlaceClassData
}

type PlaceClassOption func(*placeClassResponse)

func WithParentClass(r *http.Request, model models.PlaceClass) PlaceClassOption {
	return func(res *placeClassResponse) {
		res.included = append(res.included, placeClassData(r, model))
	}
}

func PlaceClass(r *http.Request, model models.PlaceClass, opts ...PlaceClassOption) resources.PlaceClass {
	res := &placeClassResponse{
		data: placeClassData(r, model),
	}
	for _, opt := range opts {
		opt(res)
	}

	return resources.PlaceClass{
		Data:     res.data,
		Included: res.included,
	}
}

func placeClassData(r *http.Request, model models.PlaceClass) resources.PlaceClassData {
	res := resources.PlaceClassData{
		Id:   model.ID,
		Type: "place_class",
		Attributes: resources.PlaceClassDataAttributes{
			Name:         model.Name,
			Description:  model.Description,
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

	if model.IconKey != nil {
		url := scope.ResolverURL(r, *model.IconKey)
		res.Attributes.IconUrl = &url
	}

	return res
}

type placeClassCollectionResponse struct {
	data     []resources.PlaceClassData
	included []resources.PlaceClassData
}

type PlaceClassCollectionOption func(*placeClassCollectionResponse)

func WithCollectionParentClass(r *http.Request, models []models.PlaceClass) PlaceClassCollectionOption {
	return func(res *placeClassCollectionResponse) {
		for _, model := range models {
			res.included = append(res.included, placeClassData(r, model))
		}
	}
}

func PlaceClasses(r *http.Request, page pagi.Page[[]models.PlaceClass], opts ...PlaceClassCollectionOption) resources.PlaceClassesCollection {
	data := make([]resources.PlaceClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = placeClassData(r, mod)
	}

	res := &placeClassCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(res)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	seen := make(map[string]struct{})
	included := make([]resources.PlaceClassData, 0, len(res.included))
	for _, item := range res.included {
		key := item.Id.String()
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			included = append(included, item)
		}
	}

	return resources.PlaceClassesCollection{
		Data:     res.data,
		Included: included,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
		},
	}
}

func UploadPlaceClassMediaLink(
	r *http.Request,
	class models.PlaceClass,
	uploadLink models.UploadPlaceClassMediaLinks,
) resources.UploadPlaceClassMediaLinks {
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
			placeClassData(r, class),
		},
	}
}
