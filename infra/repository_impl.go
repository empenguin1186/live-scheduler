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
	dbmap.AddTableWithName(Live{}, "Live")
	return LiveRepositoryImpl{dbmap: dbmap}
}

func (i LiveRepositoryImpl) FindByDate(time time.Time) domain.Live {
	var lives []Live
	_, err := i.dbmap.Select(&lives, "SELECT * FROM Live WHERE date = ?", time.Format(LAYOUT))
	checkErr(err, "SELECT * FROM Live QUERY failed.")
	live := lives[0]
	return live.ToModel()
}

func (i LiveRepositoryImpl) Create(live domain.Live) error {
	record := Live{
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

func (i LiveRepositoryImpl) Update(live domain.Live) error {
	record := Live{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	_, err := i.dbmap.Update(&record)
	checkErr(err, "UPDATE Live QUERY failed.")
	return err
}

func (i LiveRepositoryImpl) Delete(live domain.Live) error {
	record := Live{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	_, err := i.dbmap.Delete(&record)
	checkErr(err, "DELETE Live QUERY failed.")
	return err
}

type BandRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewBandRepositoryImpl(dbmap *gorp.DbMap) BandRepositoryImpl {
	dbmap.AddTableWithName(Band{}, "Band")
	return BandRepositoryImpl{dbmap: dbmap}
}

func (i BandRepositoryImpl) FindByLiveId(id int) []domain.Band {
	var bandRecords []Band
	_, err := i.dbmap.Select(&bandRecords, "SELECT * FROM Band WHERE live_id = ?", id)
	checkErr(err, "SELECT * FROM Band QUERY failed.")

	var bands []domain.Band
	for _, bandRecord := range bandRecords {
		bands = append(bands, bandRecord.ToModel())
	}
	return bands
}

func (i BandRepositoryImpl) Create(band domain.Band) error {
	record := Band{
		Name:   band.Name,
		LiveId: band.LiveId,
	}
	err := i.dbmap.Insert(&record)
	checkErr(err, "INSERT INTO Band QUERY failed.")
	return err
}

func (i BandRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := i.dbmap.Exec("UPDATE Band SET turn=? WHERE live_id=? AND turn=?", afterTurn, id, prevTurn)
	checkErr(err, "UPDATE Band QUERY failed.")
	return err
}

func (i BandRepositoryImpl) DeleteByIdAndTurn(id int, turn int) error {
	_, err := i.dbmap.Exec("DELETE FROM Band WHERE live_id=? AND turn=?", id, turn)
	checkErr(err, "DELETE Band QUERY failed.")
	return err
}

type BandMemberRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewBandMemberRepositoryImpl(dbmap *gorp.DbMap) BandMemberRepositoryImpl {
	dbmap.AddTableWithName(BandMember{}, "BandMember")
	return BandMemberRepositoryImpl{dbmap: dbmap}
}

func (i BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) []domain.Player {
	var memberRecords []BandMember
	_, err := i.dbmap.Select(&memberRecords, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", id, turn)
	checkErr(err, "SELECT * FROM BandMember QUERY failed.")

	var members []domain.Player
	for _, memberRecord := range memberRecords {
		members = append(members, memberRecord.ToModel())
	}
	return members
}

func (i BandMemberRepositoryImpl) Create(bandMember domain.BandMember) error {
	record := BandMember{
		LiveId:     bandMember.LiveId,
		Turn:       bandMember.Turn,
		MemberName: bandMember.MemberName,
		MemberPart: string(bandMember.MemberPart),
	}
	err := i.dbmap.Insert(&record)
	checkErr(err, "INSERT INTO BandMember QUERY failed.")
	return err
}

func (i BandMemberRepositoryImpl) Delete(bandMember domain.BandMember) error {
	record := BandMember{
		LiveId:     bandMember.LiveId,
		Turn:       bandMember.Turn,
		MemberName: bandMember.MemberName,
		MemberPart: string(bandMember.MemberPart),
	}
	_, err := i.dbmap.Delete(&record)
	checkErr(err, "INSERT INTO BandMember QUERY failed.")
	return err
}

func (i BandMemberRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := i.dbmap.Exec("UPDATE BandMember SET turn=? WHERE live_id=? AND turn=?", afterTurn, id, prevTurn)
	checkErr(err, "UPDATE BandMember QUERY failed.")
	return err
}

type PlayerRepositoryImpl struct {
	dbmap *gorp.DbMap
}

func NewPlayerRepositoryImpl(dbmap *gorp.DbMap) PlayerRepositoryImpl {
	dbmap.AddTableWithName(Player{}, "Player")
	return PlayerRepositoryImpl{dbmap: dbmap}
}

func (p PlayerRepositoryImpl) Create(player domain.Player) error {
	record := Player{
		Name: player.Name,
		Part: string(player.Part),
	}
	err := p.dbmap.Insert(&record)
	checkErr(err, "INSERT INTO Player QUERY failed.")
	return err
}

func (p PlayerRepositoryImpl) Delete(player domain.Player) error {
	record := Player{
		Name: player.Name,
		Part: string(player.Part),
	}
	_, err := p.dbmap.Delete(&record)
	checkErr(err, "INSERT INTO Player QUERY failed.")
	return err
}

func (p PlayerRepositoryImpl) FindByPart(part domain.Part) []domain.Player {
	var players []Player
	_, err := p.dbmap.Select(&players, "SELECT * FROM Player WHERE part = ?", string(part))
	checkErr(err, "SELECT * FROM Player QUERY failed.")

	var result []domain.Player
	for _, player := range players {
		result = append(result, player.ToModel())
	}
	return result
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
