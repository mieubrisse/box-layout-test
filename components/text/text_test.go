package text

import (
	"github.com/mieubrisse/box-layout-test/components/test_assertions"
	"testing"
)

func TestShortString(t *testing.T) {
	text := New("This is a short string")

	var ints []int
	var subints []int
	var subints2 []int
	ints = append(ints, subints..., subints2...)

	var assertions []test_assertions.ComponentAssertion
	assertions = append(assertions, test_assertions.GetContentSizeAssertions(6, 22, 1, 4)
	assertions = append(assertions, test_assertions.GetHeightAtWidthAssertions(
		0, 0, // invisible
		6, 4, // min content width
		8, 2, // in the middle
		22, 1, // max content width
		100, 1, // beyond max content width
	))
	assertions = append()

	[]test_assertions.ComponentAssertion{
		test_assertions.ContentSizeAssertion{
			ExpectedMinWidth:  6,
			ExpectedMaxWidth:  22,
			ExpectedMinHeight: 1,
			ExpectedMaxHeight: 4,
		},
		// Invisible
		test_assertions.HeightAtWidthAssertion{},
		// Min content width
		test_assertions.HeightAtWidthAssertion{
			Width:          6,
			ExpectedHeight: 4,
		},
		// In the middle
		test_assertions.HeightAtWidthAssertion{
			Width:          22,
			ExpectedHeight: 1,
		},
		// Max content width
		test_assertions.HeightAtWidthAssertion{
			Width:          22,
			ExpectedHeight: 1,
		},
		// Beyond max
		test_assertions.HeightAtWidthAssertion{
			Width:          100,
			ExpectedHeight: 1,
		},
	}

	// Verify that sizes don't change based off the align
	assertion.validate(t, text)

	text.SetTextAlignment(AlignCenter)
	assertion.validate(t, text)

	text.SetTextAlignment(AlignRight)
	assertion.validate(t, text)

	ass

}

func TestShortString(t *testing.T) {
	text := New("This is a short string")

	assertion := ComponentMinMaxSizeAssertion{
		MinWidth:  6,
		MaxWidth:  22,
		MinHeight: 1,
		MaxHeight: 4,
	}

	// Verify that sizes don't change based off the align
	assertion.validate(t, text)

	text.SetTextAlignment(AlignCenter)
	assertion.validate(t, text)

	text.SetTextAlignment(AlignRight)
	assertion.validate(t, text)
}
