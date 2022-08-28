package paretoselection

// import (
// 	"database/sql"
// 	"fmt"
// 	"sort"

// 	"github.com/neostefan/diet-assistant/db"
// 	"github.com/neostefan/diet-assistant/definitions"
// )

// type objectiveRange struct {
// 	max float64
// 	min float64
// }

// type smapObjectiveRanges []objectiveRange

// type chrObjective struct {
// 	chromosome definitions.Chromosome
// 	objectiveValue float64
// }

// func Pareto(sqlDb *sql.DB, g1, g2 definitions.Generation, objs definitions.AimObjectiveMap) {

// 	g1 = append(g1, g2...)
// 	generations := []definitions.Generation{
// 		g1,
// 		g2,
// 	}

// 	selectionMap := createSelectionMap(sqlDb, generations)
// 	mapLength := len(selectionMap)

// 	//see the result of this
// 	fmt.Printf("\n selectionMap of parents and children: %v with total: %v\n", selectionMap, mapLength)
// 	omap := getOldSelectionMap(selectionMap)
// 	//fmt.Printf("\n max and min of the objectives: %v \n", getSelectionMapObjectivesRange(*omap, ))
// 	pF := selection(&selectionMap)
// 	fmt.Printf("\n ParetoFront: %d \n", pF)
// 	fmt.Printf("\n New Generation: %d \n", newGeneration(pF, sqlDb, objs, *omap))

// }

// func createSelectionMap(sqlDb *sql.DB, gens []definitions.Generation) definitions.SelectionMap {
// 	selectionMap := make(definitions.SelectionMap)

// 	//try to to fix nested for loop error
// 	for _, gen := range gens {
// 		for i, c := range gen {
// 			objectives := db.GetCaloriesAndCost(sqlDb, c)
// 			selectionMap[c] = objectives
// 			fmt.Println(i, ": ", c)
// 		}
// 	}

// 	return selectionMap
// }

// func selection(smap *definitions.SelectionMap) definitions.ParetoFront {
// 	paretoFront := make(definitions.ParetoFront)
// 	i := 1
// 	counter := 0

// 	for {
// 		for chromosome, objectives := range *smap {
// 			var selectedChromosome definitions.Chromosome
// 			var selectedObjectives definitions.Objectives
// 			counter = counter + 1
// 			pFrontLength := len(paretoFront[i])
// 			cIndex := pFrontLength - 1

// 			if cIndex < 0 {
// 				paretoFront[i] = append(paretoFront[i], chromosome)
// 			} else {
// 				sMap := *smap

// 				selectedChromosome = chromosome
// 				selectedObjectives = objectives
// 				pFrontObjectives := sMap[paretoFront[i][cIndex]]
// 				fmt.Printf("selected chromosome: %v \n", chromosome)

// 				//checks if the GA objectives are fulfilled
// 				cond1, cond2 := checkConditions(selectedObjectives, pFrontObjectives)
// 				fmt.Printf("Objectives being compared, selected: %v, next one: %v \n", pFrontObjectives, selectedObjectives)

// 				//if both are satisfied add it to a paretoFront
// 				// if cond1 && cond2 {
// 				// 	paretoFront[i] = append(paretoFront[i], chromosome)
// 				// }

// 				//if one is not satisfied or dominated add the chromosome to the same front
// 				if cond1 && !cond2 || !cond1 && cond2 {
// 					if pFrontLength > 0 {
// 						checkParetoFrontElements(&paretoFront, selectedChromosome, i, chromosome, sMap)
// 					} else {
// 						paretoFront[i] = append(paretoFront[i], chromosome)
// 					}
// 				}

// 				if !cond1 && !cond2 {
// 					if pFrontLength > 0 {
// 						checkParetoFrontElements(&paretoFront, selectedChromosome, i, chromosome, sMap)
// 					} else {
// 						paretoFront[i][0] = chromosome
// 						//paretoFront = updateParetoFront(paretoFront, i, selectedChromosome, index)
// 					}
// 				}
// 			}
// 			fmt.Printf("Pareto Front: %v at iteration number: %v \n", paretoFront, counter)
// 		}

// 		counter = 0

// 		for _, v := range paretoFront[i] {
// 			delete(*smap, v)
// 		}

// 		i = i + 1

// 		if len(*smap) == 0 {
// 			break
// 		}
// 	}

// 	return paretoFront
// }

