package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
	"strings"
)

// ========================= NOTE ==============================
// Flexboxes can go in any direction. I'm going to us "main axis size"
// and "cross axis size" to refer to these. I'm going to refer to them
// as "MAS" and "CAS" throughout this piece of code.
// ========================= NOTE ==============================

type Direction int

const (
	// Row lays out the flexbox items in a row, left to right
	// The flex direction will be horizontal
	// Corresponds to "flex-direction: row" in CSS
	Row Direction = iota

	// Column lays out the flexbox items in a column, top to bottom
	// The flex direction will be vertical
	// Corresponds to "flex-direction: column" in CSS
	Column
)

// When the children don't completely fill the box, where to put teh
// Corresponds to "justify-content" in CSS
type MainAxisAlignment int

const (
	// Elements will be at the start of the flexbox (as determined by the Direction)
	// Corresponds to "flex-justify: flex-start"
	MainAxisStart MainAxisAlignment = iota

	// NOTE: in order to see this in effect, you must have
	// Corresponds to "flex-justify: center"
	MainAxisCenter

	// Elements will be pushed to the end of the flexbox (as determined by the Direction)
	// Corresponds to "flex-justify: flex-end"
	MainAxisEnd

	// TODO space-between, space-around, space-evenly: https://css-tricks.com/snippets/css/a-guide-to-flexbox/
)

// CrossAxisAlignment controls where to put children when the child's height doesn't completely fill the cross axis
// Corresponds to "align-items" in CSS
type CrossAxisAlignment int

const (
	// CrossAxisStart arranges items at the start of the cross axis, when there is extra space
	// E.g. when the flexbox direction is horizontal, this will push items to the top
	// Coreresponds to "align-items: flex-start" in CSS
	CrossAxisStart CrossAxisAlignment = iota

	// CrossAxisMiddle arranges items in the center of the cross axis, when there is extra space
	// E.g. when the flexbox direction is horizontal, this will push items to the top
	// Coreresponds to "align-items: center" in CSS
	CrossAxisMiddle

	// CrossAxisEnd arranges items at the end of the cross axis, when there is extra space
	// Coreresponds to "align-items: flex-end" in CSS
	CrossAxisEnd
)

// TODO make an interface
type Flexbox struct {
	children []flexbox_item.FlexboxItem

	// TODO make configurable on left and right
	padding int

	border lipgloss.Border

	horizontalJustify MainAxisAlignment
	verticalJustify   CrossAxisAlignment

	// -------------------- Calculation Caching -----------------------
	// Cache of the min/max widths/heights across all children
	allChildrenDimensionsCache components.DimensionsCache

	// Cached result of calculating child widths in the GetContentHeightForGivenWidth phase
	// We do this so we
	childWidthsCalculationCache calculateChildWidthsResult
}

// Convenience constructor for a box with a single element
func NewWithContent(component components.Component, opts ...flexbox_item.FlexboxItemOpt) *Flexbox {
	item := flexbox_item.NewItem(component)
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
		padding:           0,
		children:          make([]flexbox_item.FlexboxItem, 0),
		border:            lipgloss.Border{},
		horizontalJustify: MainAxisStart,
	}
}

// TODO make this configurable per side
func (b *Flexbox) SetPadding(padding int) *Flexbox {
	b.padding = padding
	return b
}

func (b *Flexbox) SetBorder(border lipgloss.Border) *Flexbox {
	b.border = border
	return b
}

func (b *Flexbox) SetChildren(children []flexbox_item.FlexboxItem) *Flexbox {
	b.children = children
	return b
}

func (b *Flexbox) SetHorizontalJustify(justify MainAxisAlignment) *Flexbox {
	b.horizontalJustify = justify
	return b
}

func (b *Flexbox) SetVerticalJustify(justify CrossAxisAlignment) *Flexbox {
	b.verticalJustify = justify
	return b
}

