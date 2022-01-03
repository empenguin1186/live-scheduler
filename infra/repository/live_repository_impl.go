package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"live-scheduler/domain/model"
	"live-scheduler/infra/repository/record"
	"log"
	"os"
	"time"
)

const LAYOUT = "2006-01-02"

type LiveRepositoryImpl struct {
	db *sql.DB
}

func NewLiveRepositoryImpl() LiveRepositoryImpl {
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/sample?parseTime=true", user, pass)
	db, err := sql.Open("mysql", dataSourceName)
	checkErr(err, "db connection initialization failed.")
	return LiveRepositoryImpl{db: db}
}

func (lri LiveRepositoryImpl) FindByDate(time time.Time) model.Live {
	dbmap := &gorp.DbMap{Db: lri.db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	defer dbmap.Db.Close()

	// Live テーブル select
	var lives []record.LiveRecord
	_, err := dbmap.Select(&lives, "SELECT * FROM Live WHERE date = ?", time.Format(LAYOUT))
	checkErr(err, "SELECT * FROM Live QUERY failed.")
	live := lives[0]

	// Band テーブル select
	var bandRecords []record.BandRecord
	_, err = dbmap.Select(&bandRecords, "SELECT * FROM Band WHERE live_id = ?", live.Id)
	checkErr(err, "SELECT * FROM Band QUERY failed.")

	// BandMember テーブル select
	var bands []model.Band
	for _, bandRecord := range bandRecords {
		var members []model.Member
		var memberRecords []record.BandMemberRecord
		_, err = dbmap.Select(&memberRecords, "SELECT * FROM BandMember WHERE live_id = ? AND turn = ?", bandRecord.LiveId, bandRecord.Turn)
		checkErr(err, "SELECT * FROM BandMember QUERY failed.")
		for _, memberRecord := range memberRecords {
			members = append(members, model.Member{Name: memberRecord.MemberName, Part: memberRecord.MemberPart})
		}
		bands = append(bands, model.Band{Name: bandRecord.Name, Order: bandRecord.Turn, Member: members})
	}
	return model.Live{Id: live.Id, Name: live.Name, Location: live.Location, Date: live.Date, PerformanceFee: live.PerformanceFee, EquipmentCost: live.EquipmentCost, Band: bands}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
