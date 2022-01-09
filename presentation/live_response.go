package presentation

import (
	"live-scheduler/domain"
	"time"
)

type LiveResponse struct {
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
	Band []*BandResponsePart `json:"band"`
}

func NewLiveResponse(liveModel *domain.LiveModel) *LiveResponse {
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
	return &LiveResponse{
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
	Member []*MemberResponsePart `json:"member"`
}

type MemberResponsePart struct {
	Name string      `json:"name"`
	Part domain.Part `json:"part"`
}
