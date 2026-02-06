package controller

import (
	"net/http"
)

// TOOD remade
func (c Controller) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	//initiator, err := contexter.AccountData(r.Context())
	//if err != nil {
	//	c.log.WithError(err).Errorf("failed to get initiator account data")
	//	c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
	//	return
	//}
	//
	//req, err := requests.UpdatePlace(r)
	//if err != nil {
	//	c.log.WithError(err).Errorf("invalid update place request data")
	//	c.responser.RenderErr(w, problems.BadRequest(err)...)
	//	return
	//}
	//
	//res, err := c.core.UpdatePlace(r.Context(), initiator.ID, req.Data.Id, place.UpdateParams{
	//	ClassID:     req.Data.Attributes.ClassId,
	//	Address:     req.Data.Attributes.Address,
	//	Name:        req.Data.Attributes.Name,
	//	Description: req.Data.Attributes.Description,
	//	Icon:        req.Data.Attributes.Icon,
	//	Banner:      req.Data.Attributes.Banner,
	//	Website:     req.Data.Attributes.Website,
	//	Phone:       req.Data.Attributes.Phone,
	//})
	//if err != nil {
	//	c.log.WithError(err).Errorf("failed to update place")
	//	switch {
	//	case errors.Is(err, errx.ErrorPlaceNotFound):
	//		c.responser.RenderErr(w, problems.NotFound("place not found"))
	//	case errors.Is(err, errx.ErrorPlaceClassNotFound):
	//		c.responser.RenderErr(w, problems.NotFound("place class not found"))
	//	case errors.Is(err, errx.ErrorNotEnoughRights):
	//		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update place"))
	//	default:
	//		c.responser.RenderErr(w, problems.InternalError())
	//	}
	//	return
	//}
	//
	//c.responser.Render(w, http.StatusOK, responses.Place(res))
}