func (b *Flexbox) GetContentMinMax() (minWidth int, maxWidth int, minHeight int, maxHeight int) {
	// TODO allow column layout

	var childrenMinWidth, childrenMaxWidth, childrenMinHeight, childrenMaxHeight int
	for _, item := range b.children {
		itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight := item.GetContentMinMax()

		// Calculate the maxes
		childrenMinWidth = utilities.GetMaxInt(childrenMinWidth, itemMinWidth)
		childrenMaxWidth = utilities.GetMaxInt(childrenMaxWidth, itemMaxWidth)
		childrenMinHeight = utilities.GetMaxInt(childrenMinHeight, itemMinHeight)
		childrenMaxHeight = utilities.GetMaxInt(childrenMaxHeight, itemMaxHeight)

	}

	b.allChildrenDimensionsCache = components.DimensionsCache{
		MinWidth:  childrenMinWidth,
		MaxWidth:  childrenMaxWidth,
		MinHeight: childrenMinHeight,
		MaxHeight: childrenMaxHeight,
	}

	additionalNonContentWidth := b.calculateAdditionalNonContentWidth()
	minWidth = childrenMinWidth + additionalNonContentWidth
	maxWidth = childrenMaxWidth + additionalNonContentWidth

	additionalNonContentHeight := b.calculateAdditionalNonContentHeight()
	minHeight = childrenMinHeight + additionalNonContentHeight
	maxHeight = childrenMaxHeight + additionalNonContentHeight

	return
}

func (b *Flexbox) GetContentHeightForGivenWidth(width int) int {
	// TODO cache this result!!!!

	// Width
	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	widthAvailableForChildren := utilities.GetMaxInt(0, width-nonContentWidthNeeded)
	calculationResult := b.calculateMainAxisWidths(widthAvailableForChildren)

	// Cache the result, so we don't have to do this again during View
	b.childWidthsCalculationCache = calculationResult

	maxDesiredItemHeight := 0
	for idx, item := range b.children {
		itemWidth := calculationResult.childWidths[idx]
		desiredItemHeight := item.GetContentHeightForGivenWidth(itemWidth)

		maxDesiredItemHeight = utilities.GetMaxInt(desiredItemHeight, maxDesiredItemHeight)
	}

	nonContentHeightNeeded := b.calculateAdditionalNonContentHeight()
	return maxDesiredItemHeight + nonContentHeightNeeded
}

