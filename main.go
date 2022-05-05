package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neostefan/diet-assistant/definitions"
	"github.com/neostefan/diet-assistant/operators"
)

func main() {
	//ings := db.ReadFromCsvFile()

	sqlDb, err := sql.Open("sqlite3", "./db/meals.db")

	if err != nil {
		fmt.Printf("Error occurred: %v", err)
	}

	//db.ShiftToDb(ings, sqlDb)


	
	population := operators.InitializePopulation(sqlDb)
	c1 := population[0]
	c2 := population[1]
	fmt.Println("Printing the initial states...")
	fmt.Printf("%d, %d \n", c1, c2)
	c1.Evolve(&c2, definitions.Crossover, 1)
	fmt.Println("Printing the mutated states...")
	fmt.Printf("%d, %d \n", c1, c2)

	defer sqlDb.Close()
}