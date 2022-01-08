package infra

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	"live-scheduler/domain"
	"log"
	"time"
)

type Dao interface {
	AddTableWithName(i interface{}, name string) *gorp.TableMap
	Insert(list ...interface{}) error
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type LiveRepositoryAlpha struct {
	db *sql.DB
}

func NewLiveRepositoryAlpha(db *sql.DB) *LiveRepositoryAlpha {
	return &LiveRepositoryAlpha{db: db}
}

func (a *LiveRepositoryAlpha) FindByDate(date *time.Time) *domain.Live {
	rows, err := a.db.Query(`SELECT * FROM Live WHERE date = ?`, date.Format("2006-01-02"))
	if err != nil {
		log.Fatal(err)
	}
	var lives []Live
	for rows.Next() {
		var id int
		var name string
		var location string
		var date time.Time
		var performanceFee int
		var equipmentCost int

		err = rows.Scan(&id, &name, &location, &date, &performanceFee, &equipmentCost)
		lives = append(lives, Live{Id: id, Name: name, Location: location, Date: date, PerformanceFee: performanceFee, EquipmentCost: equipmentCost})
	}
	live := lives[0]
	return live.ToModel()
}

func (a *LiveRepositoryAlpha) Update(live *domain.Live) error {
	panic("implement me")
}

func (a *LiveRepositoryAlpha) Delete(live *domain.Live) error {
	panic("implement me")
}

type LiveRepositoryImpl struct {
	dao Dao
}

func NewLiveRepositoryImpl(dao Dao) *LiveRepositoryImpl {
	dao.AddTableWithName(Live{}, "Live")
	return &LiveRepositoryImpl{dao: dao}
}

func (i *LiveRepositoryImpl) FindByDate(time *time.Time) *domain.Live {
	var lives []Live
	_, err := i.dao.Select(&lives, "SELECT * FROM Live WHERE date = ?", time.Format("2006-01-02"))
	checkErr(err, "SELECT * FROM Live QUERY failed.")
	live := lives[0]
	return live.ToModel()
}

func (i *LiveRepositoryImpl) Create(live *domain.Live) error {
	record := Live{
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	err := i.dao.Insert(&record)
	checkErr(err, "INSERT INTO Live QUERY failed.")
	return err
}

func (i *LiveRepositoryImpl) Update(live *domain.Live) error {
	record := Live{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	_, err := i.dao.Update(&record)
	checkErr(err, "UPDATE Live QUERY failed.")
	return err
}

func (i *LiveRepositoryImpl) Delete(live *domain.Live) error {
	record := Live{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
	_, err := i.dao.Delete(&record)
	checkErr(err, "DELETE Live QUERY failed.")
	return err
}

type BandRepositoryImpl struct {
	dao Dao
}

func NewBandRepositoryImpl(dao Dao) *BandRepositoryImpl {
	dao.AddTableWithName(Band{}, "Band")
	return &BandRepositoryImpl{dao: dao}
}

func (i BandRepositoryImpl) FindByLiveId(id int) []*domain.Band {
	var bandRecords []Band
	_, err := i.dao.Select(&bandRecords, "SELECT * FROM Band WHERE live_id = ?", id)
	checkErr(err, "SELECT * FROM Band QUERY failed.")

	var bands []*domain.Band
	for _, bandRecord := range bandRecords {
		bands = append(bands, bandRecord.ToModel())
	}
	return bands
}

func (i *BandRepositoryImpl) Create(band *domain.Band) error {
	record := Band{
		Name:   band.Name,
		LiveId: band.LiveId,
		Turn:   band.Turn,
	}
	err := i.dao.Insert(&record)
	checkErr(err, "INSERT INTO Band QUERY failed.")
	return err
}

func (i *BandRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := i.dao.Exec("UPDATE Band SET turn=? WHERE live_id=? AND turn=?", afterTurn, id, prevTurn)
	checkErr(err, "UPDATE Band QUERY failed.")
	return err
}

func (i *BandRepositoryImpl) DeleteByIdAndTurn(id int, turn int) error {
	_, err := i.dao.Exec("DELETE FROM Band WHERE live_id=? AND turn=?", id, turn)
	checkErr(err, "DELETE Band QUERY failed.")
	return err
}

type BandMemberRepositoryImpl struct {
	dao Dao
}

func NewBandMemberRepositoryImpl(dao Dao) *BandMemberRepositoryImpl {
	dao.AddTableWithName(BandMember{}, "BandMember")
	return &BandMemberRepositoryImpl{dao: dao}
}

func (i *BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) []*domain.Player {
	var memberRecords []BandMember
	_, err := i.dao.Select(&memberRecords, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", id, turn)
	checkErr(err, "SELECT * FROM BandMember QUERY failed.")

	var members []*domain.Player
	for _, memberRecord := range memberRecords {
		members = append(members, memberRecord.ToModel())
	}
	return members
}

func (i *BandMemberRepositoryImpl) Create(bandMember *domain.BandMember) error {
	record := BandMember{
		LiveId:     bandMember.LiveId,
		Turn:       bandMember.Turn,
		MemberName: bandMember.MemberName,
		MemberPart: string(bandMember.MemberPart),
	}
	err := i.dao.Insert(&record)
	checkErr(err, "INSERT INTO BandMember QUERY failed.")
	return err
}

func (i *BandMemberRepositoryImpl) Delete(bandMember *domain.BandMember) error {
	record := BandMember{
		LiveId:     bandMember.LiveId,
		Turn:       bandMember.Turn,
		MemberName: bandMember.MemberName,
		MemberPart: string(bandMember.MemberPart),
	}
	_, err := i.dao.Delete(&record)
	checkErr(err, "INSERT INTO BandMember QUERY failed.")
	return err
}

func (i *BandMemberRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := i.dao.Exec("UPDATE BandMember SET turn=? WHERE live_id=? AND turn=?", afterTurn, id, prevTurn)
	checkErr(err, "UPDATE BandMember QUERY failed.")
	return err
}

type PlayerRepositoryImpl struct {
	dao Dao
}

func NewPlayerRepositoryImpl(dao Dao) *PlayerRepositoryImpl {
	dao.AddTableWithName(Player{}, "Player")
	return &PlayerRepositoryImpl{dao: dao}
}

func (p *PlayerRepositoryImpl) Create(player *domain.Player) error {
	record := Player{
		Name: player.Name,
		Part: string(player.Part),
	}
	err := p.dao.Insert(&record)
	checkErr(err, "INSERT INTO Player QUERY failed.")
	return err
}

func (p *PlayerRepositoryImpl) Delete(player *domain.Player) error {
	record := Player{
		Name: player.Name,
		Part: string(player.Part),
	}
	_, err := p.dao.Delete(&record)
	checkErr(err, "INSERT INTO Player QUERY failed.")
	return err
}

func (p *PlayerRepositoryImpl) FindByPart(part *domain.Part) []*domain.Player {
	var players []Player
	_, err := p.dao.Select(&players, "SELECT * FROM Player WHERE part = ?", string(*part))
	checkErr(err, "SELECT * FROM Player QUERY failed.")

	var result []*domain.Player
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
