package flexbox_item

// This type represents values for a flexbox item dimension (height or width)
type FlexboxItemDimensionValue struct {
	// Given a min and a max, gets the corresponding size based on what FlexboxItemDimensionValue this is
	sizeRetriever func(min, max uint) uint

	// Whether this item should expand to consume additional free space beyond its min and max
	shouldGrow bool
}

// Returns true if the flexbox item should grow in this dimension if there's space
func (dimensionValue FlexboxItemDimensionValue) ShouldGrow() bool {
	return dimensionValue.shouldGrow
}

// Indicates a size == the minimum content size of the item, which:
// - For width is the size of the item if all wrapping opportunities are taken (basically, the length of the longest word)
// - For height is the height of the item when no word-wrapping is done
var MinContentWidth = FlexboxItemDimensionValue{
	sizeRetriever: func(min, max uint) uint {
		return min
	},
	shouldGrow: false,
}

// Indicates a size == the maximum content of the item, which is the size of the item without any wrapping applied
// - For width, this is basically, the length of the longest line
// - For height, this is the height of the item when the maximum possible word-wrapping is done
var MaxContentWidth = FlexboxItemDimensionValue{
	sizeRetriever: func(min, max uint) uint {
		return max
	},
	shouldGrow: false,
}

// Indicates a size == the maximum amount of space available (including extra space)
var MaxAvailableWidth = FlexboxItemDimensionValue{
	sizeRetriever: func(min, max uint) uint {
		return max
	},
	shouldGrow: true,
}

// Indicates a fixed size
func FixedSize(size uint) FlexboxItemDimensionValue {
	return FlexboxItemDimensionValue{
		sizeRetriever: func(min, max uint) uint {
			return size
		},
		shouldGrow: false,
	}
}
