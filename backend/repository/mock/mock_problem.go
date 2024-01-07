// Code generated by MockGen. DO NOT EDIT.
// Source: problem.go
//
// Generated by this command:
//
//	mockgen -source=problem.go -destination=./mock/mock_problem.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	repository "fjnkt98/atcodersearch/repository"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockProblemRepository is a mock of ProblemRepository interface.
type MockProblemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProblemRepositoryMockRecorder
}

// MockProblemRepositoryMockRecorder is the mock recorder for MockProblemRepository.
type MockProblemRepositoryMockRecorder struct {
	mock *MockProblemRepository
}

// NewMockProblemRepository creates a new mock instance.
func NewMockProblemRepository(ctrl *gomock.Controller) *MockProblemRepository {
	mock := &MockProblemRepository{ctrl: ctrl}
	mock.recorder = &MockProblemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProblemRepository) EXPECT() *MockProblemRepositoryMockRecorder {
	return m.recorder
}

// FetchIDs mocks base method.
func (m *MockProblemRepository) FetchIDs(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchIDs", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchIDs indicates an expected call of FetchIDs.
func (mr *MockProblemRepositoryMockRecorder) FetchIDs(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchIDs", reflect.TypeOf((*MockProblemRepository)(nil).FetchIDs), ctx)
}

// FetchIDsByContestID mocks base method.
func (m *MockProblemRepository) FetchIDsByContestID(ctx context.Context, contest_id []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchIDsByContestID", ctx, contest_id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchIDsByContestID indicates an expected call of FetchIDsByContestID.
func (mr *MockProblemRepositoryMockRecorder) FetchIDsByContestID(ctx, contest_id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchIDsByContestID", reflect.TypeOf((*MockProblemRepository)(nil).FetchIDsByContestID), ctx, contest_id)
}

// Save mocks base method.
func (m *MockProblemRepository) Save(ctx context.Context, problems []repository.Problem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, problems)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockProblemRepositoryMockRecorder) Save(ctx, problems any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockProblemRepository)(nil).Save), ctx, problems)
}