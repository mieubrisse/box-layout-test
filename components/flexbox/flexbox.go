package flexbox

import (
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
)

// NOTE: This class does some stateful caching, so when you're testing methods like "View" make sure you call the
// full flow of GetContentMinMax -> GetContentHeightForGivenWidth -> View as necessary

// TODO make an interface
type Flexbox struct {
	children []flexbox_item.FlexboxItem

	direction Direction

	horizontalAlignment AxisAlignment
	verticalAlignment   AxisAlignment

	// -------------------- Calculation Caching -----------------------
	// The widths each child desires (cached between GetContentMinMax and GetContentHeightForGivenWidth
	desiredChildWidthsCache []int

	// The actual widths each child will get (cached between GetContentHeightForGivenWidth and View)
	actualChildWidthsCache axisSizeCalculationResults

	// The desired height each child wants given its width (cached between GetContentHeightForGivenWidth and View)
	desiredChildHeightsGivenWidthCache []int
}

// Convenience constructor for a box with a single element
func NewWithContent(component components.Component, opts ...flexbox_item.FlexboxItemOpt) *Flexbox {
	item := flexbox_item.New(component)
	for _, opt := range opts {
		opt(item)
	}
	return NewWithContents(item)
}

// Convenience constructor for a box with multiple elements
func NewWithContents(items ...flexbox_item.FlexboxItem) *Flexbox {
	return New().SetChildren(items)
}

func New() *Flexbox {
	return &Flexbox{
		children:                           make([]flexbox_item.FlexboxItem, 0),
		direction:                          Row,
		horizontalAlignment:                AlignStart,
		verticalAlignment:                  AlignStart,
		desiredChildWidthsCache:            nil,
		actualChildWidthsCache:             axisSizeCalculationResults{},
		desiredChildHeightsGivenWidthCache: nil,
	}
}

func (b *Flexbox) SetChildren(children []flexbox_item.FlexboxItem) *Flexbox {
	b.children = children
	return b
}

func (b *Flexbox) SetDirection(direction Direction) *Flexbox {
	b.direction = direction
	return b
}

func (b *Flexbox) SetHorizontalAlignment(alignment AxisAlignment) *Flexbox {
	b.horizontalAlignment = alignment
	return b
}

func (b *Flexbox) SetVerticalAlignment(alignment AxisAlignment) *Flexbox {
	b.verticalAlignment = alignment
	return b
}

func (b *Flexbox) GetContentMinMax() (minWidth int, maxWidth int, minHeight int, maxHeight int) {
	// TODO allow column layout

	var childrenMinWidth, childrenMaxWidth, childrenMinHeight, childrenMaxHeight int
	b.desiredChildWidthsCache = make([]int, len(b.children))
	for idx, item := range b.children {
		itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight := item.GetContentMinMax()

		// Cache the item's max width; we'll need it in GetContentHeightForGivenWidth
		b.desiredChildWidthsCache[idx] = itemMaxWidth

		// Calculate the maxes
		childrenMinWidth = utilities.GetMaxInt(childrenMinWidth, itemMinWidth)
		childrenMaxWidth = utilities.GetMaxInt(childrenMaxWidth, itemMaxWidth)
		childrenMinHeight = utilities.GetMaxInt(childrenMinHeight, itemMinHeight)
		childrenMaxHeight = utilities.GetMaxInt(childrenMaxHeight, itemMaxHeight)
	}

	minWidth = childrenMinWidth
	maxWidth = childrenMaxWidth

	minHeight = childrenMinHeight
	maxHeight = childrenMaxHeight

	return
}

func (b *Flexbox) GetContentHeightForGivenWidth(width int) int {
	if width == 0 {
		return 0
	}

	// Width
	shouldGrowWidths := make([]bool, len(b.children))
	for idx, item := range b.children {
		shouldGrowWidths[idx] = item.GetMaxWidth().ShouldGrow()
	}
	actualWidthsCalcResults := b.direction.getActualWidths(b.desiredChildWidthsCache, shouldGrowWidths, width)

	// Cache the result, so we don't have to recalculate it in View
	b.actualChildWidthsCache = actualWidthsCalcResults

	result := 0
	desiredHeights := make([]int, len(b.children))
	for idx, item := range b.children {
		actualWidth := actualWidthsCalcResults.actualSizes[idx]
		desiredHeight := item.GetContentHeightForGivenWidth(actualWidth)

		desiredHeights[idx] = desiredHeight
		result = utilities.GetMaxInt(result, desiredHeight)
	}

	// Cache the result, so we don't have to recalculate it in View
	b.desiredChildHeightsGivenWidthCache = desiredHeights

	return result
}

func (b *Flexbox) View(width int, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	actualWidths := b.actualChildWidthsCache.actualSizes
	// widthNotUsedByChildren := utilities.GetMaxInt(0, width-b.actualChildWidthsCache.spaceUsedByChildren)

	shouldGrowHeights := make([]bool, len(b.children))
	for idx, item := range b.children {
		shouldGrowHeights[idx] = item.GetMaxHeight().ShouldGrow()
	}
	actualHeightsCalcResult := b.direction.getActualHeights(b.desiredChildHeightsGivenWidthCache, shouldGrowHeights, height)

	actualHeights := actualHeightsCalcResult.actualSizes
	// heightNotUsedByChildren := utilities.GetMaxInt(0, height-actualHeightsCalcResult.spaceUsedByChildren)

	// Now render each child
	allContentFragments := make([]string, len(b.children))
	for idx, item := range b.children {
		childWidth := actualWidths[idx]
		childHeight := actualHeights[idx]
		childStr := item.View(childWidth, childHeight)

		allContentFragments[idx] = childStr
	}

	content := b.direction.renderContentFragments(allContentFragments, width, height, b.horizontalAlignment, b.verticalAlignment)

	/*
		// Justify main axis
		switch b.horizontalAlignment {
		case AlignStart:
			pad := strings.Repeat(" ", widthNotUsedByChildren)
			content += pad
		case AlignEnd:
			pad := strings.Repeat(" ", widthNotUsedByChildren)
			content = pad + content
		case AlignCenter:
			leftPadSize := widthNotUsedByChildren / 2
			rightPadSize := widthNotUsedByChildren - leftPadSize
			leftPad := strings.Repeat(" ", leftPadSize)
			rightPad := strings.Repeat(" ", rightPadSize)

			newContentFragments := append([]string{leftPad}, allContentFragments...)
			newContentFragments = append(newContentFragments, rightPad)
			allContentFragments = newContentFragments
		}

		// Justify cross axis
		var content string
		switch b.verticalAlignment {
		case AlignStart:
			content = lipgloss.JoinHorizontal(lipgloss.Top, allContentFragments...)
			content += strings.Repeat("\n", heightNotUsedByChildren)
		case AlignEnd:
			content = lipgloss.JoinHorizontal(lipgloss.Bottom, allContentFragments...)
			content = strings.Repeat("\n", heightNotUsedByChildren) + content
		case AlignCenter:
			content = lipgloss.JoinHorizontal(lipgloss.Center, allContentFragments...)
			topPadSize := heightNotUsedByChildren / 2
			bottomPadSize := heightNotUsedByChildren - topPadSize
			topPad := strings.Repeat("\n", topPadSize)
			bottomPad := strings.Repeat("\n", bottomPadSize)
			content = topPad + content + bottomPad
		}

	*/

	return content
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
