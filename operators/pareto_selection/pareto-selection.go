package paretoselection

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/neostefan/diet-assistant/db"
	"github.com/neostefan/diet-assistant/definitions"
)

type objectiveRange struct {
	max float64
	min float64
}

//contains results of the check between pF elements and the next chromosome in the selection map
type pFcheckMapResult struct {
	result string
	pFChromosomeIndex int
	currChromosomeIndexNotInFront int
}

// type smapObjectiveRanges []objectiveRange

type chrObjective struct {
	chromosome int
	objectiveValue float64
}

func Pareto(sqlDb *sql.DB, g1, g2 definitions.Generation, objs definitions.AimObjectiveMap) definitions.Generation {

	selectionGen := make(definitions.Generation, 0)
	selectionGen = append(selectionGen, g2...)
	selectionGen = append(selectionGen, g1...)

	//! SAMPLE CHROMOSOME FOR TESTING
	// chrm := definitions.Chromosome{
	// 	6, 17, 2, 5, 4 ,8,
	// }

	// chrm02 := definitions.Chromosome{
	// 	11, 24, 3, 5, 1 ,8,
	// }

	// chrm03 := definitions.Chromosome{
	// 	2, 18, 1, 3, 3 ,8,
	// }

	// chrm04 := definitions.Chromosome{
	// 	6, 17, 3, 5, 4 ,8,
	// }

	// chrm05 := definitions.Chromosome{
	// 	11, 24, 2, 5, 1, 8,
	// }

	// chrm06 := definitions.Chromosome{
	// 	2, 18, 1, 3, 3, 8,
	// }

	// selectionGen = append(selectionGen, chrm)
	// selectionGen = append(selectionGen, chrm02)
	// selectionGen = append(selectionGen, chrm03)
	// selectionGen = append(selectionGen, chrm04)
	// selectionGen = append(selectionGen, chrm05)
	// selectionGen = append(selectionGen, chrm06)

	selectionMap := createSelectionMap(sqlDb, selectionGen)
	// mapLength := len(selectionMap)

	//see the result of this
	// fmt.Printf("\n selectionMap of parents and children: %v with total: %v\n", selectionMap, mapLength)
	omap := getOldSelectionMap(selectionMap)
	//fmt.Printf("\n max and min of the objectives: %v \n", getSelectionMapObjectivesRange(*omap, ))
	pF := selection(&selectionMap, selectionGen)
	fmt.Printf("\n ParetoFront: %d \n", pF)
	// fmt.Printf("\n New Generation: %d \n", newGeneration(pF, sqlDb, objs, *omap, selectionGen))
	
	return newGeneration(pF, sqlDb, objs, *omap, selectionGen)
}

//creates the selectionmap from the slice of parents and children
func createSelectionMap(sqlDb *sql.DB, selectionGen definitions.Generation) definitions.SelectionMap {
	selectionMap := make(definitions.SelectionMap)

	//try to fix nest
	for i, c := range selectionGen {
		objectives := db.GetCaloriesAndCost(sqlDb, c)
		selectionMap[i] = objectives
		fmt.Println(i, ": ", c)
	}

	return selectionMap
}

