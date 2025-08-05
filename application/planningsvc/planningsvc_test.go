package planningsvc

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"planning-poker/domain/planning"
	"testing"
)

type MockPlanningRepository struct {
	mock.Mock
}

func (m *MockPlanningRepository) Create(p planning.Planning) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPlanningRepository) GetById(id string) (planning.Planning, error) {
	args := m.Called(id)
	return args.Get(0).(planning.Planning), args.Error(1)
}

func (m *MockPlanningRepository) Join(planningId string, player planning.Player) error {
	args := m.Called(planningId, player)
	return args.Error(0)
}

func (m *MockPlanningRepository) Vote(planningId string, playerId string, value int) {
	m.Called(planningId, playerId, value)
}

func (m *MockPlanningRepository) RevealVotes(planningId string) (planning.Planning, error) {
	args := m.Called(planningId)
	return args.Get(0).(planning.Planning), args.Error(1)
}

func (m *MockPlanningRepository) ResetVotes(planningId string) error {
	args := m.Called(planningId)
	return args.Error(0)
}

func (m *MockPlanningRepository) Close(planningId string) {
	m.Called(planningId)
}

func TestPlanningService_Create(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	p := &planning.Planning{
		Owner: planning.Player{Name: "test-owner"},
	}

	mockRepo.On("Create", mock.AnythingOfType("planning.Planning")).Return(nil)
	mockRepo.On("Join", mock.AnythingOfType("string"), mock.AnythingOfType("planning.Player")).Return(nil)

	err := service.Create(p)

	assert.NoError(t, err)
	assert.NotEmpty(t, p.Id)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_CreateErr(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	p := &planning.Planning{
		Owner: planning.Player{Name: "test-owner"},
	}

	mockRepo.On("Create", mock.AnythingOfType("planning.Planning")).Return(errors.New("create error"))

	err := service.Create(p)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_GetById(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	expectedPlanning := planning.Planning{Id: "test-id"}
	mockRepo.On("GetById", "test-id").Return(expectedPlanning, nil)

	p, err := service.GetById("test-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedPlanning, p)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_GetByIdErr(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	mockRepo.On("GetById", "test-id").Return(planning.Planning{}, errors.New("not found"))

	_, err := service.GetById("test-id")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_Join(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	player := planning.Player{Name: "test-player"}
	planningId := uuid.NewString()

	mockRepo.On("Join", planningId, mock.AnythingOfType("planning.Player")).Return(nil)

	err := service.Join(planningId, player)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_JoinErr(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	player := planning.Player{Name: "test-player"}
	planningId := uuid.NewString()

	mockRepo.On("Join", planningId, mock.AnythingOfType("planning.Player")).Return(errors.New("join error"))

	err := service.Join(planningId, player)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_Vote(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()
	playerId := uuid.NewString()
	value := 5

	mockRepo.On("Vote", planningId, playerId, value).Return()

	err := service.Vote(planningId, playerId, value)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_RevealVotes(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()
	expectedPlanning := planning.Planning{Id: planningId, Votes: map[string]int{"player1": 5}}

	mockRepo.On("RevealVotes", planningId).Return(expectedPlanning, nil)

	p, err := service.RevealVotes(planningId)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlanning, p)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_RevealVotesErr(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()

	mockRepo.On("RevealVotes", planningId).Return(planning.Planning{}, errors.New("reveal error"))

	_, err := service.RevealVotes(planningId)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_ResetVotes(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()

	mockRepo.On("ResetVotes", planningId).Return(nil)

	err := service.ResetVotes(planningId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_ResetVotesErr(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()

	mockRepo.On("ResetVotes", planningId).Return(errors.New("reset error"))

	err := service.ResetVotes(planningId)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPlanningService_Close(t *testing.T) {
	mockRepo := new(MockPlanningRepository)
	service := NewPlanningService(mockRepo)

	planningId := uuid.NewString()

	mockRepo.On("Close", planningId).Return()

	err := service.Close(planningId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
