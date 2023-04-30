package text

import (
	"github.com/mieubrisse/box-layout-test/components/test_assertions"
	"testing"
)

func TestShortString(t *testing.T) {
	str := "This is a short string"

	sizeAssertions := test_assertions.FlattenAssertionGroups(
		test_assertions.GetDefaultAssertions(),
		test_assertions.GetContentSizeAssertions(6, 22, 1, 4),
		test_assertions.GetHeightAtWidthAssertions(
			0, 0, // invisible
			6, 4, // min content width
			8, 3, // in the middle
			22, 1, // max content width
			100, 1, // beyond max content width
		),
	)

	// Verify that the size assertions are valid at all alignments
	for _, alignment := range []TextAlignment{AlignLeft, AlignCenter, AlignRight} {
		component := New(str).SetTextAlignment(alignment)
		test_assertions.CheckAll(t, sizeAssertions, component)
	}
}
