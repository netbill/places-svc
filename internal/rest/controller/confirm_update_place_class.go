package controller

import (
	"errors"
	"net/http"

	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/pclass"
	"github.com/netbill/places-svc/internal/rest/contexter"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) ConfirmUpdatePlaceClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.ConfirmUpdatePlaceClass(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update place class request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	uploadFilesData, err := contexter.UploadContentData(r.Context())
	if err != nil {
		c.log.WithError(err).Error("failed to get upload session id")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get upload session id"))

		return
	}

	res, err := c.core.class.ConfirmUpdateSession(r.Context(), req.Data.Id, pclass.UpdateParams{
		ParentID:    req.Data.Attributes.ParentId,
		Code:        req.Data.Attributes.Code,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
		Media: pclass.UpdateMediaParams{
			UploadSessionID: uploadFilesData.GetUploadSessionID(),
			DeleteIcon:      req.Data.Attributes.DeleteIcon,
		},
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update place class")
		switch {
		case errors.Is(err, errx.ErrorPlaceClassNotFound):
			c.responser.RenderErr(w, problems.NotFound("place class not found"))
		case errors.Is(err, errx.ErrorPlaceClassCodeExists):
			c.responser.RenderErr(w, problems.Conflict("place class code already exists"))
		case errors.Is(err, errx.ErrorPlaceClassParentCycle):
			c.responser.RenderErr(w, problems.Conflict("setting this parent_id would create a cycle in the class hierarchy"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
	}

	c.responser.Render(w, http.StatusOK, responses.PlaceClass(res))
}
