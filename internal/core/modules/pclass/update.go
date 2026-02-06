package pclass

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	initiator models.Initiator,
	placeClassID uuid.UUID,
) (models.PlaceClass, models.UpdatePlaceClassMedia, error) {
	org, err := m.Get(ctx, placeClassID)
	if err != nil {
		return models.PlaceClass{}, models.UpdatePlaceClassMedia{}, err
	}

	uploadSessionID := uuid.New()
	links, err := m.bucket.GeneratePreloadLinkForPlaceClassMedia(ctx, org.ID, uploadSessionID)
	if err != nil {
		return models.PlaceClass{}, models.UpdatePlaceClassMedia{}, err
	}

	uploadToken, err := m.token.NewUploadPlaceClassMediaToken(
		initiator.GetAccountID(),
		placeClassID,
		uploadSessionID,
	)
	if err != nil {
		return models.PlaceClass{}, models.UpdatePlaceClassMedia{}, err
	}

	return org, models.UpdatePlaceClassMedia{
		Links: models.PlaceClassUploadMediaLinks{
			IconUploadURL: links.IconUploadURL,
			IconGetURL:    links.IconGetURL,
		},
		UploadSessionID: uploadSessionID,
		UploadToken:     uploadToken,
	}, nil
}

type UpdateParams struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`

	Media UpdateMediaParams `json:"media"`
}

type UpdateMediaParams struct {
	UploadSessionID uuid.UUID `json:"upload_session_id"`

	icon       *string
	DeleteIcon *bool `json:"delete_icon"`
}

func (p UpdateParams) GetUpdatedIcon() *string {
	if p.Media.DeleteIcon != nil && *p.Media.DeleteIcon {
		return nil
	}
	return p.Media.icon
}

func (m *Module) ConfirmUpdateSession(
	ctx context.Context,
	classID uuid.UUID,
	params UpdateParams,
) (class models.PlaceClass, err error) {
	class, err = m.repo.GetPlaceClass(ctx, classID)
	if err != nil {
		return models.PlaceClass{}, err
	}

	if params.ParentID != nil {
		exist, err := m.repo.CheckParentCycle(ctx, class.ID, *class.ParentID)
		if err != nil {
			return models.PlaceClass{}, err
		}
		if exist {
			return models.PlaceClass{}, errx.ErrorPlaceClassParentCycle.Raise(
				fmt.Errorf("setting parent %s for class %s would create a cycle", *class.ParentID, class.ID),
			)
		}
	}

	codeIsUsed, err := m.repo.PlaceClassExistsByCode(ctx, params.Code)
	if err != nil {
		return models.PlaceClass{}, err
	}
	if codeIsUsed {
		return models.PlaceClass{}, errx.ErrorPlaceClassCodeExists.Raise(
			fmt.Errorf("place class code already in use"),
		)
	}

	params.Media.icon = class.Icon

	if params.Media.DeleteIcon != nil || *params.Media.DeleteIcon {
		if err = m.bucket.DeletePlaceClassIcon(ctx, classID); err != nil {
			return models.PlaceClass{}, err
		}
		params.Media.icon = nil
	}

	if params.Media.DeleteIcon != nil || *params.Media.DeleteIcon == false {
		links, err := m.bucket.AcceptUpdatePlaceClassMedia(
			ctx,
			classID,
			params.Media.UploadSessionID,
		)
		if err != nil {
			return models.PlaceClass{}, err
		}

		params.Media.icon = links.Icon
	}

	if err = m.bucket.CleanPlaceClassMediaSession(ctx, classID, params.Media.UploadSessionID); err != nil {
		return models.PlaceClass{}, err
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		class, err = m.repo.UpdatePlaceClass(ctx, classID, params)
		if err != nil {
			return err
		}

		if err = m.messanger.PublishPlaceClassUpdated(ctx, class); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.PlaceClass{}, err
	}

	return class, nil
}

func (m *Module) DeleteUpdateIconInSession(
	ctx context.Context,
	placeID, uploadSessionID uuid.UUID,
) error {
	return m.bucket.CancelUpdatePlaceClassIcon(
		ctx,
		placeID,
		uploadSessionID,
	)
}

func (m *Module) CancelUpdate(
	ctx context.Context,
	placeID, uploadSessionID uuid.UUID,
) error {
	return m.bucket.CleanPlaceClassMediaSession(
		ctx,
		placeID,
		uploadSessionID,
	)
}