// //calculates the crowding distance and returns the more lonely ones or isolated ones per objective key
// func findCrowdingDistance(oSmap, pFsMap definitions.SelectionMap, keys definitions.AimObjectiveMap) []chrObjective {
// 	crDistValues := make([]float64, len(keys)) //array of the crowding distance indexed according to the sorted front
// 	crDistMap := make(map[definitions.Chromosome][]float64)
// 	crv := 0
// 	finalSolutionList := make([]chrObjective, 0)
// 	for aim, key := range keys {
// 		sPfValues := sortParetoFrontSelectionMap(key, pFsMap) // sorted values in the front
// 		generationRange := getSelectionMapObjectivesRange(oSmap, key) //max and min of the original selection map
// 		lastValueInSort := len(sPfValues) - 1  //last value in the sorted front(biggest value)

// 		if aim == definitions.MAX {

// 			fmt.Println("The selected paretoFRont and it's objective values: ", sPfValues)
// 			for i := range sPfValues {
// 				//checks if the index is between the last value and the first then calculates the crowding distance
// 				if i > 0 && i < lastValueInSort {
// 					numerator := sPfValues[i + 1].objectiveValue - sPfValues[i - 1].objectiveValue
// 					denominator := generationRange.max - generationRange.min
// 					fmt.Printf("\n Numerator: %v \n", numerator)
// 					crDistValues[crv] = numerator/denominator
// 					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], numerator/denominator)
// 				}

// 				//assigns the index of the biggest or maximum value for the key a large crowding distance
// 				if i == lastValueInSort {
// 					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], 1000.00)
// 				}
// 			}
// 		}

// 		if aim == definitions.MIN {
// 			fmt.Println("The selected paretoFRont and it's objective values: ", sPfValues)
// 			for i := range sPfValues {
// 				//checks if the index is between the last value and the first then calculates the crowding distance
// 				if i > 0 && i < lastValueInSort {
// 					numerator := sPfValues[i + 1].objectiveValue - sPfValues[i - 1].objectiveValue
// 					denominator := generationRange.max - generationRange.min
// 					crDistValues[crv] = numerator/denominator
// 					fmt.Printf("\n Numerator: %v \n", numerator)
// 					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], numerator/denominator)
// 				}

// 				//assigns the index of the biggest or maximum value for the key a large crowding distance
// 				if i == 0 {
// 					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], 1000.00)
// 				}
// 			}
// 		}

// 		crv = crv + 1
// 	}

// 	for k, v := range crDistMap {
// 		fmt.Println("The values cr distance: ", v)
// 		chr := chrObjective{
// 			chromosome: k,
// 			objectiveValue: addCrDistance(v),
// 		}

// 		finalSolutionList = append(finalSolutionList, chr)
// 	}

// 	return finalSolutionList
// }

// func addCrDistance(crDistances []float64) float64 {
// 	curr := 0.0

// 	for _, v := range crDistances {
// 		curr = curr + v
// 	}

// 	return curr
// }

// func sortParetoFrontSelectionMap(key string, pFsMap definitions.SelectionMap) []chrObjective {
// 	values := make([]float64, 0)
// 	chrObjectiveMap := make([]chrObjective, 0)

// 	for _, objectives := range pFsMap {
// 		for k, objValue := range objectives {
// 			if k == key {
// 				values = append(values, objValue)
// 			}
// 		}
// 	}

// 	sort.Float64s(values)

// 	for _, v := range values {
// 		c := findChromosomeBasedOnObjectiveValue(v, pFsMap)
// 		chrObj := chrObjective{
// 			chromosome: c,
// 			objectiveValue: v,
// 		}
// 		chrObjectiveMap = append(chrObjectiveMap, chrObj)
// 	}

// 	return chrObjectiveMap
// }

// //find which chromosome owns a particular value of calorie or price
// func findChromosomeBasedOnObjectiveValue(obj float64, oSmap definitions.SelectionMap) definitions.Chromosome {
// 	var chromosome definitions.Chromosome

// 	for c, objectives := range oSmap {
// 		for _, value := range objectives {
// 			if value == obj {
// 				chromosome = c
// 			}
// 		}
// 	}

// 	return chromosome
// }

// func getHighestCrowdingDistance(crMap *[]chrObjective) int {
// 	highestIndex := 0
// 	cRMap := *crMap

// 	for i, v := range *crMap {
// 		if i > 0 {
// 			if v.objectiveValue > cRMap[i - 1].objectiveValue {
// 				highestIndex = i
// 			}
// 		}
// 	}

// 	return highestIndex
// }

// //loops through the map for the highest crowding distance and returns the chromosome
// func pickChromosomeByHighestCrowdingDistance(crMap []chrObjective, slots int) []definitions.Chromosome {

// 	var bChromosomes []definitions.Chromosome
// 	counter := 0

