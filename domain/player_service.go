package domain

type PlayerService interface {
	Register(player *Player) error
	Delete(player *Player) error
	GetByPart(part *Part) ([]*Player, error)
}

type PlayerServiceImpl struct {
	playerRepository PlayerRepository
}

func NewPlayerServiceImpl(playerRepository PlayerRepository) *PlayerServiceImpl {
	return &PlayerServiceImpl{playerRepository: playerRepository}
}

func (p *PlayerServiceImpl) Register(player *Player) error {
	return p.playerRepository.Create(player)
}

func (p *PlayerServiceImpl) Delete(player *Player) error {
	return p.playerRepository.Delete(player)
}

func (p *PlayerServiceImpl) GetByPart(part *Part) ([]*Player, error) {
	return p.playerRepository.FindByPart(part)
}
