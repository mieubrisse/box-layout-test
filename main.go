package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/box-layout-test/bubblebath"
	"github.com/mieubrisse/box-layout-test/components/flexbox"
	"github.com/mieubrisse/box-layout-test/components/text"
	"os"
)

func main() {
	text1 := text.New("This is text 1")
	text2 := text.New("This is text 2")
	text3 := text.New("This is text 3")

	yourBox := flexbox.New()
	yourBox.SetChildren([]flexbox.FlexboxItem{
		flexbox.NewItem(text1).
			SetMinWidth(flexbox.MaxContentWidth).
			SetMaxWidth(flexbox.MaxAvailableWidth),
		flexbox.NewItem(text2).
			SetMinWidth(flexbox.MinContentWidth).
			SetMaxWidth(flexbox.FixedSizeWidth(7)).
			SetOverflowStyle(flexbox.Truncate),
		flexbox.NewItem(text3).
			SetMinWidth(flexbox.MaxContentWidth).
			SetMaxWidth(flexbox.MaxAvailableWidth),
	})

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
