package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/places-svc/pkg/resources"
	"github.com/netbill/restkit"
)

func CreatePlaceClass(r *http.Request) (req resources.CreatePlaceClass, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("place_class")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
