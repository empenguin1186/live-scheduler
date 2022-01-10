package presentation

import (
	"live-scheduler/domain"
	"time"
)

type LiveResponse struct {
	// ライブID
	Id int `json:"id"`
	// ライブ名
	Name string `json:"name"`
	// 場所
	Location string `json:"location"`
	// 日付
	Date time.Time `json:"date"`
	// 1人あたりの出演料
	PerformanceFee int `json:"performance_fee"`
	// 1バンドあたりの機材費
	EquipmentCost int `json:"equipment_cost"`
}

func (r LiveResponse) ToModel() *domain.Live {
	return &domain.Live{
		Id:             r.Id,
		Name:           r.Name,
		Location:       r.Location,
		Date:           r.Date,
		PerformanceFee: r.PerformanceFee,
		EquipmentCost:  r.EquipmentCost,
	}
}

func NewLiveResponse(live *domain.Live) *LiveResponse {
	return &LiveResponse{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
	}
}

type LiveDescResponse struct {
	// ライブ名
	Name string `json:"name"`
	// 場所
	Location string `json:"location"`
	// 日付
	Date time.Time `json:"date"`
	// 1人あたりの出演料
	PerformanceFee int `json:"performance_fee"`
	// 1バンドあたりの機材費
	EquipmentCost int `json:"equipment_cost"`
	// 出演するバンド
	Band []*BandResponsePart `json:"band,omitempty"`
}

func NewLiveDescResponse(liveModel *domain.LiveModel) *LiveDescResponse {
	bands := liveModel.Band
	var bandResponseParts []*BandResponsePart
	for _, band := range bands {
		var memberResponseParts []*MemberResponsePart
		for _, player := range band.Player {
			memberResponseParts = append(memberResponseParts, &MemberResponsePart{
				Name: player.Name,
				Part: player.Part,
			})
		}
		bandResponseParts = append(bandResponseParts, &BandResponsePart{
			Name:   band.Name,
			Turn:   band.Turn,
			Member: memberResponseParts,
		})
	}
	return &LiveDescResponse{
		Name:           liveModel.Name,
		Location:       liveModel.Location,
		Date:           liveModel.Date,
		PerformanceFee: liveModel.PerformanceFee,
		EquipmentCost:  liveModel.EquipmentCost,
		Band:           bandResponseParts,
	}
}

type BandResponsePart struct {
	// バンド名
	Name string `json:"name"`
	// 出演順
	Turn int `json:"turn"`
	// メンバー
	Member []*MemberResponsePart `json:"member,omitempty"`
}

type MemberResponsePart struct {
	Name string      `json:"name"`
	Part domain.Part `json:"part"`
}
