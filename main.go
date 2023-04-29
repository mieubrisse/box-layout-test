package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/bubblebath"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/box"
	"github.com/mieubrisse/box-layout-test/components/text"
	"os"
)

func main() {
	myText := text.New("Hello, world!")

	myBox := box.New(myText)
	myBox.SetBorder(lipgloss.NormalBorder())

	yourBox := box.New(myBox)
	yourBox.SetBorder(lipgloss.DoubleBorder())
	yourBox.SetChildSizeContraint(components.ChildSizeConstraint{
		Min: components.MinContent,
		Max: components.MaxAvailable,
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
