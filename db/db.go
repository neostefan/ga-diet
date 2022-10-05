package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/neostefan/ga-diet/definitions"
)

func CreateTables(db *sql.DB, t string) {
	var stmt *sql.Stmt
	var err error

	if t == "carbs" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS carbs (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if t == "protein" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS proteins (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if t == "oil" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS oils (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if t == "vegetables" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS vegetables (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if t == "beverage" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS beverages (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if t == "fruit" {
		stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS fruits (id INTEGER PRIMARY KEY, name VARCHAR(255), calories FLOAT, protein FLOAT, cost FLOAT, type TEXT)`)
	}

	if err != nil {
		fmt.Printf("Error occurred in initial table creation: %v", err)
	}

	defer stmt.Close()

	_, errE := stmt.Exec()

	if errE != nil {
		fmt.Printf("Error occurred in table creation: %v", errE)
	}

	defer stmt.Close()
}

func InsertIngredient(db *sql.DB, i *definitions.IngredientDetails) {

	var stmt *sql.Stmt
	var err error

	if i.Type == "carbs" {
		stmt, err = db.Prepare(`INSERT INTO carbs (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if i.Type == "protein" {
		stmt, err = db.Prepare(`INSERT INTO proteins (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if i.Type == "oil" {
		stmt, err = db.Prepare(`INSERT INTO oils (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if i.Type == "vegetables" {
		stmt, err = db.Prepare(`INSERT INTO vegetables (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if i.Type == "beverage" {
		stmt, err = db.Prepare(`INSERT INTO beverages (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if i.Type == "fruit" {
		stmt, err = db.Prepare(`INSERT INTO fruits (name, calories, protein, cost, type) VALUES (?, ?, ?, ?, ?)`)
	}

	if err != nil {
		fmt.Printf("Error occurred in initial inserting data: %v", err)
	}

	_, errE := stmt.Exec(i.Name, i.Calories, i.Protein, i.Cost, i.Type)

	if errE != nil {
		fmt.Printf("Error occurred in inserting data: %v", errE)
	}

	defer stmt.Close()
}

func GetIngredientById(id int, index int, db *sql.DB) definitions.IngredientDetails {
	ing := definitions.IngredientDetails{}

	if index == int(definitions.CARBS) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM carbs WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "carbs"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.CARBS, errE)
		}

	}

	if index == int(definitions.PROTEINS) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM proteins WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "protein"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.PROTEINS, errE)
		}
	}

	if index == int(definitions.OILS) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM oils WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "oils"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.BEVERAGES, errE)
		}
	}

	if index == int(definitions.VEGETABLES) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM vegetables WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "vegetables"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.VEGETABLES, errE)
		}

	}

	if index == int(definitions.BEVERAGES) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM beverages WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "beverages"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.FRUITS, errE)
		}

	}

	if index == int(definitions.FRUITS) {
		errE := db.QueryRow(`SELECT name, calories, protein, cost FROM fruits WHERE id = ?`, id).Scan(&ing.Name, &ing.Calories, &ing.Protein, &ing.Cost)

		ing.Type = "fruits"

		if errE != nil {
			fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.OILS, errE)
		}
	}

	return ing
}

