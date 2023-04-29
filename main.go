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
		SetPadding(2).
		SetBorder(lipgloss.NormalBorder())
	text2 := flexbox.NewWithContent(text.New("This is text 2")).
		SetPadding(2).
		SetBorder(lipgloss.DoubleBorder())
	text3 := text.New("This is text 3")

	yourBox := flexbox.NewWithContents(
		flexbox_item.NewItem(text1).
			SetMinWidth(flexbox_item.MaxContentWidth).
			SetMaxWidth(flexbox_item.MaxAvailableWidth),
		flexbox_item.NewItem(text2).
			SetMinWidth(flexbox_item.MinContentWidth).
			SetMinWidth(flexbox_item.MaxAvailableWidth).
			SetOverflowStyle(flexbox_item.Truncate),
		flexbox_item.NewItem(text3).
			SetMinWidth(flexbox_item.MaxContentWidth).
			SetMaxWidth(flexbox_item.MaxAvailableWidth),
	)

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
