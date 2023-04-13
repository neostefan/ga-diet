package operators

import (
	"database/sql"
	"fmt"

	"github.com/neostefan/ga-diet/db"
	"github.com/neostefan/ga-diet/definitions"
)

// Initializes the population
func InitializePopulation(sqlDB *sql.DB, conditions []definitions.DietCondition) definitions.Generation {

	var cr definitions.Chromosome
	var generation definitions.Generation

	//Loops through the defined population size
	for i := 0; i < definitions.PopulationSize; i++ {

		//Loops through the defined Chromosome size
		for j := 0; j < definitions.ChromosomeSize; j++ {

			//Update the Chromosome with the ingredient meant for that index
			switch {
			case j == 0:
				maxId := db.GetMaxId(sqlDB, definitions.CARBS)
				fmt.Println("CARBS MAX ID: ", maxId)
				cr[j] = db.GetRandomId(sqlDB, definitions.CARBS, conditions, maxId)
			case j == 1:
				maxId := db.GetMaxId(sqlDB, definitions.PROTEINS)
				cr[j] = db.GetRandomId(sqlDB, definitions.PROTEINS, conditions, maxId)
				// cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.PROTEINS))
			case j == 2:
				maxId := db.GetMaxId(sqlDB, definitions.VEGETABLES)
				cr[j] = db.GetRandomId(sqlDB, definitions.VEGETABLES, conditions, maxId)
				// cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.VEGETABLES))
			case j == 3:
				maxId := db.GetMaxId(sqlDB, definitions.BEVERAGES)
				cr[j] = db.GetRandomId(sqlDB, definitions.BEVERAGES, conditions, maxId)
				// cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.BEVERAGES))
			case j == 4:
				maxId := db.GetMaxId(sqlDB, definitions.FRUITS)
				cr[j] = db.GetRandomId(sqlDB, definitions.FRUITS, conditions, maxId)
				// cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.FRUITS))
			default:
				return nil
			}
		}
		generation = append(generation, cr)
	}

	return generation
}
