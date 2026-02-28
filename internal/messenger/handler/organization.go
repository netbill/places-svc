package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/modules/organization"
)

const operationOrganizationCreated = "organization_created"

func (h *Handler) OrganizationCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationOrganizationCreated).
		With("organization_id", payload.OrganizationID)

	err := h.modules.Org.Create(ctx, organization.CreateParams{
		ID:        payload.OrganizationID,
		Status:    payload.Status,
		Name:      payload.Name,
		IconKey:   payload.IconKey,
		BannerKey: payload.BannerKey,
		CreatedAt: payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorOrganizationDeleted):
		log.Debug("received organization created event for already deleted organization")
		return nil
	case errors.Is(err, errx.ErrorOrganizationAlreadyExists):
		log.Debug("received organization created event for already existing organization")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to create organization")
		return err
	default:
		log.Debug("organization created successfully")
		return nil
	}
}

const operationOrganizationUpdated = "organization_updated"

func (h *Handler) OrganizationUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationOrganizationUpdated).
		With("organization_id", payload.OrganizationID)

	err := h.modules.Org.Update(ctx, payload.OrganizationID, organization.UpdateParams{
		Name:      payload.Name,
		Status:    payload.Status,
		IconKey:   payload.IconKey,
		BannerKey: payload.BannerKey,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorOrganizationDeleted):
		log.Debug("received organization updated event for already deleted organization")
		return nil
	case errors.Is(err, errx.ErrorOrganizationAlreadyExists):
		log.Debug("received organization updated event for already existing organization")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		return err
	default:
		log.Debug("organization updated successfully")
		return nil
	}
}

const operationOrganizationDeleted = "organization_deleted"

func (h *Handler) OrganizationDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.OrganizationDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationOrganizationDeleted).
		With("organization_id", payload.OrganizationID)

	err := h.modules.Org.Delete(ctx, payload.OrganizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationDeleted):
		log.Debug("received organization deleted event for already deleted organization")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to delete organization")
		return err
	default:
		log.Debug("organization deleted successfully")
		return nil
	}
}
