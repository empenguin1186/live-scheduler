package domain

import (
	"time"
)

type LiveRepository interface {
	FindByDate(time time.Time) Live
	Create(live Live) error
	UpdateDateById(id int, time time.Time) error
	DeleteById(id int)
}

type BandRepository interface {
	FindByLiveId(id int) []Band
	Create(band Band) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
	DeleteByIdAndTurn(id int, turn int) error
}

type BandMemberRepository interface {
	FindByLiveIdAndTurn(id int, turn int) []Player
	Create(id int, turn int, name string, part Part) error
	DeleteOne(id int, turn int, name string, part Part) error
	UpdateTurn(id int, prevTurn int, afterTurn int) error
}

type PlayerRepository interface {
	Create(name string, part Part) error
	Delete(name string, part Part) error
	FindByPart(part Part) []Player
}
