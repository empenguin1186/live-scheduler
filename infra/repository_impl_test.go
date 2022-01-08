package infra

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"live-scheduler/domain"
	"regexp"
	"testing"
	"time"
)

var now = time.Now()

func TestSample(t *testing.T) {
	// given
	expected := domain.Live{
		Id:             1,
		Name:           "name",
		Location:       "location",
		Date:           now,
		PerformanceFee: 5500,
		EquipmentCost:  2000,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Live WHERE date = ?")).
		WithArgs(now.Format("2006-01-02")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "date", "performance_fee", "equipment_cost"}).
			AddRow(expected.Id, expected.Name, expected.Location, expected.Date, expected.PerformanceFee, expected.EquipmentCost))
	repository := NewLiveRepositoryImpl(db)

	// when
	actual, err := repository.FindByDate(&now)

	// then
	assert.Equal(t, &expected, actual)
	assert.Nil(t, err)
}
