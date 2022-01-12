package presentation

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"net/http"
	"time"
)

type LiveCreateRequest struct {
	// ライブ名
	Name string `json:"name" validate:"required"`
	// 場所
	Location string `json:"location" validate:"required"`
	// 日付
	Date time.Time `json:"date" validate:"required"`
	// 1人あたりの出演料
	PerformanceFee int `json:"performance_fee" validate:"required"`
	// 1バンドあたりの機材費
	EquipmentCost int `json:"equipment_cost" validate:"required"`
}

func (r LiveCreateRequest) ToModel() *domain.Live {
	return &domain.Live{
		Name:           r.Name,
		Location:       r.Location,
		Date:           r.Date,
		PerformanceFee: r.PerformanceFee,
		EquipmentCost:  r.EquipmentCost,
	}
}

type LivePatchRequest struct {
	// ライブID
	Id int `json:"id" validate:"required"`
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

func (r LivePatchRequest) ToModel() *domain.Live {
	return &domain.Live{
		Id:             r.Id,
		Name:           r.Name,
		Location:       r.Location,
		Date:           r.Date,
		PerformanceFee: r.PerformanceFee,
		EquipmentCost:  r.EquipmentCost,
	}
}

type BandCreateRequest struct {
	LiveId int    `json:"live_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Turn   int    `json:"turn" validate:"required"`
}

func (r BandCreateRequest) ToModel() *domain.Band {
	return &domain.Band{
		Name:   r.Name,
		LiveId: r.LiveId,
		Turn:   r.Turn,
	}
}

type BandPatchRequest struct {
	LiveId int    `json:"live_id"`
	Name   string `json:"name"`
	Turn   int    `json:"turn"`
}

func (r BandPatchRequest) ToModel() *domain.Band {
	return &domain.Band{
		Name:   r.Name,
		LiveId: r.LiveId,
		Turn:   r.Turn,
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
