package flexbox_item

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
)

type OverflowStyle int

const (
	Wrap OverflowStyle = iota
	Truncate
)

// These are simply conveniences for the flexbox.NewWithContent , so that it's super easy to declare a single-item box
type FlexboxItemOpt func(item FlexboxItem)

func WithMinWidth(min FlexboxItemDimensionValue) FlexboxItemOpt {
	return func(item FlexboxItem) {
		item.SetMinWidth(min)
	}
}

func WithMaxWidth(max FlexboxItemDimensionValue) FlexboxItemOpt {
	return func(item FlexboxItem) {
		item.SetMaxWidth(max)
	}
}

func WithOverflowStyle(style OverflowStyle) FlexboxItemOpt {
	return func(item FlexboxItem) {
		item.SetOverflowStyle(style)
	}
}

type FlexboxItem interface {
	components.Component

	GetComponent() components.Component

	GetMinWidth() FlexboxItemDimensionValue
	SetMinWidth(min FlexboxItemDimensionValue) FlexboxItem
	GetMaxWidth() FlexboxItemDimensionValue
	SetMaxWidth(max FlexboxItemDimensionValue) FlexboxItem

	GetMinHeight() FlexboxItemDimensionValue
	SetMinHeight(min FlexboxItemDimensionValue) FlexboxItem
	GetMaxHeight() FlexboxItemDimensionValue
	SetMaxHeight(max FlexboxItemDimensionValue) FlexboxItem

	GetOverflowStyle() OverflowStyle
	SetOverflowStyle(style OverflowStyle) FlexboxItem
}

type flexboxItemImpl struct {
	component components.Component

	// These determine how the item flexes
	// This is analogous to both "flex-basis" and "flex-grow", where:
	// - MaxAvailableWidth indicates "flex-grow: >1" (see weight below)
	// - Anything else indicates "flex-grow: 0", and sets the "flex-basis"
	minWidth  FlexboxItemDimensionValue
	maxWidth  FlexboxItemDimensionValue
	minHeight FlexboxItemDimensionValue
	maxHeight FlexboxItemDimensionValue

	overflowStyle OverflowStyle

	// Min/maxes of the inner component
	innerDimensionCache components.DimensionsCache

	// TODO weight (analogous to flex-grow)
	// When the child size constraint is set to MaxAvailableWidth, then this will be used
}

func NewItem(component components.Component) FlexboxItem {
	return &flexboxItemImpl{
		component:     component,
		minWidth:      MinContentWidth,
		maxWidth:      MaxContentWidth,
		minHeight:     MinContentWidth,
		maxHeight:     MaxContentWidth,
		overflowStyle: Wrap,
	}
}

func (item *flexboxItemImpl) GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight uint) {
	innerMinWidth, innerMaxWidth, innerMinHeight, innerMaxHeight := item.GetComponent().GetContentMinMax()
	itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight := calculateFlexboxItemContentSizesFromInnerContentSizes(
		innerMinWidth,
		innerMaxWidth,
		innerMinHeight,
		innerMaxHeight,
		item,
	)

	item.innerDimensionCache = components.DimensionsCache{
		MinWidth:  innerMinWidth,
		MaxWidth:  innerMaxWidth,
		MinHeight: innerMinHeight,
		MaxHeight: innerMaxHeight,
	}

	return itemMinWidth, itemMaxWidth, itemMinHeight, itemMaxHeight
}

func (item *flexboxItemImpl) View(width uint, height uint) string {
	component := item.GetComponent()

	var widthWhenRendering uint
	switch item.GetOverflowStyle() {
	case Wrap:
		widthWhenRendering = width
	case Truncate:
		// If truncating, the child will _think_ they have infinite space available
		// and then we'll truncate them later
		widthWhenRendering = item.innerDimensionCache.MaxWidth
	default:
		panic(fmt.Sprintf("Unknown item overflow style: %v", item.GetOverflowStyle()))
	}

	// TODO allow column format
	result := component.View(widthWhenRendering, height)

	// Truncate, in case the inner item rusn over (which will almost definitely be the case when overflowStyle = Truncate)
	result = lipgloss.NewStyle().
		Width(int(width)).
		Height(int(height)).
		MaxWidth(int(width)).
		MaxHeight(int(height)).
		Render(result)

	/*
		// Now expand, to ensure that the item takes up exactly the space we requested
		result = lipgloss.NewStyle().
			Render(result)

	*/

	return result
}

func (item *flexboxItemImpl) GetComponent() components.Component {
	return item.component
}

func (item *flexboxItemImpl) GetMinWidth() FlexboxItemDimensionValue {
	return item.minWidth
}

func (item *flexboxItemImpl) SetMinWidth(min FlexboxItemDimensionValue) FlexboxItem {
	item.minWidth = min
	return item
}

func (item *flexboxItemImpl) GetMaxWidth() FlexboxItemDimensionValue {
	return item.maxWidth
}

func (item *flexboxItemImpl) SetMaxWidth(max FlexboxItemDimensionValue) FlexboxItem {
	item.maxWidth = max
	return item
}

func (item *flexboxItemImpl) GetMinHeight() FlexboxItemDimensionValue {
	return item.minHeight
}

func (item *flexboxItemImpl) SetMinHeight(min FlexboxItemDimensionValue) FlexboxItem {
	item.minHeight = min
	return item
}

func (item *flexboxItemImpl) GetMaxHeight() FlexboxItemDimensionValue {
	return item.maxHeight
}

func (item *flexboxItemImpl) SetMaxHeight(max FlexboxItemDimensionValue) FlexboxItem {
	item.maxHeight = max
	return item
}

func (item *flexboxItemImpl) GetOverflowStyle() OverflowStyle {
	return item.overflowStyle
}

func (item *flexboxItemImpl) SetOverflowStyle(style OverflowStyle) FlexboxItem {
	item.overflowStyle = style
	return item
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================

// Rescales an item's content size based on the per-item configuration the user has set
// Max is guaranteed to be >= min
func calculateFlexboxItemContentSizesFromInnerContentSizes(
	innerMinWidth,
	innertMaxWidth,
	innerMinHeight,
	innerMaxHeight uint,
	item FlexboxItem,
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
