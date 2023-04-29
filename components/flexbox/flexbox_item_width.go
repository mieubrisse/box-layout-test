package flexbox

type FlexboxItemWidth struct {
	// Given a min and a max, gets the corresponding size based on what FlexboxItemWidth this is
	sizeRetriever func(min, max uint) uint

	// Whether this item width will expand to consume additional free space beyond its min and max
	shouldGrow bool
}

// Indicates a size == the minimum content size of the item, which is the size of the item
// if all wrapping opportunities are taken (basically, the length of the longest word)
var MinContentWidth = FlexboxItemWidth{
	sizeRetriever: func(min, max uint) uint {
		return min
	},
	shouldGrow: false,
}

// Indicates a size == the maximum content of the item, which is the size of the item without any wrapping applied
// Basically, the length of the longest line
var MaxContentWidth = FlexboxItemWidth{
	sizeRetriever: func(min, max uint) uint {
		return min
	},
	shouldGrow: false,
}

// Indicates a size == the maximum amount of space available (including extra space)
var MaxAvailableWidth = FlexboxItemWidth{
	sizeRetriever: func(min, max uint) uint {
		return max
	},
	shouldGrow: true,
}

// Indicates a fixed size
func FixedSizeWidth(size uint) FlexboxItemWidth {
	return FlexboxItemWidth{
		sizeRetriever: func(min, max uint) uint {
			return size
		},
		shouldGrow: false,
	}
}
