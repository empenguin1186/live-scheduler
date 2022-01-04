package main

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"live-scheduler/infra"
	"log"
	"net/http"
	"os"
	"time"
)

const LAYOUT = "2006-01-02"

func main() {
	// 日付を指定してそれに基づくライブの情報を取得する
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		user := os.Getenv("USER")
		pass := os.Getenv("PASS")
		dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/sample?parseTime=true", user, pass)
		db, err := sql.Open("mysql", dataSourceName)
		dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
		defer dbmap.Db.Close()
		checkErr(err, "db connection initialization failed.")

		liveRepository := infra.NewLiveRepositoryImpl(dbmap)
		bandRepository := infra.NewBandRepositoryImpl(dbmap)
		bandMemberRepository := infra.NewBandMemberRepositoryImpl(dbmap)
		liveService := domain.NewLiveService(liveRepository, bandRepository, bandMemberRepository)
		live := liveService.GetDate(time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local))

		//live := repository.FindByDate(time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local))

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

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
