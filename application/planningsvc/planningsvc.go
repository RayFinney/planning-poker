package planningsvc

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"planning-poker/domain/planning"
	"planning-poker/infra"
)

type PlanningService struct {
	planningRepository planning.Repository
	logger             *zap.Logger
}

func NewPlanningService(planningRepository planning.Repository) *PlanningService {
	return &PlanningService{
		planningRepository: planningRepository,
		logger:             infra.GetLogger(),
	}
}

// Create creates a new planning
func (svc *PlanningService) Create(p *planning.Planning) error {
	svc.logger.Debug("Creating new planning", zap.String("owner", p.Owner.Name))
	if p.Id == "" {
		p.Id = uuid.NewString()
	}
	p.Votes = make(map[string]int)
	p.HiddenVotes = make(map[string]int)
	err := svc.planningRepository.Create(*p)
	if err != nil {
		svc.logger.Error("Error creating planning", zap.Error(err))
		return err
	}
	svc.logger.Debug("Planning created successfully", zap.String("id", p.Id))
	p.Owner.IsOwner = true
	_, err = svc.Join(p.Id, &p.Owner)
	if err != nil {
		svc.logger.Error("Error joining planning", zap.String("planningId", p.Id), zap.Error(err))
		return err
	}
	return nil
}

// GetById retrieves a planning by its ID
func (svc *PlanningService) GetById(id string, playerId string) (planning.Planning, error) {
	svc.logger.Debug("Retrieving planning by ID", zap.String("id", id))
	p, err := svc.planningRepository.GetById(id)
	if err != nil {
		svc.logger.Error("Error retrieving planning", zap.String("id", id), zap.Error(err))
		return p, err
	}
	myVote, ok := p.HiddenVotes[playerId]
	if !ok {
		myVote = -1
	}
	p.MyVote = myVote
	svc.logger.Debug("Planning retrieved successfully", zap.String("id", p.Id))
	return p, nil
}

// Join allows a player to join a planning
func (svc *PlanningService) Join(planningId string, player *planning.Player) (planning.Planning, error) {
	svc.logger.Debug("Player joining planning", zap.String("planningId", planningId), zap.String("playerName", player.Name))
	if player.Name == "" {
		// TODO: generate a random player name
	}
	player.Id = uuid.NewString()
	p, err := svc.planningRepository.Join(planningId, *player)
	if err != nil {
		svc.logger.Error("Error joining planning", zap.String("planningId", planningId), zap.String("playerName", player.Name), zap.Error(err))
		return planning.Planning{}, err
	}
	svc.logger.Debug("Player joined successfully", zap.String("planningId", planningId), zap.String("playerName", player.Id))
	return p, nil
}

// Leave allows a player to leave a planning
func (svc *PlanningService) Leave(planningId string, playerId string) (planning.Planning, error) {
	svc.logger.Debug("Player leaving planning", zap.String("planningId", planningId), zap.String("playerId", playerId))
	p, err := svc.planningRepository.Leave(planningId, playerId)
	if err != nil {
		svc.logger.Error("Error leaving planning", zap.String("planningId", planningId), zap.String("playerId", playerId), zap.Error(err))
		return planning.Planning{}, err
	}
	svc.logger.Debug("Player left successfully", zap.String("planningId", planningId), zap.String("playerId", playerId))
	return p, nil
}

// Vote allows a player to vote on a planning
func (svc *PlanningService) Vote(planningId string, playerId string, value int) error {
	svc.logger.Debug("Player voting on planning", zap.String("planningId", planningId), zap.String("playerId", playerId), zap.Int("value", value))
	plan, err := svc.planningRepository.GetById(planningId)
	if err != nil {
		svc.logger.Error("Error retrieving planning for voting", zap.String("planningId", planningId), zap.Error(err))
		return err
	}
	if plan.Revealed {
		return nil
	}
	err = svc.planningRepository.Vote(planningId, playerId, value)
	if err != nil {
		svc.logger.Error("Error recording vote", zap.String("planningId", planningId), zap.String("playerId", playerId), zap.Int("value", value), zap.Error(err))
		return err
	}
	svc.logger.Debug("Vote recorded successfully", zap.String("planningId", planningId), zap.String("playerId", playerId), zap.Int("value", value))
	return nil
}

// RevealVotes reveals the votes for a planning
func (svc *PlanningService) RevealVotes(planningId string) (planning.Planning, error) {
	svc.logger.Debug("Revealing votes for planning", zap.String("planningId", planningId))
	p, err := svc.planningRepository.RevealVotes(planningId)
	if err != nil {
		svc.logger.Error("Error revealing votes", zap.String("planningId", planningId), zap.Error(err))
		return p, err
	}
	svc.logger.Debug("Votes revealed successfully", zap.String("planningId", p.Id))
	return p, nil
}

// ResetVotes resets the votes for a planning
func (svc *PlanningService) ResetVotes(planningId string) error {
	svc.logger.Debug("Resetting votes for planning", zap.String("planningId", planningId))
	err := svc.planningRepository.ResetVotes(planningId)
	if err != nil {
		svc.logger.Error("Error resetting votes", zap.String("planningId", planningId), zap.Error(err))
		return err
	}
	svc.logger.Debug("Votes reset successfully", zap.String("planningId", planningId))
	return nil
}

// Close closes a planning
func (svc *PlanningService) Close(planningId string) error {
	svc.logger.Debug("Closing planning", zap.String("planningId", planningId))
	svc.planningRepository.Close(planningId)
	svc.logger.Debug("Planning closed successfully", zap.String("planningId", planningId))
	return nil
}