func selection(smap *definitions.SelectionMap, selectionGen definitions.Generation) definitions.ParetoFront {
	paretoFront := make(definitions.ParetoFront)
	i := 1
	counter := 0

	for {
		for chromosomeIndex, objectives := range *smap {
			var selectedChromosome int //the chromosome inside of the smap we are on currently
			var selectedObjectives definitions.Objectives
			counter = counter + 1
			pFrontLength := len(paretoFront[i])
			cIndex := pFrontLength - 1

			if cIndex < 0 {
				paretoFront[i] = append(paretoFront[i], chromosomeIndex)
			} else {
				sMap := *smap
			
				selectedChromosome = chromosomeIndex
				selectedObjectives = objectives
				pFrontObjectives := sMap[paretoFront[i][cIndex]]
				fmt.Printf("selected chromosome: %v \n", selectionGen[chromosomeIndex])
				


				//checks if the GA objectives are fulfilled
				cond1, cond2 := checkConditions(selectedObjectives, pFrontObjectives)
				fmt.Printf("Objectives being compared, selected: %v, next one: %v \n", pFrontObjectives, selectedObjectives)


				//if both are satisfied add it to a paretoFront
				// if cond1 && cond2 {
				// 	paretoFront[i] = append(paretoFront[i], chromosome)
				// }

				//if one is not satisfied or dominated add the chromosome to the same front
				if cond1 && !cond2 || !cond1 && cond2 {
					if pFrontLength > 0 {
						checkParetoFrontElements(&paretoFront, selectedChromosome, i, chromosomeIndex, sMap)
					} else {
						paretoFront[i] = append(paretoFront[i], chromosomeIndex)
					}
				}

				if !cond1 && !cond2 {
					if pFrontLength > 0 {
						checkParetoFrontElements(&paretoFront, selectedChromosome, i, chromosomeIndex, sMap)
					} else {
						paretoFront[i][0] = chromosomeIndex
						//paretoFront = updateParetoFront(paretoFront, i, selectedChromosome, index)
					}
				}
			}
			fmt.Printf("Pareto Front: %v at iteration number: %v \n", paretoFront, counter)
		}

		counter = 0

		for _, v := range paretoFront[i] {
			delete(*smap, v)
		}
		
		i = i + 1

		if len(*smap) == 0 {
			break
		}
	}

	return paretoFront
}

//calculates the crowding distance and returns the more lonely ones or isolated ones per objective key
func findCrowdingDistance(oSmap, pFsMap definitions.SelectionMap, keys definitions.AimObjectiveMap) []chrObjective {
	crDistValues := make([]float64, len(keys)) //array of the crowding distance indexed according to the sorted front
	crDistMap := make(map[int][]float64) //chromosomeIndex map to the crowding distance
	crv := 0
	finalSolutionList := make([]chrObjective, 0)
	for aim, key := range keys {
		sPfValues := sortParetoFrontSelectionMap(key, pFsMap) // sorted values in the front
		generationRange := getSelectionMapObjectivesRange(oSmap, key) //max and min of the original selection map
		lastValueInSort := len(sPfValues) - 1  //last value in the sorted front(biggest value)	

		if aim == definitions.MAX {
				
			fmt.Println("The selected paretoFront and it's objective values: ", sPfValues)
			for i := range sPfValues {
				//checks if the index is between the last value and the first then calculates the crowding distance
				if i > 0 && i < lastValueInSort {
					numerator := sPfValues[i + 1].objectiveValue - sPfValues[i - 1].objectiveValue
					denominator := generationRange.max - generationRange.min
					fmt.Printf("\n Numerator: %v \n", numerator)
					crDistValues[crv] = numerator/denominator
					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], numerator/denominator)
				}
	
				//assigns the index of the biggest or maximum value for the key a large crowding distance
				if i == lastValueInSort {
					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], 1000.00)
				}
			}
		}
	
		if aim == definitions.MIN {
			fmt.Println("The selected paretoFRont and it's objective values: ", sPfValues)
			for i := range sPfValues {
				//checks if the index is between the last value and the first then calculates the crowding distance
				if i > 0 && i < lastValueInSort {
					numerator := sPfValues[i + 1].objectiveValue - sPfValues[i - 1].objectiveValue
					denominator := generationRange.max - generationRange.min
					crDistValues[crv] = numerator/denominator
					fmt.Printf("\n Numerator: %v \n", numerator)
					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], numerator/denominator)
				}
	
				//assigns the index of the biggest or maximum value for the key a large crowding distance
				if i == 0 {
					crDistMap[sPfValues[i].chromosome] = append(crDistMap[sPfValues[i].chromosome], 1000.00)
				}
			}
		}

		crv = crv + 1
	}

	for k, v := range crDistMap {
		fmt.Println("The values cr distance: ", v)
		chr := chrObjective{
			chromosome: k,
			objectiveValue: addCrDistance(v),
		}

		finalSolutionList = append(finalSolutionList, chr) 
	}

	return finalSolutionList
}

func addCrDistance(crDistances []float64) float64 {
	curr := 0.0
	
	for _, v := range crDistances {
		curr = curr + v
	}
	
	return curr
}

