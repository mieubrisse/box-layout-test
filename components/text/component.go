package text

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
	"strings"
)

type Text struct {
	text string
}

// TODO Private
func New(text string) Text {
	return Text{
		text: text,
	}
}

func (t Text) GetContentWidths() (min, max uint) {
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

func (t Text) View(width uint) string {
	return lipgloss.NewStyle().Width(int(width)).Render(t.text)
}