// 	for slots > counter {
// 		index := getHighestCrowdingDistance(&crMap)
// 		fmt.Println("The Index selected: ", index)
// 		fmt.Println("The Chromosome selected: ", crMap[index].chromosome)
// 		fmt.Println("The Counter: ", counter)
// 		bChromosomes = append(bChromosomes, crMap[index].chromosome)
// 		crMap[index] = crMap[len(crMap) - 1]
// 		crMap = crMap[:len(crMap) - 1]
// 		counter++
// 	}

// 	return bChromosomes
// }

// //returns the next generation after the pareto selection is done
// func newGeneration(pF definitions.ParetoFront, sqlDb *sql.DB, keys definitions.AimObjectiveMap, osMap definitions.SelectionMap) definitions.Generation {
// 	var slotsAvailable int
// 	nG := make(definitions.Generation, 0)
// 	pFkeys := sortParetoFrontKeys(pF)
// 	slotsAvailable = definitions.PopulationSize - len(nG)
// 	fmt.Println("Slots Available: ", slotsAvailable)
// 	fmt.Println("length of nG: ", len(nG))
// 	fmt.Println("Pareto Front: ", pF)

// 	for _, key := range pFkeys {
// 		if len(nG) < definitions.PopulationSize {
// 			if len(pF[key]) > slotsAvailable {
// 				fmt.Println("HERE!!!!!!!!")
// 				pFgeneration := []definitions.Generation{
// 					pF[key],
// 				}
// 				pFsMap := createSelectionMap(sqlDb, pFgeneration)
// 				final := findCrowdingDistance(osMap, pFsMap, keys)
// 				chromosomes := pickChromosomeByHighestCrowdingDistance(final, slotsAvailable)
// 				for _, v := range chromosomes {
// 					nG = append(nG, v)
// 				}

// 			} else {
// 				for _, chromosome := range pF[key] {
// 					nG = append(nG, chromosome)
// 					slotsAvailable = definitions.PopulationSize - len(nG)
// 				}
// 			}
// 		}
// 	}

// 	return nG
// }

// func sortParetoFrontKeys(pF definitions.ParetoFront) []int {
// 	pFkeys := make([]int, 0)
// 	for k := range pF {
// 		pFkeys = append(pFkeys, k)
// 	}

// 	sort.Ints(pFkeys)

// 	return pFkeys
// }

// func getSelectionMapObjectivesRange(smap definitions.SelectionMap, key string) objectiveRange {
// 	objectiveOne := make([]float64, 0)
// 	//objectiveTwo := make([]float64, 0)
// 	//objRanges := make(smapObjectiveRanges, 0)
// 	objOneRange := objectiveRange{}
// 	//objTwoRange := objectiveRange{}
// 	maxObjOne := 0.0
// 	//maxObjTwo := 0.0
// 	minObjOne := 0.0
// 	//minObjTwo := 0.0

// 	for _, objectives := range smap {
// 		for k, v := range objectives {
// 			//change the key to objective one
// 			if k == key {
// 				objectiveOne = append(objectiveOne, v)
// 			}

// 			//change the key to objective two
// 			// if k == "prices" {
// 			// 	objectiveTwo = append(objectiveTwo, v)
// 			// }
// 		}
// 	}

// 	for i := range objectiveOne {
// 		if i == 0 {
// 			maxObjOne = objectiveOne[i]
// 			minObjOne = objectiveOne[i]
// 			// maxObjTwo = objectiveTwo[i]
// 			// minObjTwo = objectiveTwo[i]
// 		} else {
// 			if maxObjOne < objectiveOne[i] {
// 				maxObjOne = objectiveOne[i]
// 			}

// 			if minObjOne > objectiveOne[i] {
// 				minObjOne = objectiveOne[i]
// 			}

// 			// if maxObjTwo < objectiveTwo[i] {
// 			// 	maxObjTwo = objectiveTwo[i]
// 			// }

// 			// if minObjTwo > objectiveTwo[i] {
// 			// 	minObjTwo = objectiveTwo[i]
// 			// }
// 		}

// 	}

// 	objOneRange.max = maxObjOne
// 	objOneRange.min = minObjOne
// 	// objTwoRange.max = maxObjTwo
// 	// objTwoRange.min = minObjTwo

// 	// objRanges = append(objRanges, objOneRange)
// 	// objRanges = append(objRanges, objTwoRange)

// 	return objOneRange
// }

// func getOldSelectionMap(smap definitions.SelectionMap) *definitions.SelectionMap {
// 	oldSelectionMap := make(definitions.SelectionMap)

// 	for k, v := range smap {
// 		oldSelectionMap[k] = v
// 	}

// 	return &oldSelectionMap
// }

