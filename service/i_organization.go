package service

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type IOrganizationService interface {
	Get(ctx context.Context, uid uuid.UUID) (*model.Organization, error)
	List(ctx context.Context) ([]model.Organization, error)
	HardDelete(ctx context.Context, uid uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, uReq *model.Organization) error
	Create(ctx context.Context, u model.Organization, userID uuid.UUID) (*model.Organization, error)
}
