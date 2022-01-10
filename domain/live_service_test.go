package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetByPeriod(t *testing.T) {
	// given
	live := Live{Id: 1, Name: "name", Location: "location", Date: now, PerformanceFee: 5500, EquipmentCost: 2000}

	testCase := []struct {
		testName      string
		expectedLives []*Live
		expectedError error
	}{
		{
			testName:      "正常系",
			expectedLives: []*Live{&live},
			expectedError: nil,
		},
		{
			testName:      "異常系_データアクセス時にエラー発生",
			expectedLives: []*Live{},
			expectedError: fmt.Errorf("dummy message"),
		},
	}

	for _, tc := range testCase {
		liveRepository := new(LiveRepositoryMock)
		liveRepository.On("FindByPeriod", &now, &now).Return(tc.expectedLives, tc.expectedError).Once()
		liveService := NewLiveServiceImpl(liveRepository)

		// when
		actual, err := liveService.GetByPeriod(&now, &now)

		// then
		if strings.Contains(tc.testName, "正常系") {
			assert.Equal(t, tc.expectedLives, actual, fmt.Sprintf("テスト名: %s", tc.testName))
			assert.Nil(t, err, fmt.Sprintf("テスト名: %s", tc.testName))
		} else {
			assert.Equal(t, tc.expectedError, err, fmt.Sprintf("テスト名: %s", tc.testName))
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
		liveService := NewLiveServiceImpl(liveRepository)

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
		liveService := NewLiveServiceImpl(liveRepository)

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
		liveRepository.On("Delete", live.Id).Return(tc.expectedError).Once()
		liveService := NewLiveServiceImpl(liveRepository)

		// when
		actual := liveService.Delete(live.Id)

		// then
		assert.Equal(t, tc.expectedError, actual, fmt.Sprintf("テスト名: %s", tc.testName))
	}
}
