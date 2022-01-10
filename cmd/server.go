package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"live-scheduler/infra"
	"live-scheduler/presentation"
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
		defer db.Close()
		checkErr(err, "db connection initialization failed.")

		liveRepository := infra.NewLiveRepositoryImpl(db)
		bandRepository := infra.NewBandRepositoryImpl(db)
		bandMemberRepository := infra.NewBandMemberRepositoryImpl(db)
		//playerRepository := infra.NewPlayerRepositoryImpl(db)
		liveService := domain.NewLiveDescServiceImpl(liveRepository, bandRepository, bandMemberRepository)

		var date time.Time
		err = echo.QueryParamsBinder(c).Time("date", &date, LAYOUT).BindError()
		checkErr(err, "Invalid Query Parameter")
		liveModel, err := liveService.GetByDate(&date)
		return c.JSON(http.StatusOK, presentation.NewLiveResponse(liveModel))
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
