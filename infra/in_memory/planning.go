package in_memory

import (
	"errors"
	"planning-poker/domain/planning"
	"sync"
)

type PlanningRepository struct {
	activeSessions map[string]planning.Planning
	sessionLock    sync.Mutex
}

func NewPlanningRepository() *PlanningRepository {
	return &PlanningRepository{
		activeSessions: make(map[string]planning.Planning),
	}
}

func (p *PlanningRepository) Create(planning planning.Planning) error {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	if _, ok := p.activeSessions[planning.Id]; ok {
		return errors.New("planning with this id already exists")
	}
	p.activeSessions[planning.Id] = planning
	return nil
}

func (p *PlanningRepository) GetById(id string) (planning.Planning, error) {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[id]
	if !ok {
		return planning.Planning{}, errors.New("planning with this id does not exist")
	}
	return plan, nil
}

func (p *PlanningRepository) Join(planningId string, player planning.Player) (planning.Planning, error) {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[planningId]
	if !ok {
		return planning.Planning{}, errors.New("planning with this id does not exist")
	}
	plan.Players = append(plan.Players, player)
	if player.IsOwner {
		plan.Owner = player
	}
	p.activeSessions[planningId] = plan
	return plan, nil
}

// Leave allows a player to leave a planning session
func (p *PlanningRepository) Leave(planningId string, playerId string) (planning.Planning, error) {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[planningId]
	if !ok {
		return planning.Planning{}, errors.New("planning with this id does not exist")
	}
	var updatedPlayers []planning.Player
	for _, player := range plan.Players {
		if player.Id != playerId {
			updatedPlayers = append(updatedPlayers, player)
		}
	}
	plan.Players = updatedPlayers
	if len(plan.Players) == 0 {
		delete(p.activeSessions, planningId)
		return planning.Planning{}, nil
	}
	delete(plan.Votes, playerId)
	delete(plan.HiddenVotes, playerId)
	if plan.Owner.Id == playerId {
		plan.Owner = plan.Players[0]
	}
	p.activeSessions[planningId] = plan
	return plan, nil
}

func (p *PlanningRepository) Vote(planningId string, playerId string, value int) error {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[planningId]
	if !ok {
		return errors.New("planning with this id does not exist")
	}
	plan.HiddenVotes[playerId] = value
	plan.Votes[playerId] = 0
	p.activeSessions[planningId] = plan
	return nil
}

func (p *PlanningRepository) RevealVotes(planningId string) (planning.Planning, error) {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[planningId]
	if !ok {
		return planning.Planning{}, errors.New("planning with this id does not exist")
	}
	plan.Votes = plan.HiddenVotes
	plan.Revealed = true
	p.activeSessions[planningId] = plan
	return plan, nil
}

func (p *PlanningRepository) ResetVotes(planningId string) error {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	plan, ok := p.activeSessions[planningId]
	if !ok {
		return errors.New("planning with this id does not exist")
	}
	plan.Votes = make(map[string]int)
	plan.HiddenVotes = make(map[string]int)
	plan.Revealed = false
	p.activeSessions[planningId] = plan
	return nil
}

func (p *PlanningRepository) Close(planningId string) {
	p.sessionLock.Lock()
	defer p.sessionLock.Unlock()
	delete(p.activeSessions, planningId)
}