func sortParetoFrontSelectionMap(key string, pFsMap definitions.SelectionMap) []chrObjective {
	values := make([]float64, 0)
	chrObjectiveMap := make([]chrObjective, 0)


	for _, objectives := range pFsMap {
		for k, objValue := range objectives {
			if k == key {
				values = append(values, objValue)
			}
		} 
	}

	sort.Float64s(values)

	for _, v := range values {
		c := findChromosomeBasedOnObjectiveValue(v, pFsMap)
		chrObj := chrObjective{
			chromosome: c,
			objectiveValue: v,
		}
		chrObjectiveMap = append(chrObjectiveMap, chrObj)
	}

	return chrObjectiveMap
}

//find which chromosome owns a particular value of calorie or price
func findChromosomeBasedOnObjectiveValue(obj float64, oSmap definitions.SelectionMap) int {
	var chromosomeIndex int

	for c, objectives := range oSmap {
		for _, value := range objectives {
			if value == obj {
				chromosomeIndex = c
			}
		}
	}

	return chromosomeIndex
}

func getHighestCrowdingDistance(crMap *[]chrObjective) int {
	highestIndex := 0
	cRMap := *crMap

	for i, v := range *crMap {
		if i > 0 {
			if v.objectiveValue > cRMap[i - 1].objectiveValue {
				highestIndex = i
			}
		} 
	}

	return highestIndex
}

//loops through the map for the highest crowding distance and returns the chromosome
func pickChromosomeByHighestCrowdingDistance(crMap []chrObjective, slots int) []int {
	
	var bChromosomes []int
	counter := 0

	for slots > counter {
		index := getHighestCrowdingDistance(&crMap)
		fmt.Println("The Index selected: ", index)
		fmt.Println("The Chromosome selected: ", crMap[index].chromosome)
		fmt.Println("The Counter: ", counter)
		bChromosomes = append(bChromosomes, crMap[index].chromosome)
		crMap[index] = crMap[len(crMap) - 1]
		crMap = crMap[:len(crMap) - 1]
		counter++
	}

	return bChromosomes
}

//returns the next generation after the pareto selection is done
func newGeneration(pF definitions.ParetoFront, sqlDb *sql.DB, keys definitions.AimObjectiveMap, osMap definitions.SelectionMap, selectionGen definitions.Generation) definitions.Generation {
	var slotsAvailable int
	nG := make(definitions.Generation, 0)
	pFkeys := sortParetoFrontKeys(pF)
	slotsAvailable = definitions.PopulationSize - len(nG)
	fmt.Println("Slots Available: ", slotsAvailable)
	fmt.Println("length of nG: ", len(nG))
	fmt.Println("Pareto Front: ", pF)
	
	//looping through a pareto front set of chromosome indicies
	for _, key := range pFkeys {

		//check the length of the new generation if it is less than the population size
		if len(nG) < definitions.PopulationSize {

			//check if there are still more spaces in the new generation
			if len(pF[key]) > slotsAvailable {
				fmt.Println("HERE!!!!!!!!")

				//initialize a new generation
				pFgen := make(definitions.Generation, 0)

				//convert the pareto front chromosome indicies values to their respective chromosomes
				for _, chromosomeIndex := range pF[key] {
					pFgen = append(pFgen, selectionGen[chromosomeIndex])
				}

				//create a new selection map based off the pareto front values
				pFsMap := createSelectionMap(sqlDb, pFgen)

				//finds the crowding distance based off the new selection map
				final := findCrowdingDistance(osMap, pFsMap, keys)

				//picks out the chromosome indicies by the highest crowding distance
				chromosomeIndicies := pickChromosomeByHighestCrowdingDistance(final, slotsAvailable)
				
				//convert the chromosomeIndicies back to the chromosomes
				for _, v := range chromosomeIndicies {
					nG = append(nG, pFgen[v])
				}
				
			} else {
				for _, chromosomeIndex := range pF[key] {
					nG = append(nG, selectionGen[chromosomeIndex])
					slotsAvailable = definitions.PopulationSize - len(nG)	
				}
			}	
		}
	}

	return nG
}

func sortParetoFrontKeys(pF definitions.ParetoFront) []int {
	pFkeys := make([]int, 0)
	for k := range pF {
		pFkeys = append(pFkeys, k)
	}

	sort.Ints(pFkeys)

	return pFkeys
}

