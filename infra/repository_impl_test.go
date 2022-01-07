package infra

import (
	"github.com/go-gorp/gorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"live-scheduler/domain"
	"testing"
	"time"
)

type DaoMock struct {
	mock.Mock
	Dao
}

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
	now := time.Now()
	var lives []Live
	dao := new(DaoMock)
	dao.On("Select", &lives, "SELECT * FROM Live WHERE date = ?", []interface{}{now.Format("2006-01-02")}).
		Return(nil, nil).
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
