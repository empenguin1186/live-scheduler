package model

import "time"

// Live ライブの構造体
type Live struct {
	// ライブID
	Id int
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
	Band []Band
}

// Band バンドの構造体
type Band struct {
	// バンド名
	Name string
	// ライブ ID
	LiveId int
	// 出演順
	Order int
	// メンバー
	Member []Member
}

// Part 楽器パート構造体
type Part string

const (
	Vo   = Part("Vo.")
	Gt   = Part("Gt.")
	GtVo = Part("Gt.Vo.")
	Key  = Part("Key.")
	Ba   = Part("Ba.")
	Dr   = Part("Dr.")
)

// Member Band メンバー構造体
type Member struct {
	Name string
	Part Part
}
