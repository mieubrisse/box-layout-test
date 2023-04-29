package flexbox

import (
	"github.com/mieubrisse/box-layout-test/components"
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
