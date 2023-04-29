package components

// Prebuilt constants for the child's size
type ChildWidthControl int

const (
	// Indicates a width that's equal to the minimum width of the child box's content
	MinContent ChildWidthControl = iota

	// Indicates a width that's equal to the maximum width of the child box's content
	MaxContent

	// Indicates a width that's equal to the maximum available space for the child
	// Behaves like MxContent, except that if there's more space available at render time than MaxContent,
	// the child will be given extra space
	MaxAvailable

	// TODO add ways to set fixed widths!!
)

// Controls how the parent sizes a child, including the minimum and maximum that the parent will size
// the child by
// If min > max, then the max will be defined by the min
type ChildSizeConstraint struct {
	// The width under which a child won't be asked to reflow (it'll just be truncated)
	Min ChildWidthControl

	// The width over which a child cannot expand
	Max ChildWidthControl
}
