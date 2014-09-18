package gofx

import (
	"database/sql"
	"github.com/jmoiron/modl"
	"log"
)

var dbMap *modl.DbMap

func InitDb() {
	// Connect DB
	db, err := sql.Open("mysql", "andy:@/gofx")
	err = db.Ping()
	checkErr(err, "SQL Connection failed")

	// Get db map, register and create table
	dbMap = modl.NewDbMap(db, modl.MySQLDialect{"InnoDB", "UTF8"})
	dbMap.AddTableWithName(Order{}, "orders").SetKeys(false, "Timestamp")
	err = dbMap.CreateTablesIfNotExists()
	checkErr(err, "Table Creation failed")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
