package domain

import (
	"time"
)

type LiveRepository interface {
	FindByDate(time *time.Time) *Live
	Create(live *Live) error
	Update(live *Live) error
	Delete(live *Live) error
}

type BandRepository interface {
	FindByLiveId(id int) []*Band
	Create(band *Band) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
	DeleteByIdAndTurn(id int, turn int) error
}

type BandMemberRepository interface {
	FindByLiveIdAndTurn(id int, turn int) []*Player
	Create(bandMember *BandMember) error
	Delete(bandMember *BandMember) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
}

type PlayerRepository interface {
	Create(player *Player) error
	Delete(player *Player) error
	FindByPart(part *Part) []*Player
}
