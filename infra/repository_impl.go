package infra

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	"live-scheduler/domain"
	"log"
	"time"
)

const LAYOUT = "2006-01-02"

type Dao interface {
	AddTableWithName(i interface{}, name string) *gorp.TableMap
	Insert(list ...interface{}) error
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type LiveRepositoryImpl struct {
	db *sql.DB
}

func NewLiveRepositoryImpl(db *sql.DB) *LiveRepositoryImpl {
	return &LiveRepositoryImpl{db: db}
}

func (a *LiveRepositoryImpl) FindByDate(date *time.Time) *domain.Live {
	rows, err := a.db.Query(`SELECT * FROM Live WHERE date = ?`, date.Format(LAYOUT))
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

func (a *LiveRepositoryImpl) Create(live *domain.Live) error {
	_, err := a.db.Exec(
		`INSERT INTO Live(name, location, date, performance_fee, equipment_cost) VALUES ( ?, ?, ?, ?, ? )`,
		live.Name, live.Location, live.Date.Format(LAYOUT), live.PerformanceFee, live.EquipmentCost)
	return err
}

func (a *LiveRepositoryImpl) Update(live *domain.Live) error {
	_, err := a.db.Exec(
		`UPDATE Live SET name = ?, location = ?, date = ?, performance_fee = ?, equipment_cost = ?`,
		live.Name, live.Location, live.Date.Format(LAYOUT), live.PerformanceFee, live.EquipmentCost)
	return err
}

func (a *LiveRepositoryImpl) Delete(live *domain.Live) error {
	_, err := a.db.Exec(`DELETE FROM Live WHERE id = ?`, live.Id)
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
