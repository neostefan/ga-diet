package operators

import (
	"database/sql"

	"github.com/neostefan/diet-assistant/db"
	"github.com/neostefan/diet-assistant/definitions"
	"github.com/neostefan/diet-assistant/util"
)

//Initializes the population
func InitializePopulation(sqlDB *sql.DB) definitions.Generation {

	var cr definitions.Chromosome
	var generation definitions.Generation

	//Loops through the defined population size
	for i := 0; i < definitions.PopulationSize; i++ {

		//Loops through the defined Chromosome size
		for j := 0; j < definitions.ChromosomeSize; j++ {

			//Update the Chromosome with the ingredient meant for that index
			switch {
				case j == 0:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.CARBS))
				case j == 1:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.PROTEINS))
				case j == 2:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.OILS))
				case j == 3:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.VEGETABLES))
				case j == 4:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.BEVERAGES))
				case j == 5:
					cr[j] = util.GetRandomIngredientId(db.GetMaxId(sqlDB, definitions.FRUITS))
				default:
					return nil
			}
		}
		generation = append(generation, cr)
	}

	return generation
}