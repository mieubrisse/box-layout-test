package components

type Component interface {
	// During this step (expansion), each child recursively "expands" its min and max sizes up to the parent
	// During this stage, each element is "growing"; there is no concept of a viewport
	GetContentWidths() (min, max uint)

	// During this step (compression), we now know the viewport size and each parent is trying to "compress" each
	// child to fit into the constraints
	// TDOO replace result with something that's stylable
	// The actual size that the component ought to have
	View(width uint) string
}
