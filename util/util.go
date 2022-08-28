package util

import (
	"math/rand"
	"time"
)

//Gets a random ingredient ID stored in the db
func GetRandomIngredientId(l int) (Id int) {
	seed := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(seed)
	id := rg.Intn(l + 1)

	//Checks that the random ID is not 0
	if id == 0 {
		id = 1
	}

	return id
}


//Generating random probabilties for crossover
func GetRandomProbabilty() (prob float64) {
	seed := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(seed)
	pb := rg.Float64()
	return pb
}