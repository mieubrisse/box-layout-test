package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/bubblebath"
	"github.com/mieubrisse/box-layout-test/components/flexbox"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/components/text"
	"os"
)

func main() {
	text1 := flexbox.NewWithContent(text.New("This is text 1")).
		SetBorder(lipgloss.NormalBorder())
	text2 := flexbox.NewWithContent(text.New("This is text 2"), flexbox_item.WithOverflowStyle(flexbox_item.Truncate)).
		SetBorder(lipgloss.DoubleBorder())
	text3 := flexbox.NewWithContent(text.New("This is text 3")).
		SetBorder(lipgloss.BlockBorder()).
		SetHorizontalJustify(flexbox.MainAxisCenter)

	yourBox := flexbox.NewWithContents(
		flexbox_item.NewItem(text1),
		flexbox_item.NewItem(text2),
		flexbox_item.NewItem(text3).SetMinWidth(flexbox_item.FixedSize(20)),
	).SetHorizontalJustify(flexbox.MainAxisCenter).SetVerticalJustify(flexbox.CrossAxisMiddle).SetBorder(lipgloss.NormalBorder())

	if _, err := bubblebath.RunBubbleBathProgram(
		yourBox,
		nil,
		[]tea.ProgramOption{
			tea.WithAltScreen(),
		},
	); err != nil {
		fmt.Printf("An error occurred running the program:\n%v", err)
		os.Exit(1)
	}
}
