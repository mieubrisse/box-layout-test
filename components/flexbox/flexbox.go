package flexbox

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/utilities"
	"strings"
)

type OverflowStyle int

const (
	Wrap OverflowStyle = iota
	Truncate
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
	children []FlexboxItem

	// TODO cache content widths so we don't have to burn a bunch of energy recalculating them!

	// TODO make configurable on left and right
	padding uint

	border lipgloss.Border

	horizontalJustify HorizontalJustify
}

// Convenience constructor for a box with a single element
func NewWithContent(component components.Component, opts ...FlexboxItemOpt) *Flexbox {
	item := NewItem(component)
	for _, opt := range opts {
		opt(item)
	}
	return NewWithContents(item)
}

// Convenience constructor for a box with multiple elements
func NewWithContents(items ...FlexboxItem) *Flexbox {
	return New().SetChildren(items)
}

func New() *Flexbox {
	return &Flexbox{
		padding:           0,
		children:          make([]FlexboxItem, 0),
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

func (b *Flexbox) SetChildren(children []FlexboxItem) *Flexbox {
	b.children = children
	return b
}

func (b *Flexbox) SetHorizontalJustify(justify HorizontalJustify) *Flexbox {
	b.horizontalJustify = justify
	return b
}

func (b Flexbox) GetContentWidths() (min, max uint) {
	additionalNonContentWidth := b.calculateAdditionalNonContentWidth()

	var allChildrenMin, allChildrenMax uint
	for _, item := range b.children {
		contentMin, contentMax := item.GetComponent().GetContentWidths()
		childMin, childMax := constrainItemContentSizes(contentMin, contentMax, item)
		allChildrenMin = utilities.GetMaxUint(allChildrenMin, childMin)
		allChildrenMax = utilities.GetMaxUint(allChildrenMax, childMax)
	}

	min = allChildrenMin + additionalNonContentWidth
	max = allChildrenMax + additionalNonContentWidth

	return min, max
}

// TODO set child constraints

func (b Flexbox) View(width uint) string {
	// TODO caching of views????

	// If wrap, we'll tell the child about what their real size will be
	// to give them a chance to wrap
	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	// TODO margin
	spaceAvailableForChildren := utilities.GetMaxUint(0, width-nonContentWidthNeeded)

	numChildren := len(b.children)

	// First, add up the total size the items would like
	totalMinSizeDesired := uint(0) // Under this value, the flexbox will simply truncate
	totalMaxSizeDesired := uint(0) // Above this value, only the MaxAvailableWidth items will expand
	childSizes := make([]uint, numChildren)
	maxConstrainedChildSizes := make([]uint, numChildren)
	contentMaxes := make([]uint, numChildren)
	for idx, item := range b.children {
		contentMin, contentMax := item.GetComponent().GetContentWidths()
		contrainedMin, contrainedMax := constrainItemContentSizes(contentMin, contentMax, item)
		totalMinSizeDesired += contrainedMin
		totalMaxSizeDesired += contrainedMax

		contentMaxes[idx] = contentMax
		childSizes[idx] = contrainedMin
		maxConstrainedChildSizes[idx] = contrainedMax
	}

	// When min_desired_size < space_available < max_desired_size, scale everyone up equally between their
	// min and max sizes
	if spaceAvailableForChildren > totalMinSizeDesired {
		weights := make([]uint, numChildren)
		for idx, minChildSize := range childSizes {
			maxChildSize := maxConstrainedChildSizes[idx]
			childExpansionRange := maxChildSize - minChildSize
			weights[idx] = childExpansionRange
		}

		spaceToDistributeEvenly := utilities.GetMinUint(
			spaceAvailableForChildren-totalMinSizeDesired,
			totalMaxSizeDesired-totalMinSizeDesired,
		)

		childSizes = addSpaceByWeight(spaceToDistributeEvenly, childSizes, weights)
	}

	// When width > max_desired_size, continue to scale only the elements whose max size is MaxAvailableWidth
	if spaceAvailableForChildren > totalMaxSizeDesired {
		weights := make([]uint, numChildren)
		for idx, item := range b.children {
			if item.GetMaxWidth().shouldGrow {
				// TODO use actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		spaceToDistributeToExpanders := spaceAvailableForChildren - totalMaxSizeDesired

		childSizes = addSpaceByWeight(spaceToDistributeToExpanders, childSizes, weights)
	}

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allChildStrs := make([]string, numChildren)
	for idx, item := range b.children {
		component := item.GetComponent()

		childWidth := childSizes[idx]

		var widthWhenRendering uint
		switch item.GetOverflowStyle() {
		case Wrap:
			widthWhenRendering = childWidth
		case Truncate:
			// If truncating, the child will _think_ they have infinite space available
			// and then we'll truncate them later
			widthWhenRendering = contentMaxes[idx]
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

	// TODO margin
	content := strings.Join(allChildStrs, "")

	// TODO split into left and right pad
	pad := strings.Repeat(" ", int(b.padding))

	result := b.border.Left + pad + content + pad + b.border.Right
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

// Max is guaranteed to be >= min
func constrainItemContentSizes(contentMin, contentMax uint, item FlexboxItem) (constrainedMin, constrainedMax uint) {
	constrainedMin = item.GetMinWidth().sizeRetriever(contentMin, contentMax)
	constrainedMax = item.GetMaxWidth().sizeRetriever(contentMin, contentMax)

	if constrainedMax < constrainedMin {
		constrainedMax = constrainedMin
	}

	return constrainedMin, constrainedMax
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