func (b *Flexbox) View(width int, height int) string {
	// Width
	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	widthAvailableForChildren := utilities.GetMaxInt(0, width-nonContentWidthNeeded)
	widthNotUsedByChildren := widthAvailableForChildren - b.childWidthsCalculationCache.totalWidthUsed

	// Height
	additionalNonContentHeight := b.calculateAdditionalNonContentHeight()
	heightAvailableForChildren := utilities.GetMaxInt(0, height-additionalNonContentHeight)
	childHeights, maxHeightUsedByChildren := b.calculateChildHeights(
		b.childWidthsCalculationCache.childWidths,
		heightAvailableForChildren,
	)
	heightNotUsedByChildren := heightAvailableForChildren - maxHeightUsedByChildren

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allContentFragments := make([]string, len(b.children))
	for idx, item := range b.children {
		childWidth := b.childWidthsCalculationCache.childWidths[idx]
		childHeight := childHeights[idx]
		childStr := item.View(childWidth, childHeight)

		allContentFragments[idx] = childStr
	}

	// Justify horizontally
	switch b.horizontalJustify {
	case MainAxisStart:
		pad := strings.Repeat(" ", widthNotUsedByChildren)
		allContentFragments = append(allContentFragments, pad)
	case MainAxisEnd:
		pad := strings.Repeat(" ", widthNotUsedByChildren)
		allContentFragments = append([]string{pad}, allContentFragments...)
	case MainAxisCenter:
		leftPadSize := widthNotUsedByChildren / 2
		rightPadSize := widthNotUsedByChildren - leftPadSize
		leftPad := strings.Repeat(" ", leftPadSize)
		rightPad := strings.Repeat(" ", rightPadSize)

		newContentFragments := append([]string{leftPad}, allContentFragments...)
		newContentFragments = append(newContentFragments, rightPad)
		allContentFragments = newContentFragments
	}

	// TODO allow other align types

	// Justify vertically
	var content string
	switch b.verticalJustify {
	case CrossAxisStart:
		content = lipgloss.JoinHorizontal(lipgloss.Top, allContentFragments...)
		content += strings.Repeat("\n", heightNotUsedByChildren)
	case CrossAxisEnd:
		content = lipgloss.JoinHorizontal(lipgloss.Bottom, allContentFragments...)
		content = strings.Repeat("\n", heightNotUsedByChildren) + content
	case CrossAxisMiddle:
		content = lipgloss.JoinHorizontal(lipgloss.Center, allContentFragments...)
		topPadSize := heightNotUsedByChildren / 2
		bottomPadSize := heightNotUsedByChildren - topPadSize
		topPad := strings.Repeat("\n", topPadSize)
		bottomPad := strings.Repeat("\n", bottomPadSize)
		content = topPad + content + bottomPad
	}

	result := lipgloss.NewStyle().
		Padding(b.padding).
		Border(b.border).
		Render(content)

	return result
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (b Flexbox) calculateAdditionalNonContentWidth() int {
	nonContentWidthAdditions := []int{
		// Padding
		2 * b.padding,

		// Border
		int(b.border.GetLeftSize() + b.border.GetRightSize()),
	}

	totalNonContentWidthAdded := int(0)
	for _, addition := range nonContentWidthAdditions {
		totalNonContentWidthAdded += addition
	}
	return totalNonContentWidthAdded
}

func (b Flexbox) calculateAdditionalNonContentHeight() int {
	nonContentHeightAdditions := []int{
		// Padding
		2 * b.padding,

		// Border
		int(b.border.GetTopSize() + b.border.GetBottomSize()),
	}

	totalNonContentHeightAdded := int(0)
	for _, addition := range nonContentHeightAdditions {
		totalNonContentHeightAdded += addition
	}
	return totalNonContentHeightAdded
}

// Distributes the space
// The only scenario where no space will be distributed is if there is no total weight
// If the space does get distributed, it's guaranteed to be done exactly (no more or less will remain)
func addSpaceByWeight(spaceToAllocate int, inputSizes []int, weights []int) []int {
	result := make([]int, len(inputSizes))
	for idx, inputSize := range inputSizes {
		result[idx] = inputSize
	}

	totalWeight := int(0)
	for _, weight := range weights {
		totalWeight += weight
	}

	// watch out for divide-by-zero
	if totalWeight == 0 {
		return result
	}

	desiredSpaceAllocated := float64(0)
	actualSpaceAllocated := int(0)
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
			// Round up
			actualSpaceForItem = int(desiredSpaceForItem + 1)
		} else {
			// Round down
			actualSpaceForItem = int(desiredSpaceForItem)
		}

		result[idx] += actualSpaceForItem
		desiredSpaceAllocated += desiredSpaceForItem
		actualSpaceAllocated += actualSpaceForItem
	}

	return result
}

