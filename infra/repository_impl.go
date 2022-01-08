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
		`UPDATE Live SET name = ?, location = ?, date = ?, performance_fee = ?, equipment_cost = ? WHERE id = ?`,
		live.Name, live.Location, live.Date.Format(LAYOUT), live.PerformanceFee, live.EquipmentCost, live.Id)
	return err
}

func (a *LiveRepositoryImpl) Delete(live *domain.Live) error {
	_, err := a.db.Exec(`DELETE FROM Live WHERE id = ?`, live.Id)
	return err
}

type BandRepositoryImpl struct {
	db *sql.DB
}

func NewBandRepositoryImpl(db *sql.DB) *BandRepositoryImpl {
	return &BandRepositoryImpl{db: db}
}

func (b *BandRepositoryImpl) FindByLiveId(id int) []*domain.Band {
	rows, err := b.db.Query(`SELECT * FROM Band WHERE live_id = ?`, id)
	if err != nil {
		log.Fatal(err)
	}
	var bands []*domain.Band
	for rows.Next() {
		var name string
		var liveId int
		var turn int

		// TODO エラーハンドリング
		err = rows.Scan(&name, &liveId, &turn)
		band := Band{Name: name, LiveId: liveId, Turn: turn}
		bands = append(bands, band.ToModel())
	}
	return bands
}

func (b *BandRepositoryImpl) Create(band *domain.Band) error {
	_, err := b.db.Exec(
		`INSERT INTO Band(name, live_id, turn) VALUES ( ?, ?, ? )`,
		band.Name, band.LiveId, band.Turn)
	return err
}

func (b *BandRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := b.db.Exec(
		`UPDATE Band SET turn = ? WHERE live_id = ? AND turn = ?`,
		afterTurn, id, prevTurn)
	return err
}

func (b *BandRepositoryImpl) DeleteByIdAndTurn(id int, turn int) error {
	_, err := b.db.Exec(`DELETE FROM Band WHERE live_id = ? AND turn = ?`, id, turn)
	return err
}

type BandMemberRepositoryImpl struct {
	db *sql.DB
}

func NewBandMemberRepositoryImpl(db *sql.DB) *BandMemberRepositoryImpl {
	return &BandMemberRepositoryImpl{db: db}
}

func (b *BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) []*domain.Player {
	rows, err := b.db.Query(`SELECT * FROM BandMember WHERE live_id = ? AND turn = ?`, id, turn)
	// TODO エラーハンドリング
	if err != nil {
		log.Fatal(err)
	}
	var players []*domain.Player
	for rows.Next() {
		var liveId int
		var turn int
		var name string
		var part string

		// TODO エラーハンドリング
		err = rows.Scan(&liveId, &turn, &name, &part)
		player := domain.Player{Name: name, Part: domain.Part(part)}
		players = append(players, &player)
	}
	return players
}

func (b *BandMemberRepositoryImpl) Create(bandMember *domain.BandMember) error {
	_, err := b.db.Exec(
		`INSERT INTO BandMember(live_id, turn, member_name, member_name) VALUES ( ?, ?, ?, ? )`,
		bandMember.LiveId, bandMember.Turn, bandMember.MemberName, string(bandMember.MemberPart))
	return err
}

func (b *BandMemberRepositoryImpl) Delete(bandMember *domain.BandMember) error {
	_, err := b.db.Exec(
		`DELETE FROM BandMember WHERE live_id = ? AND turn = ? AND member_name = ? AND member_name = ?`,
		bandMember.LiveId, bandMember.Turn, bandMember.MemberName, string(bandMember.MemberPart))
	return err
}

func (b *BandMemberRepositoryImpl) UpdateTurn(id int, prevTurn int, afterTurn int) error {
	_, err := b.db.Exec(`UPDATE BandMember SET turn = ? WHERE live_id = ? AND turn = ?`, afterTurn, id, prevTurn)
	return err
}

type PlayerRepositoryImpl struct {
	db *sql.DB
}

func NewPlayerRepositoryImpl(db *sql.DB) *PlayerRepositoryImpl {
	return &PlayerRepositoryImpl{db: db}
}

func (p *PlayerRepositoryImpl) Create(player *domain.Player) error {
	_, err := p.db.Exec(`INSERT INTO Player(name, part) VALUES ( ?, ? )`, player.Name, string(player.Part))
	return err
}

func (p *PlayerRepositoryImpl) Delete(player *domain.Player) error {
	_, err := p.db.Exec(`DELETE FROM Player WHERE name = ? AND part = ?`, player.Name, string(player.Part))
	return err
}

func (p *PlayerRepositoryImpl) FindByPart(part *domain.Part) []*domain.Player {
	rows, err := p.db.Query(`SELECT * FROM Player WHERE Part = ?`, string(*part))
	// TODO エラーハンドリング
	if err != nil {
		log.Fatal(err)
	}
	var players []*domain.Player
	for rows.Next() {
		var name string
		var part string

		// TODO エラーハンドリング
		err = rows.Scan(&name, &part)
		player := domain.Player{Name: name, Part: domain.Part(part)}
		players = append(players, &player)
	}
	return players
}
