package mocks

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, uid)
	// first value passed to "Return"
	var r0 *model.User
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockUserService) Register(ctx context.Context, u *model.User) error {
	// ret := m.Called(ctx, u)

	// var r0 error
	// if ret.Get(0) != nil {
	// 	r0 = ret.Get(0).(error)
	// }
	// return r0
	panic("not implemented")
}

func (s *MockUserService) Login(ctx context.Context, email string, password string) (*model.User, error) {
	panic("Not implemented")
}
