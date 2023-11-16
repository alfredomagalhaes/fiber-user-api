package repositories

import (
	"github.com/alfredomagalhaes/fiber-user-api/types"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(*types.User) error
	Find(uuid.UUID) (types.User, error)
	ListAll() ([]types.User, error)
	Update() error
	Delete() error
}
