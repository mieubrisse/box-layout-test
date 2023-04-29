package flexbox

import (
	"github.com/mieubrisse/box-layout-test/components"
)

// These are simply conveniences for the flexbox.NewWithContent , so that it's super easy to declare a single-item box
type FlexboxItemOpt func(item FlexboxItem)

func WithMinWidth(min FlexboxItemWidth) FlexboxItemOpt {
	return func(item FlexboxItem) {
		item.SetMinWidth(min)
	}
}

func WithMaxWidth(max FlexboxItemWidth) FlexboxItemOpt {
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
	// - MaxAvailableWidth indicates "flex-grow: >1" (see weight below)
	// - Anything else indicates "flex-grow: 0", and sets the "flex-basis"
	minWidth FlexboxItemWidth
	maxWidth FlexboxItemWidth

	overflowStyle OverflowStyle

	// TODO weight (analogous to flex-grow)
	// When the child size constraint is set to MaxAvailableWidth, then this will be used
}

func NewItem(component components.Component) FlexboxItem {
	return &flexboxItemImpl{
		component:     component,
		minWidth:      MinContentWidth,
		maxWidth:      MaxContentWidth,
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
