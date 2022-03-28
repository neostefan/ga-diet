package definitions

const (
	AlleleSize = 2
	ChromosomeSize = 6
	GenerationSize = 200
)

type MealDetails struct {
	Name string
	Calories float64
	Cost float64 
}

type Allele [AlleleSize]int

type Chromosome [ChromosomeSize]Allele

type Generation map[MealDetails]Chromosome

