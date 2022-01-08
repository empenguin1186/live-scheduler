package infra

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-gorp/gorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"live-scheduler/domain"
	"regexp"
	"testing"
	"time"
)

type DaoMock struct {
	mock.Mock
	Dao
}

var now = time.Now()

func (m *DaoMock) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	arguments := m.Called(i, query, args)
	return arguments.Get(0).([]interface{}), arguments.Error(1)
}

func (m *DaoMock) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	arguments := m.Called(i, name)
	return arguments.Get(0).(*gorp.TableMap)
}

func TestFindByDate(t *testing.T) {
	// given
	var lives []Live
	dao := new(DaoMock)
	dao.On("Select", &lives, "SELECT * FROM Live WHERE date = ?", []interface{}{now.Format("2006-01-02")}).
		Return(make([]interface{}, 1), nil).
		Once()
	dao.On("AddTableWithName", Live{}, "Live").
		Return(&gorp.TableMap{}).
		Once()
	liveRepository := NewLiveRepositoryImpl(dao)

	// when
	actual := liveRepository.FindByDate(&now)

	// then
	assertion := assert.New(t)
	assertion.Equal(domain.Live{}, actual)
}

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
	repository := NewLiveRepositoryAlpha(db)

	// when
	actual := repository.FindByDate(&now)

	// then
	assert.Equal(t, &expected, actual)
}
