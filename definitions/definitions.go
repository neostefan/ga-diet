package definitions

type IngredientType string
type OperationType int

const (
	ChromosomeSize = 6
	PopulationSize = 50
	GenerationSize = 200
)

const (
	CARBS IngredientType = "carbs"
	PROTEINS IngredientType = "proteins"
	VEGETABLES IngredientType = "vegetables"
	OILS IngredientType = "oils"
	FRUITS IngredientType = "fruits"
	BEVERAGES IngredientType = "beverages"
)

const (
	Mutation OperationType = iota
	Crossover 
)

type HealthyPlateContent int
type IngredientDetails struct {
	Name string
	Calories float64
	Cost float64
	Type string 
}

type MealOptions []*IngredientDetails

//Chromosome contains an array of six integers
type Chromosome [ChromosomeSize]int

type Generation []Chromosome

type Ingredient map[string]*IngredientDetails

type Ingredients []Ingredient

//swap function swap the chromosome with it's partner based on the index specified
func (c *Chromosome) swap(partner *Chromosome, index int) {
	c[index], partner[index] = partner[index], c[index]
}

//evolve function performs genetic operations; crossover and mutation
func (c *Chromosome) Evolve(partner *Chromosome, otype OperationType, index int) {
	if otype == Crossover {
		endPoint := len(c)

		for crossPoint := index + 1; crossPoint < endPoint; crossPoint++ {
			c.swap(partner, crossPoint)
		}

		return	
	}

	c.swap(partner, index)
}

