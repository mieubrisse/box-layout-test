package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/utilities"
	"strings"
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

	direction Direction

	mainAxisAlignment  MainAxisAlignment
	crossAxisAlignment CrossAxisAlignment

	// -------------------- Calculation Caching -----------------------
	// The widths each child desires (cached between GetContentMinMax and GetContentHeightForGivenWidth
	desiredChildWidthsCache []int

	// The actual widths each child will get (cached between GetContentHeightForGivenWidth and View)
	actualChildWidthsCache axisCalculationResults

	// The desired height each child wants given its width (cached between GetContentHeightForGivenWidth and View)
	desiredChildHeightsGivenWidthCache []int
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
		children:          make([]flexbox_item.FlexboxItem, 0),
		mainAxisAlignment: MainAxisStart,
	}
}

func (b *Flexbox) SetChildren(children []flexbox_item.FlexboxItem) *Flexbox {
	b.children = children
	return b
}

func (b *Flexbox) SetMainAxesAlignment(alignment MainAxisAlignment) *Flexbox {
	b.mainAxisAlignment = alignment
	return b
}

func (b *Flexbox) SetCrossAxisAlignment(alignment CrossAxisAlignment) *Flexbox {
	b.crossAxisAlignment = alignment
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
	// TODO cache this result!!!!

	// Width
	actualWidthsCalcResults := b.direction.getActualWidths(b.desiredChildWidthsCache, b.children)

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
	// Width
	widthNotUsedByChildren := width - b.actualChildWidthsCache.totalWidthUsed

	// Height
	childHeights, maxHeightUsedByChildren := b.calculateCrossAxisHeights(
		b.actualChildWidthsCache.childWidths,
		height,
	)
	heightNotUsedByChildren := height - maxHeightUsedByChildren

	// Now render each child, ensuring we expand the child's string if the resulting string is less
	allContentFragments := make([]string, len(b.children))
	for idx, item := range b.children {
		childWidth := b.actualChildWidthsCache.childWidths[idx]
		childHeight := childHeights[idx]
		childStr := item.View(childWidth, childHeight)

		allContentFragments[idx] = childStr
	}

	// Justify main axis
	switch b.mainAxisAlignment {
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

	// Justify cross axis
	var content string
	switch b.crossAxisAlignment {
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

	return content
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
