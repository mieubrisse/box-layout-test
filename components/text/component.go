package text

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/muesli/reflow/ansi"
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
	max = uint(lipgloss.Width(t.text))

	min = 0
	for _, field := range strings.Fields(t.text) {
		printableWidth := uint(ansi.PrintableRuneWidth(field))
		if printableWidth > min {
			min = printableWidth
		}
	}

	return
}

func (t textImpl) View(width uint) string {
	return lipgloss.NewStyle().Width(int(width)).Render(t.text)
}
