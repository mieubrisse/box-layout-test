package flexbox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/components/stylebox"
	"github.com/mieubrisse/box-layout-test/components/test_assertions"
	"github.com/mieubrisse/box-layout-test/components/text"
	"testing"
)

func TestColumnLayout(t *testing.T) {
	child1 := text.New("This is child 1")
	child2 := stylebox.New(text.New("This is child 2")).
		SetStyle(lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()))
	child3 := text.New("This is child 3")

	flexbox := NewWithContents(
		flexbox_item.New(child1),
		flexbox_item.New(child2),
		flexbox_item.New(child3),
	).SetHorizontalAlignment(AlignCenter).SetVerticalAlignment(AlignCenter)

	width, height := 30, 30

	assertions := test_assertions.FlattenAssertionGroups(
		test_assertions.GetDefaultAssertions(),
		test_assertions.GetContentSizeAssertions(
			7,
			17,
			5,
			20,
		),
	)

	// Need to populate the caches
	flexbox.GetContentMinMax()
	flexbox.GetContentHeightForGivenWidth(width)
	flexbox.View(width, height)

	test_assertions.CheckAll(t, assertions, flexbox)
}
