package presentation

import (
	"time"
)

type LiveResponse struct {
	// ライブ名
	Name string
	// 場所
	Location string
	// 日付
	Date time.Time
	// 1人あたりの出演料
	PerformanceFee int
	// 1バンドあたりの機材費
	EquipmentCost int
	// 出演するバンド
	Band []BandResponsePart
}

type BandResponsePart struct {
	// バンド名
	Name string
	// 出演順
	Order int
	// メンバー
	Member []MemberResponsePart
}

type MemberResponsePart struct {
	Name string
	Part string
}
