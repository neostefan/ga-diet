package nsga

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neostefan/diet-assistant/db"
	"github.com/neostefan/diet-assistant/definitions"
	"github.com/neostefan/diet-assistant/operators"
	"github.com/neostefan/diet-assistant/operators/crossover"
	"github.com/neostefan/diet-assistant/operators/mutation"
	paretoselection "github.com/neostefan/diet-assistant/operators/pareto_selection"
)

func RunAlgorithm(maxObj string, minObj string) []definitions.IngredientDetails {
	// ings := db.ReadFromCsvFile()

	sqlDb, err := sql.Open("sqlite3", "./db/meals.db")

	if err != nil {
		fmt.Printf("Error occurred: %v", err)
	}

	// db.ShiftToDb(ings, sqlDb)

	var finalIngs []definitions.IngredientDetails
	population := operators.InitializePopulation(sqlDb)
	parents := make(definitions.Generation, len(population))
	fmt.Println("Printing the initial states...")

	//where i would add the for loop
	for i := 0; i < 1; i++ {
		copy(parents, population)
		fmt.Println("Generation: ", i+1)
		fmt.Println("Parents: ")
		printChromosome(population)
		//c1.Evolve(&c2, definitions.Crossover, 1)
		population = crossover.Crossover(population, 2)
		fmt.Println("Printing the crossed states...")
		printChromosome(population)
		population = mutation.Mutation(population, 2)
		fmt.Println("Printing the mutated states...")
		printChromosome(population)
		fmt.Println("printing pareto")
		aimObj := make(definitions.AimObjectiveMap)
		aimObj[definitions.MAX] = maxObj
		aimObj[definitions.MIN] = minObj
		population = paretoselection.Pareto(sqlDb, parents, population, aimObj)
		fmt.Println("printing result of the generation")
		printChromosome(population)
		fmt.Println("printing the diet meal picked")
		selectedDiet := population[0]

		finalIngs = decodeChromosome(selectedDiet, sqlDb)
		fmt.Printf("\n The selected diet: %v", finalIngs)
	}

	defer sqlDb.Close()
	return finalIngs
}

func printChromosome(g definitions.Generation) {
	for _, c := range g {
		fmt.Printf("%d \n", c)
	}
}

func decodeChromosome(c definitions.Chromosome, sqlDB *sql.DB) []definitions.IngredientDetails {
	ings := make([]definitions.IngredientDetails, 0)

	for i, ingId := range c {

		ing := db.GetIngredientById(ingId, i, sqlDB)

		ings = append(ings, ing)
	}

	return ings
}
