package flexbox

import "github.com/mieubrisse/box-layout-test/components"

type FlexboxItem struct {
	component components.Component

	// The constraints determine how the item flexes
	// This is analogous to both "flex-basis" and "flex-grow", where:
	// - MaxAvailable indicates "flex-grow: >1" (see weight below)
	// - Anything else indicates "flex-grow: 0", and sets the "flex-basis"
	constraint *ChildSizeConstraint

	overflowStyle OverflowStyle

	// TODO weight (analogous to flex-grow)
	// When the child size constraint is set to MaxAvailable, then this will be used
}

func NewItem(component components.Component) *FlexboxItem {
	return &FlexboxItem{
		component:     component,
		constraint:    NewConstraint(),
		overflowStyle: Wrap,
	}
}

func (item *FlexboxItem) SetConstraint(constraint *ChildSizeConstraint) *FlexboxItem {
	item.constraint = constraint
	return item
}

func (item *FlexboxItem) SetOverflowStyle(style OverflowStyle) *FlexboxItem {
	item.overflowStyle = style
	return item
}
