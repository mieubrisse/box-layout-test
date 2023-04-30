package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/bubblebath"
	"github.com/mieubrisse/box-layout-test/components/flexbox"
	"github.com/mieubrisse/box-layout-test/components/flexbox_item"
	"github.com/mieubrisse/box-layout-test/components/stylebox"
	"github.com/mieubrisse/box-layout-test/components/text"
	"os"
)

var red = lipgloss.Color("#FF0000")
var blue = lipgloss.Color("#0000FF")
var green = lipgloss.Color("#00FF00")
var lightGray = lipgloss.Color("#333333")

var text1Style = lipgloss.NewStyle().Foreground(red).Background(lightGray)
var text2Style = lipgloss.NewStyle().Foreground(green).Border(lipgloss.NormalBorder())
var text3Style = lipgloss.NewStyle().Foreground(blue).Background(lightGray)

func main() {
	text1 := stylebox.New(text.New("This is text 1")).SetStyle(text1Style)
	text2 := stylebox.New(text.New("This is text 2")).SetStyle(text2Style)
	text3 := stylebox.New(
		text.New("Four score and seven years ago our fathers brought forth on this continent, " +
			"a new nation, conceived in Liberty, and dedicated to the proposition that all men " +
			"are created equal.").
			SetTextAlignment(text.AlignCenter)).SetStyle(text3Style)

	yourBox := flexbox.NewWithContents(
		flexbox_item.New(text1),
		flexbox_item.New(text2),
		flexbox_item.New(text3),
	).SetHorizontalAlignment(flexbox.AlignCenter).
		SetVerticalAlignment(flexbox.AlignCenter).SetDirection(flexbox.Column)

	appBox := stylebox.New(yourBox).SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	if _, err := bubblebath.RunBubbleBathProgram(
		appBox,
		nil,
		[]tea.ProgramOption{
			tea.WithAltScreen(),
		},
	); err != nil {
		fmt.Printf("An error occurred running the program:\n%v", err)
		os.Exit(1)
	}
}
