package box

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

type Box struct {
	child components.Component

	childSizeConstraint components.ChildSizeConstraint

	overflowStyle OverflowStyle

	// TODO make configurable on left and right
	padding uint

	// TODO give border corners
	border lipgloss.Border

	horizontalJustify HorizontalJustify

	// TODO put a cache that caches the ContentWidths in between the GetContentWidth and View step
}

func New(inner components.Component) Box {
	return Box{
		padding:       0,
		child:         inner,
		overflowStyle: Wrap,
		border:        lipgloss.Border{},
		childSizeConstraint: components.ChildSizeConstraint{
			Min: components.MinContent,
			Max: components.MaxContent,
		},
		horizontalJustify: Left,
	}
}

// TODO make this configurable per side
func (b *Box) SetPadding(padding uint) {
	b.padding = padding
}

func (b *Box) SetBorder(border lipgloss.Border) {
	b.border = border
}

func (b *Box) SetChildSizeContraint(constraint components.ChildSizeConstraint) {
	b.childSizeConstraint = constraint
}

func (b *Box) SetHorizontalJustify(justify HorizontalJustify) {
	b.horizontalJustify = justify
}

func (b Box) GetContentWidths() (min, max uint) {
	additionalNonContentWidth := b.calculateAdditionalNonContentWidth()

	childMin, childMax := b.getChildSizeRange()

	min = childMin + additionalNonContentWidth
	max = childMax + additionalNonContentWidth

	return min, max
}

func (b Box) View(width uint) string {
	// TODO caching of views????

	// If wrap, we'll tell the child about what their real size will be
	// to give them a chance to wrap
	nonContentWidthNeeded := b.calculateAdditionalNonContentWidth()
	spaceAvailableForChild := utilities.GetMaxUint(0, width-nonContentWidthNeeded)

	var widthToGiveChild uint
	switch b.overflowStyle {
	case Wrap:

		childMin, childMax := b.getChildSizeRange()
		if b.childSizeConstraint.Min == components.MaxAvailable {
			childMin = utilities.GetMaxUint(childMin, spaceAvailableForChild)
		}
		if b.childSizeConstraint.Max == components.MaxAvailable {
			childMax = utilities.GetMaxUint(childMax, spaceAvailableForChild)
		}

		widthToGiveChild = utilities.Clamp(width, childMin, childMax)
	case Truncate:
		// If truncating, the child will _think_ they have the full space available
		// and then we'll truncate them later
		// TODO cache this so we don't have to run down the tree again???
		_, innerMax := b.child.GetContentWidths()
		widthToGiveChild = innerMax
	default:
		panic(fmt.Sprintf("Unknown overflow style: %v", b.overflowStyle))
	}

	// Truncate
	truncatedChildStr := truncate.String(b.child.View(widthToGiveChild), spaceAvailableForChild)

	// Now expand, to ensure our box still remains the right size in the case of
	// small strings
	expandedChildStr := padding.String(truncatedChildStr, spaceAvailableForChild)

	// TODO split into left and right pad
	pad := strings.Repeat(" ", int(b.padding))

	result := b.border.Left + pad + expandedChildStr + pad + b.border.Right
	return result
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (b Box) calculateAdditionalNonContentWidth() uint {
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

// Get the possible size ranges for the child, using the child size constraint
// Max is guaranteed to be >= min
func (b Box) getChildSizeRange() (min, max uint) {
	innerMin, innerMax := b.child.GetContentWidths()

	switch b.childSizeConstraint.Min {
	case components.MinContent:
		min = innerMin
	case components.MaxContent, components.MaxAvailable:
		min = innerMax
	default:
		panic(fmt.Sprintf("Unknown minimum child size constraint: %v", b.childSizeConstraint.Min))
	}

	switch b.childSizeConstraint.Max {
	case components.MinContent:
		max = innerMin
	case components.MaxContent, components.MaxAvailable:
		max = innerMax
	default:
		panic(fmt.Sprintf("Unknown maximum child size constraint: %v", b.childSizeConstraint.Max))
	}

	if max < min {
		max = min
	}

	return min, max
}