// //checks if the objectives of the chromosome is better than the selected chromosome objectives
// func checkConditions(objs definitions.Objectives, sObjs definitions.Objectives) (bool, bool) {
// 	cond1 := true //by default our selected chromosome is the best
// 	cond2 := true

// 	for key, value := range objs {
// 		if key == "calories" {
// 			//maximizing calories
// 			if value > sObjs[key] {
// 				cond1 = false
// 			}
// 		}

// 		if key == "prices" {
// 			//minimizing price
// 			if value < sObjs[key] {
// 				cond2 = false
// 			}
// 		}

// 	}

// 	return cond1, cond2
// }

// //get index of a chromosome(sC) in the pareto front(sG)
// func checkForChromosome(sG definitions.Generation, sC definitions.Chromosome) (index int, ok bool) {

// 	ok = false

// 	for i, v := range sG {
// 		if v == sC {
// 			index = i
// 			ok = true
// 		}
// 	}

// 	return index, ok
// }

// //! check this possible likely hood i am removing duplicates
// // returns index of a chromosome from the pareto front for deletion
// func updateParetoFront(pF *definitions.ParetoFront, index int, nC definitions.Chromosome, frontIndex int) {

// 	Pf := *pF

// 	fmt.Println("The frontIndex: ", frontIndex)
// 	fmt.Printf("\n Pareto front i am working with: %v \n", pF)
// 	fmt.Printf("\n Pareto front in the update: %v \n", Pf)

// 	_, ok := checkForChromosome(Pf[frontIndex], nC)
// 	//oldIndex, ok := checkForChromosome(Pf[index], Pf[index][])

// 	fmt.Println("Item is in the pareto front: ", ok)

// 	if !ok && len(Pf) > 0 {
// 		Pf[frontIndex][index] = nC
// 	}

// 	if !ok && len(Pf) == 0 {
// 		Pf[frontIndex][0] = nC
// 	}

// 	// pF = &Pf
// }

// func deleteFromParetoFront(pF *definitions.ParetoFront, specidiedChromosome definitions.Chromosome, frontIndex int, cIndex int) {
// 	Pf := *pF
// 	pG := Pf[frontIndex]

// 	if len(pG) == 0 {
// 		pG = make(definitions.Generation, 0)
// 		Pf[frontIndex] = pG
// 	} else {
// 		pG[cIndex] = pG[len(pG) - 1]
// 		Pf[frontIndex] = pG[:len(pG) - 1]
// 	}

// 	// pF = &Pf
// }

// func addToParetoFront(pF *definitions.ParetoFront, nC definitions.Chromosome, index int) {

// 	Pf := *pF
// 	pG := Pf[index]
// 	_, ok := checkForChromosome(pG, nC)

// 	if !ok {
// 		pG = append(pG, nC)
// 	}

// 	Pf[index] = pG
// }

// func checkParetoFrontElements(pF *definitions.ParetoFront, sC definitions.Chromosome, index int, nC definitions.Chromosome, smap definitions.SelectionMap) {
// 	paretoFrontCheckMap := make(map[int]string)
// 	Pf := *pF
// 	for i, c := range Pf[index] {
// 		fmt.Printf("\n objectives being compared chromosome: %v selected chromosome: %v \n", smap[c], smap[sC])
// 		cond1, cond2 := checkConditions(smap[c], smap[sC])

// 		fmt.Println("Codition check result: ", cond1, cond2)

// 		if cond1 && cond2 {
// 			fmt.Printf("\n Chromosome: %d index: %v \n", c, i)
// 			paretoFrontCheckMap[i] = "true"
// 			updateParetoFront(pF, i, nC, index)
// 		}

// 		if !cond1 && cond2 || cond1 && !cond2 {
// 			fmt.Printf("\n Chromosome: %d index: %v \n", c, i)
// 			paretoFrontCheckMap[i] = "non-dominating"
// 			addToParetoFront(pF, nC, index)
// 		}

// 		if !cond1 && !cond2 {
// 			fmt.Printf("\n Chromosome: %d index: %v \n", c, i)
// 			paretoFrontCheckMap[i] = "false"
// 			deleteFromParetoFront(pF, nC, i, index)
// 		}

// 		// for k, v := range paretoFrontCheckMap {

// 		// 	//! current problem
// 		// 	//! try to see how i can adjust to multiple checks of true, non-dominating for one objective
// 		// 	if v == "false" {
// 		// 		fmt.Println("I am here!!!!")

// 		// 	}

// 		// 	if v == "non-dominating" {

// 		// 	}
// 		// }
// 	}

// 	fmt.Println("Pareto Front currently: ", Pf)
// }