// Calculates
func (b Flexbox) calculateMainAxisItemSizes(mins []int, maxes []int) ([]int, int) {
	// First, add up the space each item would like in a perfect world
	totalMinSpaceDesired := 0 // Under this value, the flexbox will simply truncate
	totalMaxSpaceDesired := 0 // Above this value, only the items that have MaxAvailable set for this dimension will expand

	// First, add up the total width the items would like in a perfect world
	childWidths := make([]int, numChildren)
	maxConstrainedChildWidths := make([]int, numChildren)
	for idx, item := range b.children {
		itemMinWidth, itemMaxWidth, _, _ := item.GetContentMinMax()

		totalMinWidthDesired += itemMinWidth
		totalMaxWidthDesired += itemMaxWidth

		childWidths[idx] = itemMinWidth
		maxConstrainedChildWidths[idx] = itemMaxWidth
	}

	// When min_desired_size < space_available < max_desired_size, scale everyone up equally between their
	// min and max sizes
	if widthAvailableForChildren > totalMinWidthDesired {
		weights := make([]int, numChildren)
		for idx, minChildSize := range childWidths {
			maxChildSize := maxConstrainedChildWidths[idx]
			childExpansionRange := maxChildSize - minChildSize
			weights[idx] = childExpansionRange
		}

		spaceToDistributeEvenly := utilities.GetMinInt(
			widthAvailableForChildren-totalMinWidthDesired,
			totalMaxWidthDesired-totalMinWidthDesired,
		)

		childWidths = addSpaceByWeight(spaceToDistributeEvenly, childWidths, weights)
	}

	// When width > max_desired_size, continue to scale only the elements whose max size is MaxAvailable
	if widthAvailableForChildren > totalMaxWidthDesired {
		weights := make([]int, numChildren)
		for idx, item := range b.children {
			if item.GetMaxWidth().ShouldGrow() {
				// TODO use actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		spaceToDistributeToExpanders := widthAvailableForChildren - totalMaxWidthDesired

		childWidths = addSpaceByWeight(spaceToDistributeToExpanders, childWidths, weights)
	}

	// Finally, ensure that the child widths don't exceed our available space
	totalWidthUsed := int(0)
	widthAvailable := widthAvailableForChildren
	for idx, childWidth := range childWidths {
		actualWidth := utilities.GetMinInt(widthAvailable, childWidth)
		childWidths[idx] = actualWidth

		widthAvailable -= actualWidth
		totalWidthUsed += actualWidth
	}

	return calculateChildWidthsResult{
		childWidths:    childWidths,
		totalWidthUsed: totalWidthUsed,
	}
}

// Calculates the sizes for children
func (b Flexbox) calculateMainAxisWidths(widthAvailableForChildren int) calculateChildWidthsResult {

	numChildren := len(b.children)

	// First, add up the total width the items would like in a perfect world
	totalMinWidthDesired := int(0) // Under this value, the flexbox will simply truncate
	totalMaxWidthDesired := int(0) // Above this value, only the MaxAvailable items will expand
	childWidths := make([]int, numChildren)
	maxConstrainedChildWidths := make([]int, numChildren)
	for idx, item := range b.children {
		itemMinWidth, itemMaxWidth, _, _ := item.GetContentMinMax()

		totalMinWidthDesired += itemMinWidth
		totalMaxWidthDesired += itemMaxWidth

		childWidths[idx] = itemMinWidth
		maxConstrainedChildWidths[idx] = itemMaxWidth
	}

	// When min_desired_size < space_available < max_desired_size, scale everyone up equally between their
	// min and max sizes
	if widthAvailableForChildren > totalMinWidthDesired {
		weights := make([]int, numChildren)
		for idx, minChildSize := range childWidths {
			maxChildSize := maxConstrainedChildWidths[idx]
			childExpansionRange := maxChildSize - minChildSize
			weights[idx] = childExpansionRange
		}

		spaceToDistributeEvenly := utilities.GetMinInt(
			widthAvailableForChildren-totalMinWidthDesired,
			totalMaxWidthDesired-totalMinWidthDesired,
		)

		childWidths = addSpaceByWeight(spaceToDistributeEvenly, childWidths, weights)
	}

	// When width > max_desired_size, continue to scale only the elements whose max size is MaxAvailable
	if widthAvailableForChildren > totalMaxWidthDesired {
		weights := make([]int, numChildren)
		for idx, item := range b.children {
			if item.GetMaxWidth().ShouldGrow() {
				// TODO use actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		spaceToDistributeToExpanders := widthAvailableForChildren - totalMaxWidthDesired

		childWidths = addSpaceByWeight(spaceToDistributeToExpanders, childWidths, weights)
	}

	// Finally, ensure that the child widths don't exceed our available space
	totalWidthUsed := int(0)
	widthAvailable := widthAvailableForChildren
	for idx, childWidth := range childWidths {
		actualWidth := utilities.GetMinInt(widthAvailable, childWidth)
		childWidths[idx] = actualWidth

		widthAvailable -= actualWidth
		totalWidthUsed += actualWidth
	}

	return calculateChildWidthsResult{
		childWidths:    childWidths,
		totalWidthUsed: totalWidthUsed,
	}
}

func (b Flexbox) calculateChildHeights(childWidths []int, heightAvailable int) ([]int, int) {
	// TODO cache these results???
	childHeights := make([]int, len(b.children))
	maxHeightUsed := 0
	for idx, item := range b.children {
		width := childWidths[idx]
		height := item.GetContentHeightForGivenWidth(width)

		if item.GetMaxHeight().ShouldGrow() {
			height = utilities.GetMaxInt(height, heightAvailable)
		}

		// Ensure we don't overrun
		height = utilities.GetMinInt(heightAvailable, height)

		childHeights[idx] = height
		maxHeightUsed = utilities.GetMaxInt(height, maxHeightUsed)
	}
	return childHeights, maxHeightUsed
}
