// Code generated by MockGen. DO NOT EDIT.
// Source: contest.go
//
// Generated by this command:
//
//	mockgen -source=contest.go -destination=./mock/mock_contest.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	repository "fjnkt98/atcodersearch/repository"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockContestRepository is a mock of ContestRepository interface.
type MockContestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockContestRepositoryMockRecorder
}

// MockContestRepositoryMockRecorder is the mock recorder for MockContestRepository.
type MockContestRepositoryMockRecorder struct {
	mock *MockContestRepository
}

// NewMockContestRepository creates a new mock instance.
func NewMockContestRepository(ctrl *gomock.Controller) *MockContestRepository {
	mock := &MockContestRepository{ctrl: ctrl}
	mock.recorder = &MockContestRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContestRepository) EXPECT() *MockContestRepositoryMockRecorder {
	return m.recorder
}

// FetchCategories mocks base method.
func (m *MockContestRepository) FetchCategories(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchCategories", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchCategories indicates an expected call of FetchCategories.
func (mr *MockContestRepositoryMockRecorder) FetchCategories(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchCategories", reflect.TypeOf((*MockContestRepository)(nil).FetchCategories), ctx)
}

// FetchContestIDs mocks base method.
func (m *MockContestRepository) FetchContestIDs(ctx context.Context, categories []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchContestIDs", ctx, categories)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchContestIDs indicates an expected call of FetchContestIDs.
func (mr *MockContestRepositoryMockRecorder) FetchContestIDs(ctx, categories any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchContestIDs", reflect.TypeOf((*MockContestRepository)(nil).FetchContestIDs), ctx, categories)
}

// Save mocks base method.
func (m *MockContestRepository) Save(ctx context.Context, contests []repository.Contest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, contests)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockContestRepositoryMockRecorder) Save(ctx, contests any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockContestRepository)(nil).Save), ctx, contests)
}
