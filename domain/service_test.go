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

func (m *LiveRepositoryMock) Delete(live *Live) error {
	args := m.Called(live)
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
		// テスト名
		testName string
		// LiveRepository のモックを返す関数
		liveRepository func() *LiveRepositoryMock
		// BandRepository のモックを返す関数
		bandRepository func() *BandRepositoryMock
		// BandMemberRepository のモックを返す関数
		bandMemberRepository func() *BandMemberRepositoryMock
		// 戻り値の期待値(LiveModel)
		expectedLive LiveModel
		// 戻り値の期待値(error)
		expectedError error
	}{
		{
			testName: "正常系",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindByPeriod", &now, &now).Return([]*Live{&live}, nil).Once()
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
			testName: "異常系_Liveレコード取得処理でエラー発生",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindByPeriod", &now, &now).Return([]*Live{}, expectedError).Once()
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
			testName: "異常系_Liveレコードが空",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindByPeriod", &now, &now).Return([]*Live{}, nil).Once()
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
			expectedError: fmt.Errorf("live not found"),
		},
		{
			testName: "異常系_Bandレコード取得処理でエラー発生",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindByPeriod", &now, &now).Return([]*Live{&live}, nil).Once()
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
			testName: "異常系_BandMemberレコード取得処理でエラー発生",
			liveRepository: func() *LiveRepositoryMock {
				liveRepository := new(LiveRepositoryMock)
				liveRepository.On("FindByPeriod", &now, &now).Return([]*Live{&live}, nil).Once()
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
		liveService := NewLiveServiceImpl(tc.liveRepository(), tc.bandRepository(), tc.bandMemberRepository())

		// when
		actual, err := liveService.GetByDate(&now)

		// then
		assertion := assert.New(t)
		if strings.Contains(tc.testName, "正常系") {
			assertion.Equal(&tc.expectedLive, actual, fmt.Sprintf("テスト名: %s", tc.testName))
			assertion.Nil(err)
		} else {
			assertion.Equal(tc.expectedError, err, fmt.Sprintf("テスト名: %s", tc.testName))
		}
	}
}

var testCase = []struct {
	testName      string
	expectedError error
}{
	{
		testName:      "正常系",
		expectedError: nil,
	},
	{
		testName:      "異常系",
		expectedError: fmt.Errorf("dummy message"),
	},
}

func TestRegister(t *testing.T) {
	// given
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}

	for _, tc := range testCase {
		liveRepository := new(LiveRepositoryMock)
		liveRepository.On("Create", &live).Return(tc.expectedError).Once()
		bandRepository := new(BandRepositoryMock)
		bandMemberRepository := new(BandMemberRepositoryMock)
		liveService := NewLiveServiceImpl(liveRepository, bandRepository, bandMemberRepository)

		// when
		actual := liveService.Register(&live)

		// then
		assert.Equal(t, tc.expectedError, actual, fmt.Sprintf("テスト名: %s", tc.testName))
	}
}

func TestUpdate(t *testing.T) {
	// given
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}

	for _, tc := range testCase {
		liveRepository := new(LiveRepositoryMock)
		liveRepository.On("Update", &live).Return(tc.expectedError).Once()
		bandRepository := new(BandRepositoryMock)
		bandMemberRepository := new(BandMemberRepositoryMock)
		liveService := NewLiveServiceImpl(liveRepository, bandRepository, bandMemberRepository)

		// when
		actual := liveService.Update(&live)

		// then
		assert.Equal(t, tc.expectedError, actual, fmt.Sprintf("テスト名: %s", tc.testName))
	}
}

func TestDelete(t *testing.T) {
	// given
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}

	for _, tc := range testCase {
		liveRepository := new(LiveRepositoryMock)
		liveRepository.On("Delete", &live).Return(tc.expectedError).Once()
		bandRepository := new(BandRepositoryMock)
		bandMemberRepository := new(BandMemberRepositoryMock)
		liveService := NewLiveServiceImpl(liveRepository, bandRepository, bandMemberRepository)

		// when
		actual := liveService.Delete(&live)

		// then
		assert.Equal(t, tc.expectedError, actual, fmt.Sprintf("テスト名: %s", tc.testName))
	}
}
