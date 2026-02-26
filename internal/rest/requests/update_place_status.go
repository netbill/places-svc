package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit"
)

func UpdatePlaceStatus(r *http.Request) (req resources.UpdatePlaceStatus, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("update_place")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "place_id") != req.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "query place_id and body data/id do not match")
	}

	return req, errs.Filter()
}
