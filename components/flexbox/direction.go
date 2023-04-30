package flexbox

import "github.com/mieubrisse/box-layout-test/components/flexbox_item"

// The direction that the flexbox ought to be layed out in
type Direction interface {
	// The functions in this interface are used internally in the flexbox to do its calculations

	/*
		getMainAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

		getCrossAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

	*/

	getActualWidths(desiredWidths []int, items []flexbox_item.FlexboxItem, widthAvailable int) axisSizeCalculationResults

	getActualHeights(desiredHeights []int, items []flexbox_item.FlexboxItem, heightAvailable int) axisSizeCalculationResults

	renderContentFragments(contentFragments []string) string
}

// Row lays out the flexbox items in a row, left to right
// The flex direction will be horizontal
// Corresponds to "flex-direction: row" in CSS
func Row() Direction {
	return &directionImpl{
		actualWidthCalculator:  calculateActualMainAxisSizes,
		actualHeightCalculator: calculateActualCrossAxisSizes,
	}
}

// Column lays out the flexbox items in a column, top to bottom
// The flex direction will be vertical
// Corresponds to "flex-direction: column" in CSS
func Column() Direction {
	return &directionImpl{
		actualWidthCalculator:  calculateActualCrossAxisSizes,
		actualHeightCalculator: calculateActualMainAxisSizes,
	}

}

// TODO column

// ====================================================================================================
//
//	Private
//
// ====================================================================================================
type directionImpl struct {
	actualWidthCalculator  axisSizeCalculator
	actualHeightCalculator axisSizeCalculator
}

func (r directionImpl) getActualWidths(desiredWidths []int, items []flexbox_item.FlexboxItem, widthAvailable int) axisSizeCalculationResults {
	return r.actualWidthCalculator(
		items,
		desiredWidths,
		widthAvailable,
		func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
			return item.GetMaxWidth()
		},
	)
}

func (r directionImpl) getActualHeights(desiredHeights []int, items []flexbox_item.FlexboxItem, heightAvailable int) axisSizeCalculationResults {
	return r.actualHeightCalculator(
		items,
		desiredHeights,
		heightAvailable,
		func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
			return item.GetMaxHeight()
		},
	)
}

func (r directionImpl) renderContentFragments(contentFragments []string) string {
	//TODO implement me
	panic("implement me")
}
