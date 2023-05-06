package imports

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neostefan/ga-diet/db"
)

func ImportDataFromCsvToDb() {
	sqlDB, errdb := sql.Open("sqlite3", "../db/meals.db")

	if errdb != nil {
		panic(errdb)
	}

	data := db.ReadFromCsvFile()

	db.ShiftToDb(data, sqlDB)

	defer sqlDB.Close()
}
