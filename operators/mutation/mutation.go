package mutation

import (
	"github.com/neostefan/ga-diet/definitions"
	"github.com/neostefan/ga-diet/util"
)

func mutate(c1, c2 *definitions.Chromosome, index int) {
	c1.Swap(c2, index)
}

func Mutation(g definitions.Generation, index int) definitions.Generation {
	i := 0
	gen := definitions.Generation{}
	lastC := len(g) - 1
	pi := 0

	//creates the probabilities array of each chromosome
	probabilitySize := definitions.PopulationSize / 2
	if definitions.PopulationSize%2 != 0 {
		probabilitySize = (definitions.PopulationSize - 1) / 2
	}

	probabilities := make([]float64, probabilitySize)

	for p := 0; p < probabilitySize; p++ {
		probabilities[p] = util.GetRandomProbabilty()
	}

	//! Previous code for mutation check crossover.go for issue explanation
	for len(g) != len(gen) {
		if i < lastC {
			c1 := &g[i]
			c2 := &g[i+1]

			//0.5 is the minimum allowance for mutation to occur
			if probabilities[pi] > 0.5 {
				mutate(c1, c2, index)
			}

			gen = append(gen, *c1, *c2)

			pi = pi + 1
			i = i + 2
		} else {
			gen = append(gen, g[lastC])
		}
	}

	return gen
}
