package mocks

import (
	"context"

	model "github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenId uuid.UUID) (*model.TokenData, error) {
	panic("not impleneted")
	// ret := m.Called(ctx, u, prevTokenID)
	// // first value passed to "Return"
	// var r0 *service.TokenData
	// if ret.Get(0) != nil {
	// 	// we can just return this if we know we won't be passing function to "Return"
	// 	r0 = ret.Get(0).(*service.TokenData)
	// }

	// var r1 error

	// if ret.Get(1) != nil {
	// 	r1 = ret.Get(1).(error)
	// }
	// _, _ = r0, r1
	// return r0, r1
	// return &service.TokenData{AuthToken: "authToken",
	// 	RefreshToken: "refreshToken"}, nil
}
