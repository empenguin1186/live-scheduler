package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	repository2 "live-scheduler/infra/repository"
	"net/http"
	"time"
)

const LAYOUT = "2006-01-02"

func main() {
	// 日付を指定してそれに基づくライブの情報を取得する
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		repository := repository2.NewLiveRepositoryImpl()
		live := repository.FindByDate(time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local))

		//response := presentation.LiveResponse{
		//	Name:           live.Name,
		//	Location:       live.Location,
		//	Date:           live.Date,
		//	PerformanceFee: live.PerformanceFee,
		//	EquipmentCost:  live.EquipmentCost,
		//	Band:           live.Band,
		//}
		return c.JSON(http.StatusOK, live)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
