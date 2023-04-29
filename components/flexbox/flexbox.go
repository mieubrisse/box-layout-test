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
	Center
	Right
)

// TODO make an interface
type Flexbox struct {
	children []flexbox_item.FlexboxItem

	// Cache of the min/max widths/heights across all children
	allChildrenDimensionsCache components.DimensionsCache

	// TODO make configurable on left and right
	padding uint

	border lipgloss.Border

	horizontalJustify HorizontalJustify
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
func (b *Flexbox) SetPadding(padding uint) *Flexbox {
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

func (b *Flexbox) GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight uint) {
	// TODO allow column layout

	var childrenMinWidth, childrenMaxWidth, childrenMinHeight, childrenMaxHeight uint
	for _, item := range b.children {
		itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight := item.GetContentMinMax()

		// Calculate the maxes
		childrenMinWidth = utilities.GetMaxUint(childrenMinWidth, itemMinWidth)
		childrenMaxWidth = utilities.GetMaxUint(childrenMaxWidth, itemMaxWidth)
		childrenMinHeight = utilities.GetMaxUint(childrenMinHeight, itemMinHeight)
		childrenMaxHeight = utilities.GetMaxUint(childrenMaxHeight, itemMaxHeight)

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

func (b Flexbox) View(width uint, height uint) string {
	// TODO caching of views????

	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	// TODO margin
	widthAvailableForChildren := utilities.GetMaxUint(0, width-nonContentWidthNeeded)

	childWidths, totalWidthUsedByChildren := b.calculateChildWidths(widthAvailableForChildren)
	widthNotUsedByChildren := widthAvailableForChildren - totalWidthUsedByChildren

	// TODO allow different types of expansion in cross axis
	additionalNonContentHeight := b.calculateAdditionalNonContentHeight()
	availableChildHeight := utilities.GetMaxUint(0, height-additionalNonContentHeight)

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allContentFragments := make([]string, len(b.children))
	for idx, item := range b.children {
		childWidth := childWidths[idx]
		childStr := item.View(childWidth, availableChildHeight)
		allContentFragments[idx] = childStr
	}

	switch b.horizontalJustify {
	case Left:
		pad := strings.Repeat(" ", int(widthNotUsedByChildren))
		allContentFragments = append(allContentFragments, pad)
	case Right:
		pad := strings.Repeat(" ", int(widthNotUsedByChildren))
		allContentFragments = append([]string{pad}, allContentFragments...)
	case Center:
		leftPadSize := widthNotUsedByChildren / 2
		rightPadSize := widthNotUsedByChildren - leftPadSize
		leftPad := strings.Repeat(" ", int(leftPadSize))
		rightPad := strings.Repeat(" ", int(rightPadSize))

		newContentFragments := append([]string{leftPad}, allContentFragments...)
		newContentFragments = append(newContentFragments, rightPad)
		allContentFragments = newContentFragments
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, allContentFragments...)

	switch b.horizontalJustify {

	}

	result := lipgloss.NewStyle().
		Padding(int(b.padding)).
		Border(b.border).
		Render(content)

	return result
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (b Flexbox) calculateAdditionalNonContentWidth() uint {
	nonContentWidthAdditions := []uint{
		// Padding
		2 * b.padding,

		// Border
		uint(b.border.GetLeftSize() + b.border.GetRightSize()),
	}

	totalNonContentWidthAdded := uint(0)
	for _, addition := range nonContentWidthAdditions {
		totalNonContentWidthAdded += addition
	}
	return totalNonContentWidthAdded
}

func (b Flexbox) calculateAdditionalNonContentHeight() uint {
	nonContentHeightAdditions := []uint{
		// Padding
		2 * b.padding,

		// Border
		uint(b.border.GetTopSize() + b.border.GetBottomSize()),
	}

	totalNonContentHeightAdded := uint(0)
	for _, addition := range nonContentHeightAdditions {
		totalNonContentHeightAdded += addition
	}
	return totalNonContentHeightAdded
}

// Distributes the space
// The only scenario where no space will be distributed is if there is no total weight
// If the space does get distributed, it's guaranteed to be done exactly (no more or less will remain)
func addSpaceByWeight(spaceToAllocate uint, inputSizes []uint, weights []uint) []uint {
	result := make([]uint, len(inputSizes))
	for idx, inputSize := range inputSizes {
		result[idx] = inputSize
	}

	totalWeight := uint(0)
	for _, weight := range weights {
		totalWeight += weight
	}

	// watch out for divide-by-zero
	if totalWeight == 0 {
		return result
	}

	desiredSpaceAllocated := float64(0)
	actualSpaceAllocated := uint(0)
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
		var actualSpaceForItem uint
		if desiredSpaceAllocated < float64(actualSpaceAllocated) {
			// Round up
			actualSpaceForItem = uint(desiredSpaceForItem + 1)
		} else {
			// Round down
			actualSpaceForItem = uint(desiredSpaceForItem)
		}

		result[idx] += actualSpaceForItem
		desiredSpaceAllocated += desiredSpaceForItem
		actualSpaceAllocated += actualSpaceForItem
	}

	return result
}

func (b Flexbox) calculateChildWidths(widthAvailableForChildren uint) ([]uint, uint) {

	numChildren := len(b.children)

	// First, add up the total width the items would like in a perfect world
	totalMinWidthDesired := uint(0) // Under this value, the flexbox will simply truncate
	totalMaxWidthDesired := uint(0) // Above this value, only the MaxAvailableWidth items will expand
	childWidths := make([]uint, numChildren)
	maxConstrainedChildWidths := make([]uint, numChildren)
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
		weights := make([]uint, numChildren)
		for idx, minChildSize := range childWidths {
			maxChildSize := maxConstrainedChildWidths[idx]
			childExpansionRange := maxChildSize - minChildSize
			weights[idx] = childExpansionRange
		}

		spaceToDistributeEvenly := utilities.GetMinUint(
			widthAvailableForChildren-totalMinWidthDesired,
			totalMaxWidthDesired-totalMinWidthDesired,
		)

		childWidths = addSpaceByWeight(spaceToDistributeEvenly, childWidths, weights)
	}

	// When width > max_desired_size, continue to scale only the elements whose max size is MaxAvailableWidth
	if widthAvailableForChildren > totalMaxWidthDesired {
		weights := make([]uint, numChildren)
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
	totalWidthUsed := uint(0)
	widthAvailable := widthAvailableForChildren
	for idx, childWidth := range childWidths {
		actualWidth := utilities.GetMinUint(widthAvailable, childWidth)
		childWidths[idx] = actualWidth

		widthAvailable -= actualWidth
		totalWidthUsed += actualWidth
	}

	return childWidths, totalWidthUsed
}
