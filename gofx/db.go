package gofx

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
)

var dbMap *gorp.DbMap

func InitDb() {
	// Connect DB
	db, err := sql.Open("mysql", "andy:@/gofx")
	err = db.Ping()
	checkErr(err, "SQL Connection failed")

	// Get db map, register and create table
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbMap.AddTableWithName(Order{}, "orders").SetKeys(false, "ts")
	dbMap.AddTableWithName(Account{}, "accounts").SetKeys(false, "id")
	err = dbMap.CreateTablesIfNotExists()
	checkErr(err, "Table Creation failed")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
