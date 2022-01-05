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
		//playerRepository := infra.NewPlayerRepositoryImpl(dbmap)
		liveService := domain.NewLiveService(liveRepository, bandRepository, bandMemberRepository)

		var date time.Time
		err = echo.QueryParamsBinder(c).Time("date", &date, LAYOUT).BindError()
		checkErr(err, "Invalid Query Parameter")
		liveModel := liveService.GetByDate(&date)

		//response := presentation.LiveResponse{
		//	Name:           live.Name,
		//	Location:       live.Location,
		//	Date:           live.Date,
		//	PerformanceFee: live.PerformanceFee,
		//	EquipmentCost:  live.EquipmentCost,
		//	Band:           live.Band,
		//}

		return c.JSON(http.StatusOK, liveModel)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
