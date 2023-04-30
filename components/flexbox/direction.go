package flexbox

import "github.com/mieubrisse/box-layout-test/components/flexbox_item"

// The direction that the flexbox ought to be layed out in
type Direction interface {
	// The functions in this interface are used internally in the flexbox to do its calculations

	getMainAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

	getCrossAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

	getActualWidths(desiredWidths []int, items []flexbox_item.FlexboxItem) axisCalculationResults

	getDesiredHeights(width int, items []flexbox_item.FlexboxItem) []int
}

// Row lays out the flexbox items in a row, left to right
// The flex direction will be horizontal
// Corresponds to "flex-direction: row" in CSS
func Row() Direction {
	return &directionImpl{
		mainAxisMaxDimensionValueGetter: func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
			return item.GetMaxHeight()
		},
		crossAxisMaxDimensionValueGetter: func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
			return item.GetMaxWidth()
		},
	}
}

// Column lays out the flexbox items in a column, top to bottom
// The flex direction will be vertical
// Corresponds to "flex-direction: column" in CSS
func Column() Direction {

}

// TODO column

// ====================================================================================================
//
//	Private
//
// ====================================================================================================
type itemDimensionValueGetter func(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue
type directionImpl struct {
	mainAxisMaxDimensionValueGetter  itemDimensionValueGetter
	crossAxisMaxDimensionValueGetter itemDimensionValueGetter
}

func (d directionImpl) getMainAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
	return d.mainAxisMaxDimensionValueGetter(item)
}

func (d directionImpl) getCrossAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue {
	return d.crossAxisMaxDimensionValueGetter(item)
}

func (d directionImpl) getDesiredWidths(items []flexbox_item.FlexboxItem) []int {
	//TODO implement me
	panic("implement me")
}

func (d directionImpl) getDesiredHeights(width int, items []flexbox_item.FlexboxItem) []int {
	//TODO implement me
	panic("implement me")
}
