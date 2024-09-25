package repository

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
)

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type IOrganizationRepository interface {
	Create(ctx context.Context, org model.Organization, userID uuid.UUID) (*model.Organization, error)
	Update(ctx context.Context, id uuid.UUID, org *model.Organization) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Organization, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]model.Organization, error)
	ListWhereCreatorID(ctx context.Context, creatorID uuid.UUID) ([]model.Organization, error)
}
