package crossover

import (
	"github.com/neostefan/diet-assistant/definitions"
	"github.com/neostefan/diet-assistant/util"
)

//cross function performs genetic operations; crossover and mutation
func cross(c1, c2 *definitions.Chromosome, index int) {

	endPoint := len(c1)

	for crossPoint := index + 1; crossPoint < endPoint; crossPoint++ {
		c1.Swap(c2, crossPoint)
	}

}

//TODO generate probabilities and use them to determine genetic operations for both crossover and mutation
/* function makes use of the cross function to perform crossover takes previous generation
index to cross*/
func Crossover(g definitions.Generation, index int) definitions.Generation {
	i := 0
	gen := definitions.Generation{}
	lastC := len(g) - 1
	pi := 0

	probabilitySize := definitions.PopulationSize / 2
	if definitions.PopulationSize % 2 != 0 {
		probabilitySize = (definitions.PopulationSize - 1) / 2
	}
	
	probabilities := make([]float64, probabilitySize)

	for p := 0; p < probabilitySize; p++ {
		probabilities[p] = util.GetRandomProbabilty()
	}

	//! Old CrossOver Code the challenge was that population sizes of odd numbers were exceeding their ranges
	for len(gen) != len(g) {
		if i < lastC {
			c1 := &g[i]
			c2 := &g[i + 1]

			//0.25 is the min allowance for crossover to occur
			if probabilities[pi] > 0.1 {
				cross(c1, c2, index)
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