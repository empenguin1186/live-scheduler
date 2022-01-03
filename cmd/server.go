package main

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"live-scheduler/infra/repository/record"
	"live-scheduler/presentation"
	"net/http"
	"os"
	"time"
)

const LAYOUT = "2006-01-02"

func main() {
	// 日付を指定してそれに基づくライブの情報を取得する
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// データベース接続処理
		user := os.Getenv("USER")
		pass := os.Getenv("PASS")
		dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/sample?parseTime=true", user, pass)
		db, err := sql.Open("mysql", dataSourceName)
		dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
		defer dbmap.Db.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Live テーブル select
		var lives []record.LiveRecord
		_, err = dbmap.Select(&lives, "SELECT * FROM Live WHERE date = ?", time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local).Format(LAYOUT))
		if err != nil {
			log.Fatal(err)
		}
		live := lives[0]

		// Band テーブル select
		var bands []record.BandRecord
		_, err = dbmap.Select(&bands, "SELECT * FROM Band WHERE live_id = ?", live.Id)
		if err != nil {
			log.Fatal(err)
		}

		// BandMember テーブル select
		var bandResponsePart []presentation.BandResponsePart
		for _, band := range bands {
			var memberResponsePart []presentation.MemberResponsePart
			var members []record.BandMemberRecord
			_, err = dbmap.Select(&members, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", band.LiveId, band.Turn)
			if err != nil {
				log.Fatal(err)
			}
			for _, member := range members {
				memberResponsePart = append(memberResponsePart, presentation.MemberResponsePart{Name: member.MemberName, Part: member.MemberPart})
			}
			bandResponsePart = append(bandResponsePart, presentation.BandResponsePart{Name: band.Name, Order: band.Turn, Member: memberResponsePart})
		}

		response := presentation.LiveResponse{
			Name:           live.Name,
			Location:       live.Location,
			Date:           live.Date,
			PerformanceFee: live.PerformanceFee,
			EquipmentCost:  live.EquipmentCost,
			//Band:           []presentation.BandResponsePart{{"hoge", 1, []presentation.MemberResponsePart{{"hoge", model.Ba}}}},
			Band: bandResponsePart,
		}
		return c.JSON(http.StatusOK, response)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
