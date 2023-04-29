package flexbox

import (
	"github.com/mieubrisse/box-layout-test/components"
)

// Prebuilt constants for the child's size
type FlexboxItemWidth int

const (
	// Indicates a width that's equal to the minimum width of the child box's content
	MinContent FlexboxItemWidth = iota

	// Indicates a width that's equal to the maximum width of the child box's content
	MaxContent

	// Indicates a width that's equal to the maximum available space for the child
	// Behaves like MxContent, except that if there's more space available at render time than MaxContent,
	// the child will be given extra space
	MaxAvailable

	// TODO add ways to set fixed widths!!
)

type FlexboxItem interface {
	GetComponent() components.Component

	GetMinWidth() FlexboxItemWidth
	SetMinWidth(min FlexboxItemWidth) FlexboxItem

	GetMaxWidth() FlexboxItemWidth
	SetMaxWidth(max FlexboxItemWidth) FlexboxItem

	GetOverflowStyle() OverflowStyle
	SetOverflowStyle(style OverflowStyle) FlexboxItem
}

type flexboxItemImpl struct {
	component components.Component

	// These determine how the item flexes
	// This is analogous to both "flex-basis" and "flex-grow", where:
	// - MaxAvailable indicates "flex-grow: >1" (see weight below)
	// - Anything else indicates "flex-grow: 0", and sets the "flex-basis"
	minWidth FlexboxItemWidth
	maxWidth FlexboxItemWidth

	overflowStyle OverflowStyle

	// TODO weight (analogous to flex-grow)
	// When the child size constraint is set to MaxAvailable, then this will be used
}

func NewItem(component components.Component) FlexboxItem {
	return &flexboxItemImpl{
		component:     component,
		minWidth:      MinContent,
		maxWidth:      MaxContent,
		overflowStyle: Wrap,
	}
}

func (item *flexboxItemImpl) GetComponent() components.Component {
	return item.component
}

func (item *flexboxItemImpl) GetMinWidth() FlexboxItemWidth {
	return item.minWidth
}

func (item *flexboxItemImpl) SetMinWidth(min FlexboxItemWidth) FlexboxItem {
	item.minWidth = min
	return item
}

func (item *flexboxItemImpl) GetMaxWidth() FlexboxItemWidth {
	return item.maxWidth
}

func (item *flexboxItemImpl) SetMaxWidth(max FlexboxItemWidth) FlexboxItem {
	item.maxWidth = max
	return item
}

func (item *flexboxItemImpl) GetOverflowStyle() OverflowStyle {
	return item.overflowStyle
}

func (item *flexboxItemImpl) SetOverflowStyle(style OverflowStyle) FlexboxItem {
	item.overflowStyle = style
	return item
}
