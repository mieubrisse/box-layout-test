package flexbox

import (
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
)

type itemDimensionValueGetter func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

type axisSizeCalculator func(
	items []flexbox_item.FlexboxItem,
	desiredSizes []int,
	spaceAvailable int,
	maxDimensionValueGetter itemDimensionValueGetter,
) axisSizeCalculationResults

type axisSizeCalculationResults struct {
	actualSizes []int

	spaceUsedByChildren int
}

// TODO move to be a function on the axis?
func calculateActualCrossAxisSizes(
	items []flexbox_item.FlexboxItem,
	desiredSizes []int,
	// How much space is available in the cross axis
	spaceAvailable int,
	// How to get the item's desired max value in the cross axis
	crossAxisMaxValueGetter itemDimensionValueGetter,
) axisSizeCalculationResults {
	actualSizes := make([]int, len(items))

	// The space used in the cross axis is the max across all children
	maxSpaceUsed := 0
	for idx, item := range items {
		actualSize := desiredSizes[idx]
		if crossAxisMaxValueGetter(item).ShouldGrow() {
			actualSize = utilities.GetMaxInt(actualSize, spaceAvailable)
		}

		// Ensure we don't overrun
		actualSize = utilities.GetMinInt(actualSize, spaceAvailable)

		actualSizes[idx] = actualSize
		maxSpaceUsed = utilities.GetMaxInt(actualSize, maxSpaceUsed)
	}
	return axisSizeCalculationResults{
		actualSizes:         actualSizes,
		spaceUsedByChildren: maxSpaceUsed,
	}
}

func calculateActualMainAxisSizes(
	items []flexbox_item.FlexboxItem,
	desiredSizes []int,
	spaceAvailable int,
	mainAxisValueGetter itemDimensionValueGetter,
) axisSizeCalculationResults {
	totalDesiredSize := 0
	for _, desiredSize := range desiredSizes {
		totalDesiredSize += desiredSize
	}

	actualSizes := desiredSizes
	freeSpace := spaceAvailable - totalDesiredSize
	// The "grow" case
	if freeSpace > 0 {
		weights := make([]int, len(items))
		for idx, item := range items {
			if mainAxisValueGetter(item).ShouldGrow() {
				// TODO deal with actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		actualSizes = distributeSpaceByWeight(freeSpace, desiredSizes, weights)
		// The "shrink" case
	} else if freeSpace < 0 {
		// We use desired sizes as the weight, so that
		actualSizes = distributeSpaceByWeight(freeSpace, desiredSizes, desiredSizes)
	}

	totalSpaceUsed := 0
	for _, spaceUsedByChild := range actualSizes {
		totalSpaceUsed += spaceUsedByChild
	}

	return axisSizeCalculationResults{
		actualSizes:         actualSizes,
		spaceUsedByChildren: totalSpaceUsed,
	}
}

// Distributes the space (which can be negative) across the children, using the weight as a bias for how to allocate
// The only scenario where no space will be distributed is if there is no total weight
// If the space does get distributed, it's guaranteed to be done exactly (no more or less will remain)
func distributeSpaceByWeight(spaceToAllocate int, inputSizes []int, weights []int) []int {
	result := make([]int, len(inputSizes))
	for idx, inputSize := range inputSizes {
		result[idx] = inputSize
	}

	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// watch out for divide-by-zero
	if totalWeight == 0 {
		return result
	}

	desiredSpaceAllocated := float64(0)
	actualSpaceAllocated := 0
	for idx, size := range inputSizes {
		result[idx] = size

		// Dump any remaining space for the last item (it should always be at most 1
		// in any direction)
		if idx == len(inputSizes)-1 {
			result[idx] += spaceToAllocate - actualSpaceAllocated
			break
		}

		weight := weights[idx]
		share := float64(weight) / float64(totalWeight)

		// Because we can only display lines in integer numbers, but flexing
		// will yield float scale ratios, no matter what space we give each item
		// our integer value will always be off from the float value
		// This algorithm is to ensure that we're always rounding in the direction
		// that pushes us closer to our desired allocation (rather than naively rounding up or down)
		desiredSpaceForItem := share * float64(spaceToAllocate)
		var actualSpaceForItem int
		if desiredSpaceAllocated < float64(actualSpaceAllocated) {
			// If we're under our desired allocation, round up to try and get closer
			actualSpaceForItem = int(desiredSpaceForItem + 1)
		} else {
			// If we're at or over our desired allocation, round down (so we either stay there or get closer by undershooting)
			actualSpaceForItem = int(desiredSpaceForItem)
		}

		result[idx] += actualSpaceForItem
		desiredSpaceAllocated += desiredSpaceForItem
		actualSpaceAllocated += actualSpaceForItem
	}

	return result
}
