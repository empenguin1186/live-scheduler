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
	"os"
)

const LAYOUT = "2006-01-02"

func main() {
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/sample?parseTime=true", user, pass)
	db, err := sql.Open("mysql", dataSourceName)
	defer db.Close()
	checkErr(err, "db connection initialization failed.")

	liveRepository := infra.NewLiveRepositoryImpl(db)
	bandRepository := infra.NewBandRepositoryImpl(db)
	bandMemberRepository := infra.NewBandMemberRepositoryImpl(db)
	liveDescService := domain.NewLiveDescServiceImpl(liveRepository, bandRepository, bandMemberRepository)
	liveService := domain.NewLiveServiceImpl(liveRepository)

	e := echo.New()
	server := presentation.NewLiveServer(e, liveService, liveDescService)
	server.Start()
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
