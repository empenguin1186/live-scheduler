package domain

import (
	"time"
)

type LiveRepository interface {
	FindById(id int) (*Live, error)
	FindByPeriod(start *time.Time, end *time.Time) ([]*Live, error)
	Create(live *Live) error
	Update(live *Live) error
	Delete(id int) error
}

type BandRepository interface {
	FindByLiveId(id int) ([]*Band, error)
	Create(band *Band) error
	Update(id int, turn int, band *Band) error
	Delete(id int, turn int) error
}

type BandMemberRepository interface {
	FindByLiveIdAndTurn(id int, turn int) ([]*Player, error)
	Create(bandMember *BandMember) error
	Delete(bandMember *BandMember) error
	Update(bandMember *BandMember, id int, turn int) error
}

type PlayerRepository interface {
	Create(player *Player) error
	Delete(player *Player) error
	FindByPart(part *Part) ([]*Player, error)
}
