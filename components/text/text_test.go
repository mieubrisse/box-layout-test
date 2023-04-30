package text

import (
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/stretchr/testify/require"
	"testing"
)

type componentMinMaxSizeAssertion struct {
	minWidth  int
	maxWidth  int
	minHeight int
	maxHeight int
}

func (assertion componentMinMaxSizeAssertion) validate(t *testing.T, component components.Component) {
	minWidth, maxWidth, minHeight, maxHeight := component.GetContentMinMax()
	require.Equal(t, assertion.minWidth, minWidth)
	require.Equal(t, assertion.maxWidth, maxWidth)
	require.Equal(t, assertion.minHeight, minHeight)
	require.Equal(t, assertion.maxHeight, maxHeight)
}

type componentHeightAtWidthAssertion struct {
	width  int
	height int
}

func (assertion componentHeightAtWidthAssertion) validate(t *testing.T, component components.Component) {
	height := component.GetContentHeightForGivenWidth(assertion.width)
	require.Equal(
		t,
		assertion.height,
		height,
	)
}

func TestShortString(t *testing.T) {
	text := New("This is a short string")

	assertion := componentMinMaxSizeAssertion{
		minWidth:  6,
		maxWidth:  22,
		minHeight: 1,
		maxHeight: 4,
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

	assertion := componentMinMaxSizeAssertion{
		minWidth:  6,
		maxWidth:  22,
		minHeight: 1,
		maxHeight: 4,
	}

	// Verify that sizes don't change based off the align
	assertion.validate(t, text)

	text.SetTextAlignment(AlignCenter)
	assertion.validate(t, text)

	text.SetTextAlignment(AlignRight)
	assertion.validate(t, text)
}
