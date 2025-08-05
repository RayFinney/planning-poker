package planning

type Repository interface {
	Create(planning Planning) error
	GetById(id string) (Planning, error)
	Join(planningId string, player Player) (Planning, error)
	Leave(planningId string, playerId string) (Planning, error)
	Vote(planningId string, playerId string, value int) error
	RevealVotes(planningId string) (Planning, error)
	ResetVotes(planningId string) error
	Close(planningId string)
}
