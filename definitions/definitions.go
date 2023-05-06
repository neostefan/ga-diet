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

//FG Classification
//Beverage

type diabetesConstraint struct {
	FGClassification string // food group classification
}

type ulcerConstraint struct {
	Flavor    string
	Taste     string
	Cullinary string
}

//Age constraints definitions

const (
	CALORIES = "calories"
	PROTEIN  = "protein"
	PRICE    = "cost"
)

type DietCondition int

const (
	DIABETES DietCondition = iota
	ULCER
)

const (
	MAX Aim = iota
	MIN
)

const (
	ChromosomeSize = 5
	PopulationSize = 30
	GenerationSize = 50
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
	Name      string
	Calories  float64
	Carbs     float64
	Fat       float64
	Protein   float64
	Cost      float64
	FoodGroup string
	Cullinary string
	Allergen  string
	Taste     string
	Flavor    string
	Type      string
}

type MealOptions []*IngredientDetails

// Chromosome contains an array of five integers
type Chromosome [ChromosomeSize]int

type Generation []Chromosome

type Ingredient map[string]*IngredientDetails

type Ingredients []Ingredient

// initialize an ulcer constraint
func UlcerConstraint() ulcerConstraint {
	return ulcerConstraint{
		Taste:     "spicy", //although taste and flavor are same for situations that require multiple limits they are seperated
		Flavor:    "acidic",
		Cullinary: "fried",
	}
}

// initialize diabetic constraint
func DiabetesConstraint() diabetesConstraint {
	return diabetesConstraint{
		FGClassification: "whole",
	}
}

// swap function swap the chromosome with it's partner based on the index specified
func (c *Chromosome) Swap(partner *Chromosome, index int) {
	c[index], partner[index] = partner[index], c[index]
}
