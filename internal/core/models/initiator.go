package models

import (
	"github.com/google/uuid"
)

type AccountClaims interface {
	GetAccountID() uuid.UUID
}
