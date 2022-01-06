package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type LiveRepositoryMock struct {
	mock.Mock
	LiveRepository
}

func (m *LiveRepositoryMock) FindByDate(time *time.Time) *Live {
	args := m.Called(time)
	return args.Get(0).(*Live)
}

type BandRepositoryMock struct {
	mock.Mock
	BandRepository
}

func (m *BandRepositoryMock) FindByLiveId(id int) []*Band {
	args := m.Called(id)
	return args.Get(0).([]*Band)
}

type BandMemberRepositoryMock struct {
	mock.Mock
	BandMemberRepository
}

func (m *BandMemberRepositoryMock) FindByLiveIdAndTurn(id int, turn int) []*Player {
	args := m.Called(id, turn)
	return args.Get(0).([]*Player)
}

func TestSomething(t *testing.T) {
	// given
	now := time.Now()
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}
	liveRepository := new(LiveRepositoryMock)
	liveRepository.On("FindByDate", &now).Return(&live).Once()

	bands := []*Band{&Band{Name: "band1", LiveId: 1, Turn: 1}, &Band{Name: "band2", LiveId: 1, Turn: 2}}
	bandRepository := new(BandRepositoryMock)
	bandRepository.On("FindByLiveId", 1).Return(bands).Once()

	players1 := []*Player{&Player{Name: "player1", Part: Ba}, &Player{Name: "player2", Part: Dr}}
	players2 := []*Player{&Player{Name: "player3", Part: Gt}, &Player{Name: "player4", Part: Key}}
	bandMemberRepository := new(BandMemberRepositoryMock)
	bandMemberRepository.
		On("FindByLiveIdAndTurn", 1, 1).
		Return(players1).Once().
		On("FindByLiveIdAndTurn", 1, 2).
		Return(players2).Once()

	expected := LiveModel{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
		Band: []*BandModel{
			&BandModel{Name: "band1", LiveId: 1, Turn: 1, Player: players1},
			&BandModel{Name: "band2", LiveId: 1, Turn: 2, Player: players2},
		},
	}

	liveService := NewLiveService(liveRepository, bandRepository, bandMemberRepository)

	// when
	actual := liveService.GetByDate(&now)

	// then
	assertion := assert.New(t)
	assertion.Equal(&expected, actual)
}
