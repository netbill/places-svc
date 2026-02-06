package controller

import (
	"net/http"
)

// TODO remade
func (c Controller) UpdatePlaceClass(w http.ResponseWriter, r *http.Request) {
	//req, err := requests.UpdatePlaceClass(r)
	//if err != nil {
	//	c.log.WithError(err).Errorf("invalid update place class request")
	//	c.responser.RenderErr(w, problems.BadRequest(err)...)
	//	return
	//}
	//
	//res, err := c.core.UpdatePlaceClass(r.Context(), req.Data.Id, pclass.UpdateParams{
	//	ParentID:    req.Data.Attributes.ParentId,
	//	Code:        req.Data.Attributes.Code,
	//	Name:        req.Data.Attributes.Name,
	//	Description: req.Data.Attributes.Description,
	//	Icon:        req.Data.Attributes.Icon,
	//})
	//if err != nil {
	//	c.log.WithError(err).Errorf("failed to update place class")
	//	switch {
	//	case errors.Is(err, errx.ErrorPlaceClassNotFound):
	//		c.responser.RenderErr(w, problems.NotFound("place class not found"))
	//	case errors.Is(err, errx.ErrorPlaceClassCodeExists):
	//		c.responser.RenderErr(w, problems.Conflict("place class code already exists"))
	//	case errors.Is(err, errx.ErrorPlaceClassParentCycle):
	//		c.responser.RenderErr(w, problems.Conflict("setting this parent_id would create a cycle in the class hierarchy"))
	//	default:
	//		c.responser.RenderErr(w, problems.InternalError())
	//	}
	//}
	//
	//c.responser.Render(w, http.StatusOK, responses.PlaceClass(res))
}
