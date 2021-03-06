package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
	"time"
)

type LiveRepositoryMock struct {
	mock.Mock
	LiveRepository
}

func (m *LiveRepositoryMock) FindById(id int) (*Live, error) {
	args := m.Called(id)
	return args.Get(0).(*Live), args.Error(1)
}

func (m *LiveRepositoryMock) FindByPeriod(start *time.Time, end *time.Time) ([]*Live, error) {
	args := m.Called(start, end)
	return args.Get(0).([]*Live), args.Error(1)
}

func (m *LiveRepositoryMock) Create(live *Live) error {
	args := m.Called(live)
	return args.Error(0)
}

func (m *LiveRepositoryMock) Update(live *Live) error {
	args := m.Called(live)
	return args.Error(0)
}

func (m *LiveRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type BandRepositoryMock struct {
	mock.Mock
	BandRepository
}

func (m *BandRepositoryMock) FindByLiveId(id int) ([]*Band, error) {
	args := m.Called(id)
	return args.Get(0).([]*Band), args.Error(1)
}

type BandMemberRepositoryMock struct {
	mock.Mock
	BandMemberRepository
}

func (m *BandMemberRepositoryMock) FindByLiveIdAndTurn(id int, turn int) ([]*Player, error) {
	args := m.Called(id, turn)
	return args.Get(0).([]*Player), args.Error(1)
}

var now = time.Now()

func TestGetByDate(t *testing.T) {
	// given
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}
	bands := []*Band{&Band{Name: "band1", LiveId: 1, Turn: 1}, &Band{Name: "band2", LiveId: 1, Turn: 2}}
	players1 := []*Player{&Player{Name: "player1", Part: Ba}, &Player{Name: "player2", Part: Dr}}
	players2 := []*Player{&Player{Name: "player3", Part: Gt}, &Player{Name: "player4", Part: Key}}
	expectedLive := LiveModel{
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
	expectedError := fmt.Errorf("dummy message")

	// given
	tests := []struct {
		// ????????????
		testName string
		// LiveRepository ???????????????????????????
		liveRepository func() *LiveRepositoryMock
		// BandRepository ???????????????????????????
		bandRepository func() *BandRepositoryMock
		// BandMemberRepository ???????????????????????????
		bandMemberRepository func() *BandMemberRepositoryMock
		// ?????????????????????(LiveModel)
		expectedLive LiveModel
		// ?????????????????????(error)
		expectedError error
	}{
		{
			testName: "?????????",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindById", live.Id).Return(&live, nil).Once()
				return liveRepository
			},
			bandRepository: func() *BandRepositoryMock {
				bandRepository := new(BandRepositoryMock)
				bandRepository.On("FindByLiveId", 1).Return(bands, nil).Once()
				return bandRepository
			},
			bandMemberRepository: func() *BandMemberRepositoryMock {
				bandMemberRepository := new(BandMemberRepositoryMock)
				bandMemberRepository.
					On("FindByLiveIdAndTurn", 1, 1).Return(players1, nil).Once().
					On("FindByLiveIdAndTurn", 1, 2).Return(players2, nil).Once()
				return bandMemberRepository
			},
			expectedLive:  expectedLive,
			expectedError: nil,
		},
		{
			testName: "?????????_Live??????????????????????????????????????????",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindById", live.Id).Return(&Live{}, expectedError).Once()
				return liveRepository
			},
			bandRepository: func() *BandRepositoryMock {
				bandRepository := new(BandRepositoryMock)
				bandRepository.On("FindByLiveId", mock.Anything).Times(0)
				return bandRepository
			},
			bandMemberRepository: func() *BandMemberRepositoryMock {
				bandMemberRepository := new(BandMemberRepositoryMock)
				bandMemberRepository.On("FindByLiveIdAndTurn", mock.Anything, mock.Anything).Times(0)
				return bandMemberRepository
			},
			expectedLive:  LiveModel{},
			expectedError: expectedError,
		},
		{
			testName: "?????????_Band??????????????????????????????????????????",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindById", live.Id).Return(&live, nil).Once()
				return liveRepository
			},
			bandRepository: func() *BandRepositoryMock {
				bandRepository := new(BandRepositoryMock)
				bandRepository.On("FindByLiveId", 1).Return([]*Band{}, expectedError).Once()
				return bandRepository
			},
			bandMemberRepository: func() *BandMemberRepositoryMock {
				bandMemberRepository := new(BandMemberRepositoryMock)
				bandMemberRepository.On("FindByLiveIdAndTurn", mock.Anything, mock.Anything).Times(0)
				return bandMemberRepository
			},
			expectedLive:  LiveModel{},
			expectedError: expectedError,
		},
		{
			testName: "?????????_BandMember??????????????????????????????????????????",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindById", live.Id).Return(&live, nil).Once()
				return liveRepository
			},
			bandRepository: func() *BandRepositoryMock {
				bandRepository := new(BandRepositoryMock)
				bandRepository.On("FindByLiveId", 1).Return(bands, nil).Once()
				return bandRepository
			},
			bandMemberRepository: func() *BandMemberRepositoryMock {
				bandMemberRepository := new(BandMemberRepositoryMock)
				bandMemberRepository.
					On("FindByLiveIdAndTurn", 1, 1).Return(players1, nil).Once().
					On("FindByLiveIdAndTurn", 1, 2).Return([]*Player{}, expectedError).Once()
				return bandMemberRepository
			},
			expectedLive:  LiveModel{},
			expectedError: expectedError,
		},
	}

	for _, tc := range tests {
		liveDescService := NewLiveDescServiceImpl(tc.liveRepository(), tc.bandRepository(), tc.bandMemberRepository())

		// when
		actual, err := liveDescService.GetById(live.Id)

		// then
		assertion := assert.New(t)
		if strings.Contains(tc.testName, "?????????") {
			assertion.Equal(&tc.expectedLive, actual, fmt.Sprintf("????????????: %s", tc.testName))
			assertion.Nil(err)
		} else {
			assertion.Equal(tc.expectedError, err, fmt.Sprintf("????????????: %s", tc.testName))
		}
	}
}
