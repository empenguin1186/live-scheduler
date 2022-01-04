package infra

import (
	"github.com/go-gorp/gorp"
	"live-scheduler/domain"
	"log"
	"time"
)

const LAYOUT = "2006-01-02"

type LiveRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewLiveRepositoryImpl(dbmap *gorp.DbMap) LiveRepositoryImpl {
	return LiveRepositoryImpl{dbmap: dbmap}
}

func (i LiveRepositoryImpl) FindByDate(time time.Time) domain.Live {
	var lives []LiveRecord
	_, err := i.dbmap.Select(&lives, "SELECT * FROM Live WHERE date = ?", time.Format(LAYOUT))
	checkErr(err, "SELECT * FROM Live QUERY failed.")
	live := lives[0]
	return live.ToModel()
}

func (i LiveRepositoryImpl) Create(live domain.Live) error {
	record := LiveRecordForInsert{
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	err := i.dbmap.Insert(&record)
	checkErr(err, "INSERT INTO Live QUERY failed.")
	return err
}

func (i LiveRepositoryImpl) UpdateDateById(id int, time time.Time) error {
	panic("implement me")
}

func (i LiveRepositoryImpl) DeleteById(id int) {
	panic("implement me")
}

type BandRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewBandRepositoryImpl(dbmap *gorp.DbMap) BandRepositoryImpl {
	return BandRepositoryImpl{dbmap: dbmap}
}

func (i BandRepositoryImpl) FindByLiveId(id int) []domain.Band {
	var bandRecords []BandRecord
	_, err := i.dbmap.Select(&bandRecords, "SELECT * FROM Band WHERE live_id = ?", id)
	checkErr(err, "SELECT * FROM Band QUERY failed.")

	var bands []domain.Band
	for _, bandRecord := range bandRecords {
		bands = append(bands, bandRecord.ToModel())
	}
	return bands
}

func (i BandRepositoryImpl) Create(band domain.Band) error {
	panic("implement me")
}

func (i BandRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	panic("implement me")
}

func (i BandRepositoryImpl) DeleteByIdAndTurn(id int, turn int) error {
	panic("implement me")
}

type BandMemberRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewBandMemberRepositoryImpl(dbmap *gorp.DbMap) BandMemberRepositoryImpl {
	return BandMemberRepositoryImpl{dbmap: dbmap}
}

func (i BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) []domain.Player {
	var memberRecords []BandMemberRecord
	_, err := i.dbmap.Select(&memberRecords, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", id, turn)
	checkErr(err, "SELECT * FROM BandMember QUERY failed.")

	var members []domain.Player
	for _, memberRecord := range memberRecords {
		members = append(members, memberRecord.ToModel())
	}
	return members
}

func (i BandMemberRepositoryImpl) Create(id int, turn int, name string, part domain.Part) error {
	panic("implement me")
}

func (i BandMemberRepositoryImpl) DeleteOne(id int, turn int, name string, part domain.Part) error {
	panic("implement me")
}

func (i BandMemberRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	panic("implement me")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