func GetObjectiveConstraintValues(db *sql.DB, c definitions.Chromosome, constraints definitions.AimObjectiveMap) definitions.Objectives {
	objectiveConstraintMap := make(map[string]float64)
	calories := make([]float64, 0)
	proteins := make([]float64, 0)
	prices := make([]float64, 0)
	var calorie float64
	var price float64
	var protein float64
	sumCalories := 0.0
	sumCost := 0.0
	sumProteins := 0.0

	for i, id := range c {
		if i == int(definitions.CARBS) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM carbs WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {
				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.CARBS, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}

		if i == int(definitions.PROTEINS) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM proteins WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {

				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.PROTEINS, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}

		if i == int(definitions.OILS) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM oils WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {
				fmt.Println(id)
				fmt.Println(GetMaxId(db, definitions.OILS))
				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.BEVERAGES, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}

		if i == int(definitions.VEGETABLES) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM vegetables WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {
				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.VEGETABLES, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}

		if i == int(definitions.BEVERAGES) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM beverages WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {
				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.FRUITS, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}

		if i == int(definitions.FRUITS) {
			errE := db.QueryRow(`SELECT calories, cost, protein FROM fruits WHERE id = ?`, id).Scan(&calorie, &price, &protein)

			if errE != nil {
				fmt.Printf("Error in executing query from %b, with message: %s \n", definitions.OILS, errE)
			}

			calories = append(calories, calorie)
			prices = append(prices, price)
			proteins = append(proteins, protein)
		}
	}

	for i := 0; i < len(c); i++ {
		sumCalories = sumCalories + calories[i]
		sumCost = sumCost + prices[i]
		sumProteins = sumProteins + proteins[i]
	}

	for _, v := range constraints {
		if v == definitions.CALORIES {
			objectiveConstraintMap[definitions.CALORIES] = sumCalories
		}

		if v == definitions.PROTEIN {
			objectiveConstraintMap[definitions.PROTEIN] = sumProteins
		}

		if v == definitions.PRICE {
			objectiveConstraintMap[definitions.PRICE] = sumCost
		}
	}

	return objectiveConstraintMap
}

func ShiftToDb(i definitions.Ingredients, db *sql.DB) {
	for _, v := range i {
		for _, m := range v {
			fmt.Println(m.Type)
			if m.Type == "carbs" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}

			if m.Type == "protein" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}

			if m.Type == "oil" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}

			if m.Type == "vegetables" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}

			if m.Type == "beverage" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}

			if m.Type == "fruit" {
				CreateTables(db, m.Type)
				InsertIngredient(db, m)
			}
		}
	}

	fmt.Println("DONE WRITING TO DB")
}

func ReadFromCsvFile() definitions.Ingredients {
	mealOpts := make(definitions.Ingredients, 0)
	f, err := os.Open("./stigler.csv")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, errR := r.ReadAll()

	if errR != nil {
		fmt.Println(errR)
	}

	if errR == io.EOF {
		fmt.Println("DONE!!!!")
	}

	for k, v := range records {
		ings := make(definitions.Ingredient)
		ingD := new(definitions.IngredientDetails)
		for i, m := range v {
			if k != 0 {
				if i == 2 {
					ingD.Cost = ParseStringToFloat(m)
				}

				if i == 3 {
					ingD.Calories = ParseStringToFloat(m)
				}

				if i == 4 {
					ingD.Protein = ParseStringToFloat(m)
				}

				if i == 0 {
					ingD.Name = m
				}

				if i == 12 {
					ingD.Type = m
				}
			}
		}
		ings[ingD.Name] = ingD
		mealOpts = append(mealOpts, ings)
	}

	return mealOpts
}

func ParseStringToFloat(s string) float64 {
	float, err := strconv.ParseFloat(s, 64)

	if err != nil {
		log.Fatalf("Err %v", err)
	}

	return float
}

func ParseStringToInt(s string) int {
	num, err := strconv.Atoi(s)

	if err != nil {
		log.Fatalf("Err %v", err)
	}

	return num
}

func GetMaxId(db *sql.DB, t definitions.IngredientType) int {
	var maxId int
	var err error
	if t == definitions.CARBS {
		err = db.QueryRow(`SELECT MAX(Id) FROM carbs`).Scan(&maxId)
	}

	if t == definitions.PROTEINS {
		err = db.QueryRow(`SELECT MAX(Id) FROM proteins`).Scan(&maxId)
	}

	if t == definitions.OILS {
		err = db.QueryRow(`SELECT MAX(Id) FROM oils`).Scan(&maxId)
	}

	if t == definitions.VEGETABLES {
		err = db.QueryRow(`SELECT MAX(Id) FROM vegetables`).Scan(&maxId)
	}

	if t == definitions.BEVERAGES {
		err = db.QueryRow(`SELECT MAX(Id) FROM beverages`).Scan(&maxId)
	}

	if t == definitions.FRUITS {
		err = db.QueryRow(`SELECT MAX(Id) FROM fruits`).Scan(&maxId)
	}

	if err != nil {
		fmt.Printf("Error occurred! err value: %v", err)
	}

	return maxId
}
