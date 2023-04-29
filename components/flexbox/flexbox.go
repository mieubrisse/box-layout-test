package flexbox

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/utilities"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
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

type FlexboxItem struct {
	component components.Component

	// The constraints determine how the item flexes
	// This is analogous to both "flex-basis" and "flex-grow", where:
	// - MaxAvailable indicates "flex-grow: >1" (see weight below)
	// - Anything else indicates "flex-grow: 0", and sets the "flex-basis"
	constraint components.ChildSizeConstraint

	overflowStyle OverflowStyle

	// TODO weight (analogous to flex-grow)
	// When the child size constraint is set to MaxAvailable, then this will be used
}

// TODO constructor

type Flexbox struct {
	children []FlexboxItem

	// TODO make configurable on left and right
	padding uint

	// TODO give border corners
	border lipgloss.Border

	horizontalJustify HorizontalJustify

	// TODO put a cache that caches the ContentWidths in between the GetContentWidth and View step
}

func New(inner components.Component) Flexbox {
	return Flexbox{
		padding:           0,
		children:          make([]FlexboxItem, 0),
		border:            lipgloss.Border{},
		horizontalJustify: Left,
	}
}

// TODO make this configurable per side
func (b *Flexbox) SetPadding(padding uint) {
	b.padding = padding
}

func (b *Flexbox) SetBorder(border lipgloss.Border) {
	b.border = border
}

func (b *Flexbox) SetChildren(children []FlexboxItem) {
	b.children = children
}

func (b *Flexbox) SetHorizontalJustify(justify HorizontalJustify) {
	b.horizontalJustify = justify
}

func (b Flexbox) GetContentWidths() (min, max uint) {
	additionalNonContentWidth := b.calculateAdditionalNonContentWidth()

	var allChildrenMin, allChildrenMax uint
	for _, item := range b.children {
		childMin, childMax := getChildSizeRangeUsingConstraints(item)
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
	minChildSizeDesired := uint(0) // Under this value, the flexbox will simply truncate
	maxChildSizeDesired := uint(0) // Above this value, only the MaxAvailable items will expand
	childSizes := make([]uint, numChildren)
	maxChildSizes := make([]uint, numChildren)
	for idx, item := range b.children {
		childMin, childMax := getChildSizeRangeUsingConstraints(item)
		minChildSizeDesired += childMin
		maxChildSizeDesired += childMax

		childSizes[idx] = childMin
		maxChildSizes[idx] = childMax
	}

	// When min_desired_size < space_available < max_desired_size, scale everyone up equally between their
	// min and max sizes
	if spaceAvailableForChildren > minChildSizeDesired && spaceAvailableForChildren < maxChildSizeDesired {
		weights := make([]uint, numChildren)
		for idx := range b.children {
			// TODO use actual weights
			weights[idx] = 1
		}

		spaceToDistributeEvenly := maxChildSizeDesired - spaceAvailableForChildren

		childSizes = addSpaceByWeight(spaceToDistributeEvenly, childSizes, weights)
	}

	// When width > max_desired_size, continue to scale only the elements whose max size is MaxAvailable
	if spaceAvailableForChildren > maxChildSizeDesired {
		weights := make([]uint, numChildren)
		for idx, item := range b.children {
			if item.constraint.Max == components.MaxAvailable {
				// TODO use actual weights
				weights[idx] = 1
				continue
			}

			weights[idx] = 0
		}

		spaceToDistributeToExpanders := maxChildSizeDesired - spaceAvailableForChildren

		childSizes = addSpaceByWeight(spaceToDistributeToExpanders, childSizes, weights)
	}

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allChildStrs := make([]string, numChildren)
	for idx, item := range b.children {
		component := item.component

		childWidth := childSizes[idx]

		var widthWhenRendering uint
		switch item.overflowStyle {
		case Wrap:
			widthWhenRendering = childWidth
		case Truncate:
			// If truncating, the child will _think_ they have the full space available
			// and then we'll truncate them later
			widthWhenRendering = maxChildSizes[idx]
		default:
			panic(fmt.Sprintf("Unknown overflow style: %v", item.overflowStyle))
		}

		childStr := component.View(widthWhenRendering)

		// Truncate, in case any children are over
		childStr = truncate.String(childStr, childWidth)

		// Now expand, to ensure that children with MaxAvailable get expanded
		childStr = padding.String(childStr, childWidth)

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

// Get the possible size ranges for the child, using the child size constraints
// Max is guaranteed to be >= min
func getChildSizeRangeUsingConstraints(item FlexboxItem) (min, max uint) {
	innerMin, innerMax := item.component.GetContentWidths()

	switch item.constraint.Min {
	case components.MinContent:
		min = innerMin
	case components.MaxContent, components.MaxAvailable:
		min = innerMax
	default:
		panic(fmt.Sprintf("Unknown minimum component size constraint: %v", item.constraint.Min))
	}

	switch item.constraint.Max {
	case components.MinContent:
		max = innerMin
	case components.MaxContent, components.MaxAvailable:
		max = innerMax
	default:
		panic(fmt.Sprintf("Unknown maximum component size constraint: %v", item.constraint.Max))
	}

	if max < min {
		max = min
	}

	return min, max
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
