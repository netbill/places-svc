package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/internal/repository/pgdb"
)

func (s Service) UpsertProfile(ctx context.Context, profile models.Profile) error {
	_, err := s.profilesQ(ctx).Upsert(ctx, pgdb.Profile{
		AccountID: profile.AccountID,
		Username:  profile.Username,
		Official:  profile.Official,
		Pseudonym: profile.Pseudonym,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) error {
	_, err := s.profilesQ(ctx).FilterByAccountID(accountID).UpdateUsername(username).UpdateOne(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteProfile(ctx context.Context, accountID uuid.UUID) error {
	return s.profilesQ(ctx).FilterByAccountID(accountID).Delete(ctx)
}
