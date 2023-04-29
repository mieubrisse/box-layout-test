package components

// Simple struct for caching the min/max width/height of anything
type DimensionsCache struct {
	MinWidth  uint
	MaxWidth  uint
	MinHeight uint
	MaxHeight uint
}
