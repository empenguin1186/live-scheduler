package infra

import (
	"database/sql"
	"live-scheduler/domain"
	"time"
)

const LAYOUT = "2006-01-02"

type LiveRepositoryImpl struct {
	db *sql.DB
}

func NewLiveRepositoryImpl(db *sql.DB) *LiveRepositoryImpl {
	return &LiveRepositoryImpl{db: db}
}

func (i *LiveRepositoryImpl) FindById(id int) (*domain.Live, error) {
	var live domain.Live
	err := i.db.QueryRow(`SELECT * FROM Live WHERE id = ?`, id).
		Scan(&live.Id, &live.Name, &live.Location, &live.Date, &live.PerformanceFee, &live.EquipmentCost)
	if err != nil {
		return nil, err
	}
	return &live, nil
}

func (i *LiveRepositoryImpl) FindByPeriod(start *time.Time, end *time.Time) ([]*domain.Live, error) {
	rows, err := i.db.Query(
		`SELECT * FROM Live WHERE date >= ? AND date <= ?`,
		start.Format(LAYOUT), end.Format(LAYOUT))
	if err != nil {
		return nil, err
	}
	var lives []*domain.Live
	for rows.Next() {
		var id, performanceFee, equipmentCost int
		var name, location string
		var date time.Time

		err = rows.Scan(&id, &name, &location, &date, &performanceFee, &equipmentCost)
		lives = append(lives, &domain.Live{
			Id:             id,
			Name:           name,
			Location:       location,
			Date:           date,
			PerformanceFee: performanceFee,
			EquipmentCost:  equipmentCost,
		})
	}
	return lives, nil
}

func (i *LiveRepositoryImpl) Create(live *domain.Live) error {
	_, err := i.db.Exec(
		`INSERT INTO Live(name, location, date, performance_fee, equipment_cost) VALUES ( ?, ?, ?, ?, ? )`,
		live.Name, live.Location, live.Date.Format(LAYOUT), live.PerformanceFee, live.EquipmentCost)
	return err
}

func (i *LiveRepositoryImpl) Update(live *domain.Live) error {
	_, err := i.db.Exec(
		`UPDATE Live SET name = ?, location = ?, date = ?, performance_fee = ?, equipment_cost = ? WHERE id = ?`,
		live.Name, live.Location, live.Date.Format(LAYOUT), live.PerformanceFee, live.EquipmentCost, live.Id)
	return err
}

func (i *LiveRepositoryImpl) Delete(id int) error {
	_, err := i.db.Exec(`DELETE FROM Live WHERE id = ?`, id)
	return err
}

type BandRepositoryImpl struct {
	db *sql.DB
}

func NewBandRepositoryImpl(db *sql.DB) *BandRepositoryImpl {
	return &BandRepositoryImpl{db: db}
}

func (b *BandRepositoryImpl) FindByLiveId(id int) ([]*domain.Band, error) {
	rows, err := b.db.Query(`SELECT * FROM Band WHERE live_id = ?`, id)
	if err != nil {
		return nil, err
	}
	var bands []*domain.Band
	for rows.Next() {
		var name string
		var liveId, turn int

		err = rows.Scan(&name, &liveId, &turn)
		if err != nil {
			return nil, err
		}
		band := domain.Band{Name: name, LiveId: liveId, Turn: turn}
		bands = append(bands, &band)
	}
	return bands, nil
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

func (b *BandMemberRepositoryImpl) FindByLiveIdAndTurn(id int, turn int) ([]*domain.Player, error) {
	rows, err := b.db.Query(`SELECT * FROM BandMember WHERE live_id = ? AND turn = ?`, id, turn)
	if err != nil {
		return nil, err
	}
	var players []*domain.Player
	for rows.Next() {
		var liveId, turn int
		var name, part string

		err = rows.Scan(&liveId, &turn, &name, &part)
		if err != nil {
			return nil, err
		}
		player := domain.Player{Name: name, Part: domain.Part(part)}
		players = append(players, &player)
	}
	return players, nil
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

func (p *PlayerRepositoryImpl) FindByPart(part *domain.Part) ([]*domain.Player, error) {
	rows, err := p.db.Query(`SELECT * FROM Player WHERE Part = ?`, string(*part))
	if err != nil {
		return nil, err
	}
	var players []*domain.Player
	for rows.Next() {
		var name, part string

		err = rows.Scan(&name, &part)
		if err != nil {
			return nil, err
		}
		player := domain.Player{Name: name, Part: domain.Part(part)}
		players = append(players, &player)
	}
	return players, err
}
