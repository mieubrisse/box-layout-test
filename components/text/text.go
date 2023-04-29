package text

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/wordwrap"
	"strings"
)

// Analogous to the <p> tag in HTML
type Text interface {
	components.Component

	GetContents() string
	SetContents(str string) Text
}

type textImpl struct {
	text string
}

func New(text string) Text {
	return &textImpl{
		text: text,
	}
}

func (t textImpl) GetContents() string {
	return t.text
}

func (t *textImpl) SetContents(str string) Text {
	t.text = str
	return t
}

func (t textImpl) GetContentWidths() (min, max uint) {
	min = 0
	for _, field := range strings.Fields(t.text) {
		printableWidth := uint(ansi.PrintableRuneWidth(field))
		if printableWidth > min {
			min = printableWidth
		}
	}

	max = uint(lipgloss.Width(t.text))

	return
}

func (t textImpl) GetContentHeightGivenWidth(width uint) uint {
	wrappedText := wordwrap.String(t.text, int(width))
	return uint(lipgloss.Height(wrappedText))
}

func (t textImpl) View(width uint, height uint) string {
	return lipgloss.NewStyle().
		Width(int(width)).
		// The only overflow behaviour we can support is truncate
		MaxHeight(int(height)).
		Render(t.text)
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================
