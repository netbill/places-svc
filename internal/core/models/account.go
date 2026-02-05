package models

import (
	"github.com/google/uuid"
)

type Initiator interface {
	GetAccountID() uuid.UUID
}
