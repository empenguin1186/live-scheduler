package infra

import (
	"live-scheduler/domain"
	"time"
)

type Live struct {
	Id             int       `db:"id, primarykey, autoincrement"`
	Name           string    `db:"name"`
	Location       string    `db:"location"`
	Date           time.Time `db:"date"`
	PerformanceFee int       `db:"performance_fee"`
	EquipmentCost  int       `db:"equipment_cost"`
}

func (l *Live) ToModel() *domain.Live {
	return &domain.Live{Id: l.Id, Name: l.Name, Location: l.Location, Date: l.Date, PerformanceFee: l.PerformanceFee, EquipmentCost: l.EquipmentCost}
}

type Band struct {
	Name   string `db:"name"`
	LiveId int    `db:"live_id, primarykey"`
	Turn   int    `db:"turn"`
}

func (r *Band) ToModel() *domain.Band {
	return &domain.Band{Name: r.Name, LiveId: r.LiveId, Turn: r.Turn}
}

type BandMember struct {
	LiveId     int    `db:"live_id, primarykey"`
	Turn       int    `db:"turn, primarykey"`
	MemberName string `db:"member_name, primarykey"`
	MemberPart string `db:"member_part, primarykey"`
}

func (r *BandMember) ToModel() *domain.Player {
	return &domain.Player{Name: r.MemberName, Part: domain.Part(r.MemberPart)}
}

type Player struct {
	Name string `db:"name, primarykey"`
	Part string `db:"part, primarykey"`
}

func (p *Player) ToModel() *domain.Player {
	return &domain.Player{Name: p.Name, Part: domain.Part(p.Part)}
}
