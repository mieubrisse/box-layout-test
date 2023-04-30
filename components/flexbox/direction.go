package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
)

// The direction that the flexbox ought to be layed out in
type Direction interface {
	// The functions in this interface are used internally in the flexbox to do its calculations

	/*
		getMainAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

		getCrossAxisMaxDimensionValue(item flexbox_item.FlexboxItem) flexbox_item.FlexboxItemDimensionValue

	*/

	getContentSizes(items []flexbox_item.FlexboxItem) (minWidth, maxWidth, minHeight, maxHeight int)

	getActualWidths(desiredWidths []int, shouldGrow []bool, widthAvailable int) axisSizeCalculationResults

	getActualHeights(desiredHeights []int, shouldGrow []bool, heightAvailable int) axisSizeCalculationResults

	renderContentFragments(contentFragments []string, width int, height int, horizontalAlignment AxisAlignment, verticalAlignment AxisAlignment) string
}

// Row lays out the flexbox items in a row, left to right
// The flex direction will be horizontal
// Corresponds to "flex-direction: row" in CSS
var Row = &directionImpl{
	actualWidthCalculator:  calculateActualMainAxisSizes,
	actualHeightCalculator: calculateActualCrossAxisSizes,
	minMaxWidthCombiner:    mainAxisDimensionMinMaxCombiner,
	minMaxHeightCombiner:   crossAxisDimensionMinMaxCombiner,
	contentFragmentRenderer: func(contentFragments []string, width int, height int, horizontalAlign AxisAlignment, verticalAlign AxisAlignment) string {
		joined := lipgloss.JoinHorizontal(lipgloss.Position(verticalAlign), contentFragments...)
		horizontallyPlaced := lipgloss.PlaceHorizontal(width, lipgloss.Position(horizontalAlign), joined)
		return lipgloss.PlaceVertical(height, lipgloss.Position(verticalAlign), horizontallyPlaced)
	},
}

// Column lays out the flexbox items in a column, top to bottom
// The flex direction will be vertical
// Corresponds to "flex-direction: column" in CSS
var Column = &directionImpl{
	actualWidthCalculator:  calculateActualCrossAxisSizes,
	actualHeightCalculator: calculateActualMainAxisSizes,
	minMaxWidthCombiner:    crossAxisDimensionMinMaxCombiner,
	minMaxHeightCombiner:   mainAxisDimensionMinMaxCombiner,
	contentFragmentRenderer: func(contentFragments []string, width int, height int, horizontalAlign AxisAlignment, verticalAlign AxisAlignment) string {
		joined := lipgloss.JoinVertical(lipgloss.Position(horizontalAlign), contentFragments...)
		horizontallyPlaced := lipgloss.PlaceHorizontal(width, lipgloss.Position(horizontalAlign), joined)
		return lipgloss.PlaceVertical(height, lipgloss.Position(verticalAlign), horizontallyPlaced)
	},
}

// ====================================================================================================
//
//	Private
//
// ====================================================================================================
type directionImpl struct {
	actualWidthCalculator   axisSizeCalculator
	actualHeightCalculator  axisSizeCalculator
	minMaxWidthCombiner     axisDimensionMinMaxCombiner
	minMaxHeightCombiner    axisDimensionMinMaxCombiner
	contentFragmentRenderer func(contentFragments []string, width int, height int, horizontalAlign AxisAlignment, verticalAlign AxisAlignment) string
}

func (a directionImpl) getContentSizes(items []flexbox_item.FlexboxItem) (int, int, int, int) {
	childMinWidths := make([]int, len(items))
	childMaxWidths := make([]int, len(items))
	childMinHeights := make([]int, len(items))
	childMaxHeights := make([]int, len(items))
	for idx, item := range items {
		childMinWidths[idx], childMaxWidths[idx], childMinHeights[idx], childMaxHeights[idx] = item.GetContentMinMax()
	}

	minWidth, maxWidth := a.minMaxWidthCombiner(childMinWidths, childMaxWidths)
	minHeight, maxHeight := a.minMaxHeightCombiner(childMinHeights, childMaxHeights)

	return minWidth, maxWidth, minHeight, maxHeight
}

func (r directionImpl) getActualWidths(desiredWidths []int, shouldGrow []bool, widthAvailable int) axisSizeCalculationResults {
	return r.actualWidthCalculator(
		desiredWidths,
		shouldGrow,
		widthAvailable,
	)
}

func (r directionImpl) getActualHeights(desiredHeights []int, shouldGrow []bool, heightAvailable int) axisSizeCalculationResults {
	return r.actualHeightCalculator(
		desiredHeights,
		shouldGrow,
		heightAvailable,
	)
}

func (r directionImpl) renderContentFragments(contentFragments []string, width int, height int, horizontalAlign AxisAlignment, verticalAlign AxisAlignment) string {
	return r.contentFragmentRenderer(contentFragments, width, height, horizontalAlign, verticalAlign)
}