func getSelectionMapObjectivesRange(smap definitions.SelectionMap, key string) objectiveRange {
	objectiveOne := make([]float64, 0)
	//objectiveTwo := make([]float64, 0)
	//objRanges := make(smapObjectiveRanges, 0)
	objOneRange := objectiveRange{}
	//objTwoRange := objectiveRange{}
	maxObjOne := 0.0
	//maxObjTwo := 0.0
	minObjOne := 0.0
	//minObjTwo := 0.0

	for _, objectives := range smap {
		for k, v := range objectives {
			//change the key to objective one 
			if k == key {
				objectiveOne = append(objectiveOne, v)
			}

			//change the key to objective two
			// if k == "prices" {
			// 	objectiveTwo = append(objectiveTwo, v)
			// }
		}
	}

	for i := range objectiveOne {
		if i == 0 {
			maxObjOne = objectiveOne[i]
			minObjOne = objectiveOne[i]
			// maxObjTwo = objectiveTwo[i]
			// minObjTwo = objectiveTwo[i]
		} else {
			if maxObjOne < objectiveOne[i] {
				maxObjOne = objectiveOne[i]
			}

			if minObjOne > objectiveOne[i] {
				minObjOne = objectiveOne[i]
			}

			// if maxObjTwo < objectiveTwo[i] {
			// 	maxObjTwo = objectiveTwo[i]
			// }

			// if minObjTwo > objectiveTwo[i] {
			// 	minObjTwo = objectiveTwo[i]
			// }
		}
		
	}

	objOneRange.max = maxObjOne
	objOneRange.min = minObjOne
	// objTwoRange.max = maxObjTwo
	// objTwoRange.min = minObjTwo

	// objRanges = append(objRanges, objOneRange)
	// objRanges = append(objRanges, objTwoRange)

	return objOneRange
}

func getOldSelectionMap(smap definitions.SelectionMap) *definitions.SelectionMap {
	oldSelectionMap := make(definitions.SelectionMap)

	for k, v := range smap {
		oldSelectionMap[k] = v
	}

	return &oldSelectionMap
}

//checks if the objectives of the chromosome is better than the selected chromosome objectives
func checkConditions(objs definitions.Objectives, sObjs definitions.Objectives) (bool, bool) {
	cond1 := true //by default our selected chromosome is the best
	cond2 := true

	for key, value := range objs {
		if key == "calories" {
			//maximizing calories
			if value > sObjs[key] {
				cond1 = false
			}
		}

		if key == "prices" {
			//minimizing price
			if value < sObjs[key] {
				cond2 = false
			}
		}

	}

	return cond1, cond2
}

//get index of the chromosomeIndex in the pareto front(sG)
func checkForChromosome(frontSliceOfChromosomeIndex []int, chromosomeIndex int) (index int, ok bool) {

	ok = false

	for i, v := range frontSliceOfChromosomeIndex {
		if v == chromosomeIndex {
			index = i
			ok = true
		}
	}

	return index, ok
}

//! check this possible likely hood i am removing duplicates
// returns index of a chromosome from the pareto front for deletion 
func updateParetoFront(pF *definitions.ParetoFront, index int, frontIndex int, checkResult *pFcheckMapResult) {	
	
	Pf := *pF
	
	fmt.Println("The frontIndex: ", frontIndex)
	fmt.Printf("\n Pareto front i am working with: %v \n", pF)
	fmt.Printf("\n Pareto front in the update: %v \n", Pf)

	_, ok := checkForChromosome(Pf[frontIndex], checkResult.currChromosomeIndexNotInFront)
	//oldIndex, ok := checkForChromosome(Pf[index], Pf[index][])

	fmt.Println("Item is in the pareto front: ", ok)

	//if the chromosome index in the selection map is not in the pareto front and pareto front is not empty
	if !ok && len(Pf) > 0 {
		//delete the front chromosome that does not satisfy the equation
		deleteFromParetoFront(pF, frontIndex, index)

		//add the chromsome not in the pareto front that satisfies the conditions to the front
		addToParetoFront(pF, checkResult.currChromosomeIndexNotInFront, frontIndex)
		// Pf[frontIndex][index] = checkResult.currChromosomeIndexNotInFront
	}

	if !ok && len(Pf) == 0 {
		Pf[frontIndex][0] = checkResult.currChromosomeIndexNotInFront
	}

	// pF = &Pf
}

