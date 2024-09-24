package mocks

import (
	"context"

	model "github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProgramRepository struct {
	mock.Mock
}

func (m *MockProgramRepository) GetLimitOffset(ctx context.Context, limit int, offset int) ([]model.MProgram, error) {
	ret := m.Called(ctx, limit, offset)

	var r0 []model.MProgram
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]model.MProgram)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockProgramRepository) FindByID(ctx context.Context, programID uuid.UUID) (*model.MProgram, error) {
	ret := m.Called(ctx, programID)

	var r0 *model.MProgram
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.MProgram)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
