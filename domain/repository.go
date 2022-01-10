package domain

import (
	"time"
)

type LiveRepository interface {
	FindByPeriod(start *time.Time, end *time.Time) ([]*Live, error)
	Create(live *Live) error
	Update(live *Live) error
	Delete(live *Live) error
}

type BandRepository interface {
	FindByLiveId(id int) ([]*Band, error)
	Create(band *Band) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
	DeleteByIdAndTurn(id int, turn int) error
}

type BandMemberRepository interface {
	FindByLiveIdAndTurn(id int, turn int) ([]*Player, error)
	Create(bandMember *BandMember) error
	Delete(bandMember *BandMember) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
}

type PlayerRepository interface {
	Create(player *Player) error
	Delete(player *Player) error
	FindByPart(part *Part) ([]*Player, error)
}
