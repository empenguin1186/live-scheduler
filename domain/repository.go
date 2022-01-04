package domain

import (
	"time"
)

type LiveRepository interface {
	FindByDate(time time.Time) Live
}

type BandRepository interface {
	FindByLiveId(id int) []Band
}

type BandMemberRepository interface {
	FindByLiveIdAndTurn(id int, turn int) []Member
}
