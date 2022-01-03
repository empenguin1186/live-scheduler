package record

import (
	"live-scheduler/domain/model"
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

type BandRecord struct {
	Name   string `db:"name"`
	LiveId int    `db:"live_id, primarykey"`
	Turn   int    `db:"turn"`
}

type BandMemberRecord struct {
	LiveId     int        `db:"live_id, primarykey"`
	Turn       int        `db:"turn, primarykey"`
	MemberName string     `db:"member_name, primarykey"`
	MemberPart model.Part `db:"member_part, primarykey"`
}
