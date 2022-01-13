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

func main() {
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/sample?parseTime=true", user, pass)
	db, err := sql.Open("mysql", dataSourceName)
	defer db.Close()
	if err != nil {
		log.Fatalln("db connection initialization failed.", err)
	}

	liveRepository := infra.NewLiveRepositoryImpl(db)
	bandRepository := infra.NewBandRepositoryImpl(db)
	bandMemberRepository := infra.NewBandMemberRepositoryImpl(db)
	liveDescService := domain.NewLiveDescServiceImpl(liveRepository, bandRepository, bandMemberRepository)
	liveService := domain.NewLiveServiceImpl(liveRepository)
	bandService := domain.NewBandServiceImpl(bandRepository)
	bandMemberService := domain.NewBandMemberServiceImpl(bandMemberRepository)

	e := echo.New()
	server := presentation.NewLiveServer(e, liveService, liveDescService, bandService, bandMemberService)
	server.Start()
}
