package infra

import (
	"live-scheduler/domain"
	"time"
)

type LiveRecord struct {
	Id             int       `db:"id, primarykey, autoincrement"`
	Name           string    `db:"name"`
	Location       string    `db:"location"`
	Date           time.Time `db:"date"`
	PerformanceFee int       `db:"performance_fee"`
	EquipmentCost  int       `db:"equipment_cost"`
}

type LiveRecordForInsert struct {
	Name           string    `db:"name"`
	Location       string    `db:"location"`
	Date           time.Time `db:"date"`
	PerformanceFee int       `db:"performance_fee"`
	EquipmentCost  int       `db:"equipment_cost"`
}

func (r LiveRecord) ToModel() domain.Live {
	return domain.Live{Id: r.Id, Name: r.Name, Location: r.Location, Date: r.Date, PerformanceFee: r.PerformanceFee, EquipmentCost: r.EquipmentCost}
}

type BandRecord struct {
	Name   string `db:"name"`
	LiveId int    `db:"live_id, primarykey"`
	Turn   int    `db:"turn"`
}

func (r BandRecord) ToModel() domain.Band {
	return domain.Band{Name: r.Name, LiveId: r.LiveId, Turn: r.Turn}
}

type BandMemberRecord struct {
	LiveId     int         `db:"live_id, primarykey"`
	Turn       int         `db:"turn, primarykey"`
	MemberName string      `db:"member_name, primarykey"`
	MemberPart domain.Part `db:"member_part, primarykey"`
}

func (r BandMemberRecord) ToModel() domain.Player {
	return domain.Player{Name: r.MemberName, Part: r.MemberPart}
}
