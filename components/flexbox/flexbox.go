package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
	"strings"
)

// TODO make this "main axis justify"
// When the child doesn't completely fill the box, where to put the child
type HorizontalJustify int

const (
	Left HorizontalJustify = iota

	// NOTE: in order to see this in effect, you must have
	Center
	Right
)

// When the child's height doesn't completely fill the box, where to put the child
type VerticalJustify int

const (
	Top VerticalJustify = iota
	Middle
	Bottom
)

// TODO make an interface
type Flexbox struct {
	children []flexbox_item.FlexboxItem

	// TODO make configurable on left and right
	padding int

	border lipgloss.Border

	horizontalJustify HorizontalJustify
	verticalJustify   VerticalJustify

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
		horizontalJustify: Left,
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

func (b *Flexbox) SetHorizontalJustify(justify HorizontalJustify) *Flexbox {
	b.horizontalJustify = justify
	return b
}

func (b *Flexbox) SetVerticalJustify(justify VerticalJustify) *Flexbox {
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
	calculationResult := b.calculateChildWidths(widthAvailableForChildren)

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
	case Left:
		pad := strings.Repeat(" ", widthNotUsedByChildren)
		allContentFragments = append(allContentFragments, pad)
	case Right:
		pad := strings.Repeat(" ", widthNotUsedByChildren)
		allContentFragments = append([]string{pad}, allContentFragments...)
	case Center:
		leftPadSize := widthNotUsedByChildren / 2
		rightPadSize := widthNotUsedByChildren - leftPadSize
		leftPad := strings.Repeat(" ", leftPadSize)
		rightPad := strings.Repeat(" ", rightPadSize)

		newContentFragments := append([]string{leftPad}, allContentFragments...)
		newContentFragments = append(newContentFragments, rightPad)
		allContentFragments = newContentFragments
	}

	// TODO allow other align types
	content := lipgloss.JoinHorizontal(lipgloss.Top, allContentFragments...)

	// Justify vertically
	switch b.verticalJustify {
	case Top:
		content += strings.Repeat("\n", heightNotUsedByChildren)
	case Bottom:
		content = strings.Repeat("\n", heightNotUsedByChildren) + content
	case Middle:
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

func (b Flexbox) calculateChildWidths(widthAvailableForChildren int) calculateChildWidthsResult {

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
