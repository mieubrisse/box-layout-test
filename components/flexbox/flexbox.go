package flexbox

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
	"strings"
)

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

	// Cache containing the results of calculating min/max width/height for each child via GetContentMinMax
	childrenMinMaxDimensionsCache []itemMinMaxDimensionsCache

	// TODO cache content widths so we don't have to burn a bunch of energy recalculating them!

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

func (b Flexbox) GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight uint) {

	// TODO allow column layout

	var childrenMinWidth, childrenMaxWidth, childrenMinHeight, childrenMaxHeight uint
	newCache := make([]itemMinMaxDimensionsCache, len(b.children))
	for idx, item := range b.children {
		innerMinWidth, innerMaxWidth, innerMinHeight, innerMaxHeight := item.GetComponent().GetContentMinMax()
		itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight := calculateFlexboxItemContentSizesFromInnerContentSizes(
			innerMinWidth,
			innerMaxWidth,
			innerMinHeight,
			innerMaxHeight,
			item,
		)

		// Calculate the maxes
		childrenMinWidth = utilities.GetMaxUint(childrenMinWidth, itemMinWidth)
		childrenMaxWidth = utilities.GetMaxUint(childrenMaxWidth, itemMaxWidth)
		childrenMinHeight = utilities.GetMaxUint(childrenMinHeight, itemMinHeight)
		childrenMaxHeight = utilities.GetMaxUint(childrenMaxHeight, itemMaxHeight)

		newCache[idx] = itemMinMaxDimensionsCache{
			minWidth:  itemMinWidth,
			maxWidth:  itemMaxWidth,
			minHeight: itemMinHeight,
			maxHeight: itemMaxHeight,
		}
	}

	additionalNonContentWidth := b.calculateAdditionalNonContentWidth()
	minWidth = childrenMinWidth + additionalNonContentWidth
	maxWidth = childrenMaxWidth + additionalNonContentWidth

	additionalNonContentHeight := b.calculateAdditionalNonContentHeight()
	minHeight = childrenMinHeight + additionalNonContentHeight
	maxHeight = childrenMaxHeight + additionalNonContentHeight

	return
}

func (b *Flexbox) GetContentHeightGivenWidth(width uint) uint {
	//TODO implement me
	panic("implement me")
}

func (b Flexbox) View(width uint, height int) string {
	// TODO caching of views????

	childWidths := b.calculateChildWidths(width)

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allChildStrs := make([]string, numChildren)
	for idx, item := range b.children {
		component := item.GetComponent()

		childWidth := childWidths[idx]

		var widthWhenRendering uint
		switch item.GetOverflowStyle() {
		case flexbox_item.Wrap:
			widthWhenRendering = childWidth
		case flexbox_item.Truncate:
			// If truncating, the child will _think_ they have infinite space available
			// and then we'll truncate them later
			widthWhenRendering = contentWidthMaxes[idx]
		default:
			panic(fmt.Sprintf("Unknown item overflow style: %v", item.GetOverflowStyle()))
		}

		childStr := component.View(widthWhenRendering)

		// Truncate, in case any children are over
		childStr = lipgloss.NewStyle().
			MaxWidth(int(childWidth)).
			MaxHeight(1).
			Render(childStr)

		// Now expand, to ensure that children with MaxAvailableWidth get expanded
		padNeeded := int(childWidth) - lipgloss.Width(childStr)
		childStr += strings.Repeat(" ", padNeeded)

		allChildStrs[idx] = childStr
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, allChildStrs...)

	result := lipgloss.NewStyle().
		Padding(int(b.padding)).
		Border(b.border).
		Render(content)

	pad := strings.Repeat(" ", int(b.padding))

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

// Rescales an item's content size based on the per-item configuration the user has set
// Max is guaranteed to be >= min
func calculateFlexboxItemContentSizesFromInnerContentSizes(
	innerMinWidth,
	innertMaxWidth,
	innerMinHeight,
	innerMaxHeight uint,
	item flexbox_item.FlexboxItem,
) (itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight uint) {
	itemMinWidth = item.GetMinWidth().sizeRetriever(innerMinWidth, innertMaxWidth)
	itemMaxWidth = item.GetMaxWidth().sizeRetriever(innerMinWidth, innertMaxWidth)

	if itemMaxWidth < itemMinWidth {
		itemMaxWidth = itemMinWidth
	}

	itemMinHeight = item.GetMinHeight().sizeRetriever(innerMinHeight, innerMaxHeight)
	itemMaxHeight = item.GetMaxHeight().sizeRetriever(innerMinHeight, innerMaxHeight)

	if itemMaxHeight < itemMinHeight {
		itemMaxHeight = itemMinHeight
	}

	return
}

// Distributes the space
// The given space is guaranteed to be exactly distributed (no more or less will remain)
func addSpaceByWeight(spaceToAllocate uint, inputSizes []uint, weights []uint) []uint {
	totalWeight := uint(0)
	for _, weight := range weights {
		totalWeight += weight
	}

	result := make([]uint, len(inputSizes))

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

func (b Flexbox) calculateChildWidths(flexboxWidth uint) []uint {
	// If wrap, we'll tell the child about what their real size will be
	// to give them a chance to wrap
	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	// TODO margin
	widthAvailableForChildren := utilities.GetMaxUint(0, flexboxWidth-nonContentWidthNeeded)

	numChildren := len(b.children)

	// First, add up the total size the items would like
	totalMinWidthDesired := uint(0) // Under this value, the flexbox will simply truncate
	totalMaxWidthDesired := uint(0) // Above this value, only the MaxAvailableWidth items will expand
	childWidths := make([]uint, numChildren)
	maxConstrainedChildWidths := make([]uint, numChildren)
	contentWidthMaxes := make([]uint, numChildren)
	for idx, item := range b.children {
		contentMin, contentMax := item.GetComponent().GetContentMinMax()
		contrainedMin, contrainedMax := calculateFlexboxItemContentSizesFromInnerContentSizes(contentMin, contentMax, item)
		totalMinWidthDesired += contrainedMin
		totalMaxWidthDesired += contrainedMax

		contentWidthMaxes[idx] = contentMax
		childWidths[idx] = contrainedMin
		maxConstrainedChildWidths[idx] = contrainedMax
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
			if item.GetMaxWidth().shouldGrow {
				// TODO use actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		spaceToDistributeToExpanders := widthAvailableForChildren - totalMaxWidthDesired

		childWidths = addSpaceByWeight(spaceToDistributeToExpanders, childWidths, weights)
	}

	return childWidths
}
