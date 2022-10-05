package definitions

type IngredientType int32

type Objectives map[string]float64

// type SelectionMap map[Chromosome]Objectives
type SelectionMap map[int]Objectives

// type ParetoFront map[int]Generation

type ParetoFront map[int][]int

type Aim int32

type AimObjectiveMap map[Aim]string

type ChromosomeCache struct {
	Chromosome Chromosome
	Objectives Objectives
}

const (
	CALORIES = "calories"
	PROTEIN  = "protein"
	PRICE    = "cost"
)

const (
	MAX Aim = iota
	MIN
)

const (
	ChromosomeSize = 6
	PopulationSize = 3
	GenerationSize = 200
)

const (
	CARBS IngredientType = iota
	PROTEINS
	OILS
	VEGETABLES
	BEVERAGES
	FRUITS
)

// const (
// 	Mutation OperationType = iota
// 	Crossover
// )

type HealthyPlateContent int
type IngredientDetails struct {
	Name     string
	Calories float64
	Protein  float64
	Cost     float64
	Type     string
}

type MealOptions []*IngredientDetails

//Chromosome contains an array of six integers
type Chromosome [ChromosomeSize]int

type Generation []Chromosome

type Ingredient map[string]*IngredientDetails

type Ingredients []Ingredient

//swap function swap the chromosome with it's partner based on the index specified
func (c *Chromosome) Swap(partner *Chromosome, index int) {
	c[index], partner[index] = partner[index], c[index]
}
