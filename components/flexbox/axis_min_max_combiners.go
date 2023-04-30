package flexbox

import "github.com/mieubrisse/box-layout-test/utilities"

// Combines the mins & maxes of all the items to return a min & max for the parent
type axisDimensionMinMaxCombiner func(mins []int, maxes []int) (min, max int)

// Combines mins & maxes for items on the cross axis to get a min & max parent
func crossAxisDimensionMinMaxCombiner(mins []int, maxes []int) (int, int) {
	min, max := 0, 0
	for idx := range mins {
		min = utilities.GetMaxInt(min, mins[idx])
		max = utilities.GetMaxInt(max, maxes[idx])
	}
	return min, max
}

// Combines mins & maxes for items on the cross axis to get a min & max parent
func mainAxisDimensionMinMaxCombiner(mins []int, maxes []int) (int, int) {
	min, max := 0, 0
	for idx := range mins {
		min += mins[idx]
		max += maxes[idx]
	}
	return min, max
}
