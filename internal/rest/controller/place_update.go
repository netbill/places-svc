package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/rest/requests"
	"github.com/netbill/places-svc/internal/rest/responses"
	"github.com/netbill/places-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationUpdatePlace = "update_place"

func (c *PlaceController) Update(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlace)

	req, err := requests.UpdatePlace(r)
	if err != nil {
		log.WithError(err).Warn("invalid update place request")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("place_id", req.Data.Id)

	res, err := c.place.Update(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		core.UpdatePlaceParams{
			ClassID:     req.Data.Attributes.ClassId,
			Name:        req.Data.Attributes.Name,
			Address:     req.Data.Attributes.Address,
			Description: req.Data.Attributes.Description,
			Website:     req.Data.Attributes.Website,
			Phone:       req.Data.Attributes.Phone,
			IconKey:     req.Data.Attributes.IconKey,
			BannerKey:   req.Data.Attributes.BannerKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to update place")
		render.ResponseError(w, problems.Forbidden("not enough rights to update place"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Conflict("organization is suspended"))
	case errors.Is(err, errx.ErrorPlaceClassIsDeprecated):
		log.WithError(err).Warn("place class is deprecated")
		render.ResponseError(w, problems.Conflict("place class is deprecated"))
	case errors.Is(err, errx.ErrorPlaceIconIsInvalid):
		log.WithError(err).Warn("upload icon is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"icon": err,
		})...)
	case errors.Is(err, errx.ErrorPlaceBannerIsInvalid):
		log.WithError(err).Warn("upload banner is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"banner": err,
		})...)
	case err != nil:
		log.WithError(err).Error("failed to update place")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place updated")
		render.Response(w, http.StatusOK, responses.Place(r, res))
	}
}

const operationUpdatePlaceStatus = "update_place_status"

func (c *PlaceController) Activate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceStatus)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	res, err := c.place.Activate(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to update place status")
		render.ResponseError(w, problems.Forbidden("not enough rights to update place status"))
	case err != nil:
		log.WithError(err).Error("failed to update place status")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place activated updated")
		render.Response(w, http.StatusOK, responses.Place(r, res))
	}
}

func (c *PlaceController) Deactivate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceStatus)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	res, err := c.place.Deactivate(r.Context(), scope.AccountActor(r), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.WithError(err).Warn("not enough rights to update place status")
		render.ResponseError(w, problems.Forbidden("not enough rights to update place status"))
	case err != nil:
		log.WithError(err).Error("failed to update place status")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("place deactivated updated")
		render.Response(w, http.StatusOK, responses.Place(r, res))
	}
}

const operationUpdatePlaceVerify = "update_place_verify"

func (c *PlaceController) Verify(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceVerify)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	res, err := c.place.Verify(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case err != nil:
		log.WithError(err).Error("failed to update place verify")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Place(r, res))
	}
}

const operationUpdatePlaceUnverify = "update_place_unverify"

func (c *PlaceController) Unverify(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdatePlaceUnverify)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		log.WithError(err).Info("invalid place id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/place_id": fmt.Errorf("invalid place id: %s", chi.URLParam(r, "place_id")),
		})...)
		return
	}

	log = log.WithField("place_id", placeID)

	res, err := c.place.Unverify(r.Context(), placeID)
	switch {
	case errors.Is(err, errx.ErrorPlaceNotExists):
		log.WithError(err).Warn("place not found")
		render.ResponseError(w, problems.NotFound("place not found"))
	case err != nil:
		log.WithError(err).Error("failed to update place unverify")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Place(r, res))
	}
}
