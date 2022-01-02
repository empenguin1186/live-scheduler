package main

import (
	"live-scheduler/domain/model"
	"live-scheduler/presentation"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	// 日付を指定してそれに基づくライブの情報を取得する
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		response := presentation.LiveResponse{
			Name:           "Name",
			Location:       "Location",
			Date:           time.Now(),
			PerformanceFee: 5500,
			EquipmentCost:  2000,
			Band:           []presentation.BandResponsePart{{"hoge", 1, []presentation.MemberResponsePart{{"hoge", model.Ba}}}},
		}
		return c.JSON(http.StatusOK, response)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
