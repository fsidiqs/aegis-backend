package mocks

import (
	"context"

	model "github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/repository"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock type for model.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) WithTrx() repository.IUserRepository {
	// return &MockUserRepository{}
	panic("unimplemented")
}

func (m *MockUserRepository) CommitTrx() error {
	panic("unimplemented")
}

func (m *MockUserRepository) RollbackTrx() {
	panic("unimplemented")
}

// FindByID is mock of UserRepository FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, uid uint) (*model.User, error) {
	ret := m.Called(ctx, uid)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Create is a mock for UserRepository Create
func (m *MockUserRepository) Create(ctx context.Context, u model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// FindByEmail is mock of UserRepository.FindByEmail
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := m.Called(ctx, email)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Update is mock of UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// UpdateImage is mock of UserRepository.UpdateImage, duh
func (m *MockUserRepository) UpdateImage(
	ctx context.Context,
	uid uint,
	imageURL string,
) (*model.User, error) {
	ret := m.Called(ctx, uid, imageURL)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
