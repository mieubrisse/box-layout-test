package components

type Component interface {
	// This is used during the X-expansion phase, where each child "expands" its min and max widths up to its parent
	// During this stage, each element is growing in the X direction; there is no concept of a viewport
	GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight uint)

	// This is used during the Y-expansion phase, where each child "expands" its height up to its parent
	// During this stage, each element is growing in the Y direction; there is no concept of a viewport
	GetContentHeightGivenWidth(width uint) uint

	// TDOO replace result with something that's stylable
	// The actual size that the component ought to have
	View(width uint, height uint) string
}
