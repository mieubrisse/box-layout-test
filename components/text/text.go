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

func (t *textImpl) GetContentMinMax() (minWidth int, maxWidth int, minHeight int, maxHeight int) {
	minWidth = 0
	for _, field := range strings.Fields(t.text) {
		printableWidth := ansi.PrintableRuneWidth(field)
		if printableWidth > minWidth {
			minWidth = printableWidth
		}
	}

	maxWidth = lipgloss.Width(t.text)

	minHeight = lipgloss.Height(t.text)

	minWidthWrapped := wordwrap.String(t.text, maxWidth)
	maxHeight = lipgloss.Height(minWidthWrapped)

	return
}

func (t textImpl) View(width int, height int) string {
	return lipgloss.NewStyle().
		Width(width).
		// The only overflow behaviour we can support is truncate
		MaxHeight(height).
		Render(t.text)
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================
