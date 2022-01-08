package infra

import (
	"live-scheduler/domain"
	"time"
)

type Live struct {
	Id             int
	Name           string
	Location       string
	Date           time.Time
	PerformanceFee int
	EquipmentCost  int
}

func (l *Live) ToModel() *domain.Live {
	return &domain.Live{Id: l.Id, Name: l.Name, Location: l.Location, Date: l.Date, PerformanceFee: l.PerformanceFee, EquipmentCost: l.EquipmentCost}
}

type Band struct {
	Name   string
	LiveId int
	Turn   int
}

func (r *Band) ToModel() *domain.Band {
	return &domain.Band{Name: r.Name, LiveId: r.LiveId, Turn: r.Turn}
}

type BandMember struct {
	LiveId     int
	Turn       int
	MemberName string
	MemberPart string
}

func (r *BandMember) ToModel() *domain.Player {
	return &domain.Player{Name: r.MemberName, Part: domain.Part(r.MemberPart)}
}

type Player struct {
	Name string
	Part string
}

func (p *Player) ToModel() *domain.Player {
	return &domain.Player{Name: p.Name, Part: domain.Part(p.Part)}
}
