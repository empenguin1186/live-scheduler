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
	Band []BandResponsePart `json:"band"`
}

type BandResponsePart struct {
	// バンド名
	Name string `json:"name"`
	// 出演順
	Order int `json:"order"`
	// メンバー
	Member []MemberResponsePart `json:"member"`
}

type MemberResponsePart struct {
	Name string      `json:"name"`
	Part domain.Part `json:"part"`
}
