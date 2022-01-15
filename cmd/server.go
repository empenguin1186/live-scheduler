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
	playerRepository := infra.NewPlayerRepositoryImpl(db)

	liveDescService := domain.NewLiveDescServiceImpl(liveRepository, bandRepository, bandMemberRepository)
	liveService := domain.NewLiveServiceImpl(liveRepository)
	bandService := domain.NewBandServiceImpl(bandRepository)
	bandMemberService := domain.NewBandMemberServiceImpl(bandMemberRepository)
	playerService := domain.NewPlayerServiceImpl(playerRepository)

	e := echo.New()
	handler := presentation.NewLiveHandler(liveService, liveDescService, bandService, bandMemberService, playerService)
	e.Validator = presentation.NewCustomValidator()

	e.GET("/live", handler.GetLives)
	e.GET("/live/:id", handler.GetLive)
	e.POST("/live", handler.PostLive)
	e.PATCH("/live", handler.PatchLive)
	e.DELETE("/live/:id", handler.DeleteLive)

	e.GET("/live/:id/band", handler.GetBand)
	e.POST("/live/:id/band", handler.PostBand)
	e.PATCH("/live/:live_id/band/:turn", handler.PatchBand)
	e.DELETE("/live/:live_id/band/:turn", handler.DeleteBand)

	e.GET("/member", handler.GetPart)
	e.POST("/member/create", handler.PostPart)
	e.POST("/member/delete", handler.DeletePart)

	e.Logger.Fatal(e.Start(":1323"))
}