//deletes a chromosomeIndex stored in the pareto front by the paretofront index and index of the chromosome in the frontslice
func deleteFromParetoFront(pF *definitions.ParetoFront, frontIndex int, frontChromosomeIndex int) {
	Pf := *pF
	frontSliceOfChromosomeIndex := Pf[frontIndex]

	if len(frontSliceOfChromosomeIndex) == 0 {
		frontSliceOfChromosomeIndex = make([]int, 0)
		Pf[frontIndex] = frontSliceOfChromosomeIndex
	} else {
		frontSliceOfChromosomeIndex[frontChromosomeIndex] = frontSliceOfChromosomeIndex[len(frontSliceOfChromosomeIndex) - 1]
		Pf[frontIndex] = frontSliceOfChromosomeIndex[:len(frontSliceOfChromosomeIndex) - 1]
	}

	// pF = &Pf
}

func addToParetoFront(pF *definitions.ParetoFront, chromosomeIndex int, frontIndex int) {
	
	Pf := *pF
	frontSliceOfChromosomeIndex := Pf[frontIndex]
	_, ok := checkForChromosome(frontSliceOfChromosomeIndex, chromosomeIndex)

	if !ok {
		frontSliceOfChromosomeIndex = append(frontSliceOfChromosomeIndex, chromosomeIndex)	
	}
	
	Pf[frontIndex] = frontSliceOfChromosomeIndex
}

//pFChroIndex sC
func checkParetoFrontElements(pF *definitions.ParetoFront, pFChromosomeIndex int, index int, currChromosomeIndex int, smap definitions.SelectionMap) {
	Pf := *pF
	chromsomeCheckMap := map[string]map[int]int{}
	currChromosomeIndexTopFchromosomeIndex := map[int]int{}
	totalItemsInPf := len(Pf[index]) - 1

	for i, c := range Pf[index] {
		fmt.Printf("\n objectives being compared chromosome: %v selected chromosome: %v \n", smap[c], smap[pFChromosomeIndex])
		cond1, cond2 := checkConditions(smap[c], smap[pFChromosomeIndex])

		fmt.Println("Codition check result: ", cond1, cond2)

		if cond1 && cond2 {
			fmt.Printf("\n Chromosome: %d index: %v \n", c, i)
			paretoFrontCheckResult := &pFcheckMapResult{}
			paretoFrontCheckResult.result = "true"
			paretoFrontCheckResult.currChromosomeIndexNotInFront = currChromosomeIndex
			paretoFrontCheckResult.pFChromosomeIndex = i
			currChromosomeIndexTopFchromosomeIndex[currChromosomeIndex] = i
			chromsomeCheckMap["true"] = currChromosomeIndexTopFchromosomeIndex
			updateParetoFront(pF, i, index, paretoFrontCheckResult)
		}

		if !cond1 && cond2 || cond1 && !cond2 {
			fmt.Printf("\n Chromosome: %d index: %v \n", c, i)
			currChromosomeIndexTopFchromosomeIndex[currChromosomeIndex] = i
			chromsomeCheckMap["non-dominating"] = currChromosomeIndexTopFchromosomeIndex
			addToParetoFront(pF, currChromosomeIndex, index)
		}

		//if the selected chromosome is worse off than any in the pareto front stop the check
		if !cond1 && !cond2 {
			break
		}
	}

	//loop through the chromosome check map
	for cond, chromosomeMapToPf := range chromsomeCheckMap {

		//get the number of total items for the curr chromosome compared with the pareto front elements in the check map
		chromosomeMapToPfLength := len(chromosomeMapToPf) - 1

		//check if the chromosome indicies in checkmap is less than current indicies in the pareto front
		if chromosomeMapToPfLength < totalItemsInPf {

			//remove all elements that are non-dominating with the current chromosome based on the checkmap
			if cond == "non-dominating" {
				for _, pFindex := range chromosomeMapToPf {
					deleteFromParetoFront(pF, index, pFindex)
				}
			}
		} 
	}

	fmt.Println("Pareto Front currently: ", Pf)
}