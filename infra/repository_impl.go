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

type BandMemberRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewBandMemberRepositoryImpl(dbmap *gorp.DbMap) BandMemberRepositoryImpl {
	return BandMemberRepositoryImpl{dbmap: dbmap}
}

func (i BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) []domain.Member {
	var memberRecords []BandMemberRecord
	_, err := i.dbmap.Select(&memberRecords, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", id, turn)
	checkErr(err, "SELECT * FROM BandMember QUERY failed.")

	var members []domain.Member
	for _, memberRecord := range memberRecords {
		members = append(members, memberRecord.ToModel())
	}
	return members
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
