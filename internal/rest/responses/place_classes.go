package responses

import (
	"net/http"

	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
	"github.com/netbill/restkit/pagi"
)

func PlaceClass(model models.PlaceClass) resources.PlaceClass {
	res := resources.PlaceClass{
		Data: resources.PlaceClassData{
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
		},
	}
	if model.ParentID != nil {
		res.Data.Relationships = &resources.PlaceClassDataRelationships{
			Parent: &resources.PlaceClassDataRelationshipsParent{
				Id:   *model.ParentID,
				Type: "class",
			},
		}
	}

	return res
}

func PlaceClasses(r *http.Request, page pagi.Page[[]models.PlaceClass]) resources.PlaceClassesCollection {
	data := make([]resources.PlaceClassData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = PlaceClass(mod).Data
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

func OpenUpdatePlaceClassSession(
	class models.PlaceClass,
	uploadLinks models.UpdatePlaceClassMedia,
) resources.OpenUpdatePlaceClassMediaLinks {
	return resources.OpenUpdatePlaceClassMediaLinks{
		Data: resources.OpenUpdatePlaceClassMediaLinksData{
			Id:   uploadLinks.UploadSessionID,
			Type: "update_placeClass_session",
			Attributes: resources.OpenUpdatePlaceClassMediaLinksDataAttributes{
				UploadToken:   uploadLinks.UploadToken,
				IconGetUrl:    uploadLinks.Links.IconGetURL,
				IconUploadUrl: uploadLinks.Links.IconUploadURL,
			},
			Relationships: resources.OpenUpdatePlaceClassMediaLinksDataRelationships{
				PlaceClass: &resources.OpenUpdatePlaceClassMediaLinksDataRelationshipsPlaceClass{
					Data: resources.OpenUpdatePlaceClassMediaLinksDataRelationshipsPlaceClassData{
						Id:   class.ID,
						Type: "placeClass",
					},
				},
			},
		},
		Included: []resources.PlaceClassData{
			PlaceClass(class).Data,
		},
	}
}
