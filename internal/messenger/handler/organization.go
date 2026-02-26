package handler

import (
	"context"
	"encoding/json"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

func (h *Handler) OrganizationCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.organization.Create(ctx, models.Organization{
		ID:        payload.OrganizationID,
		Status:    payload.Status,
		Name:      payload.Name,
		IconKey:   payload.IconKey,
		BannerKey: payload.BannerKey,
		Version:   payload.Version,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.CreatedAt,
	})
}

func (h *Handler) OrganizationUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.organization.Update(ctx, payload.OrganizationID, organization.UpdateParams{
		Name:      payload.Name,
		Status:    payload.Status,
		IconKey:   payload.IconKey,
		BannerKey: payload.BannerKey,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
}

func (h *Handler) OrganizationDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.organization.Delete(ctx, payload.OrganizationID)
}
