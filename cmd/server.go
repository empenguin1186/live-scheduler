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
	"strconv"
	"time"
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
	e.Validator = presentation.NewCustomValidator()
	e.GET("/live", func(context echo.Context) error {
		var start time.Time
		err = echo.QueryParamsBinder(context).Time("start", &start, LAYOUT).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var end time.Time
		err = echo.QueryParamsBinder(context).Time("end", &end, LAYOUT).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		lives, err := liveService.GetByPeriod(&start, &end)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		var liveResponse []*presentation.LiveResponse
		for _, e := range lives {
			liveResponse = append(liveResponse, presentation.NewLiveResponse(e))
		}
		return context.JSON(http.StatusOK, liveResponse)
	})

	e.GET("/live/:id", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		liveModel, err := liveDescService.GetById(int(liveId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, presentation.NewLiveDescResponse(liveModel))
	})

	e.POST("/live", func(context echo.Context) error {
		live := new(presentation.LiveCreateRequest)
		if err = context.Bind(live); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = context.Validate(live); err != nil {
			return err
		}
		err := liveService.Register(live.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, live)
	})

	e.PATCH("/live", func(context echo.Context) error {
		live := new(presentation.LivePatchRequest)
		if err = context.Bind(live); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = context.Validate(live); err != nil {
			return err
		}
		err := liveService.Update(live.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, live)
	})

	e.DELETE("/live/:id", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = liveService.Delete(int(liveId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
