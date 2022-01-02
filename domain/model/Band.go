package model

import "time"

// Live ライブの構造体
type Live struct {
	// ライブID
	id int
	// ライブ名
	name string
	// 場所
	location string
	// 日付
	date time.Time
	// 1人あたりの出演料
	performanceFee int
	// 1バンドあたりの機材費
	equipmentCost int
	// 出演するバンド
	band []Band
}

// Band バンドの構造体
type Band struct {
	// バンド名
	name string
	// ライブ ID
	liveId int
	// 出演順
	order int
	// メンバー
	member []Member
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
	name string
	part Part
}
