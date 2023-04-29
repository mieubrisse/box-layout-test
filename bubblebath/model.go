package bubblebath

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/box"
)

type BubbleBathOption func(*bubbleBathModel)

func WithInitCmd(cmd tea.Cmd) BubbleBathOption {
	return func(model *bubbleBathModel) {
		model.initCmd = cmd
	}
}

func WithQuitSequences(quitSequenceSet map[string]bool) BubbleBathOption {
	return func(model *bubbleBathModel) {
		model.quitSequenceSet = quitSequenceSet
	}
}

var defaultQuitSequenceSet = map[string]bool{
	"ctrl+c": true,
	"ctrl+d": true,
}

type bubbleBathModel struct {
	// The tea.Cmd that will be fired upon initialization
	initCmd tea.Cmd

	// Sequences matching String() of tea.KeyMsg that will quit the program
	quitSequenceSet map[string]bool

	// We put the user's app in a box here so we can give the user top-level control over how their app expands/contracts
	// relative to the terminal
	appBox components.Component

	app components.Component

	width  uint
	height uint
}

// NewBubbleBathModel creates a new tea.Model for tea.NewProgram based off the given InteractiveComponent
func NewBubbleBathModel(app components.Component, options ...BubbleBathOption) tea.Model {
	appBox := box.New(app)
	appBox.SetChildSizeContraint(components.ChildSizeConstraint{
		Min: components.MinContent,
		Max: components.MaxAvailable,
	})
	result := &bubbleBathModel{
		initCmd:         nil,
		quitSequenceSet: defaultQuitSequenceSet,
		appBox:          appBox,
		app:             app,
		width:           0,
		height:          0,
	}
	for _, opt := range options {
		opt(result)
	}
	return result
}

func (b bubbleBathModel) Init() tea.Cmd {
	return b.initCmd
}

func (b bubbleBathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if _, found := b.quitSequenceSet[msg.String()]; found {
			return b, tea.Quit

		}
	case tea.WindowSizeMsg:
		// b.appComponent.Resize(msg.Width, msg.Height)
		b.width = uint(msg.Width)
		b.height = uint(msg.Height)
		return b, nil
	}

	// return b, b.appComponent.Update(msg)
	return b, nil
}

func (b bubbleBathModel) View() string {
	return b.appBox.View(b.width)
}

/*
func (b bubbleBathModel) GetAppComponent() InteractiveComponent {
	return b.appComponent
}
*/

func RunBubbleBathProgram[T components.Component](
	appComponent T,
	bubbleBathOptions []BubbleBathOption,
	teaOptions []tea.ProgramOption,
) (T, error) {
	model := NewBubbleBathModel(appComponent, bubbleBathOptions...)

	finalModel, err := tea.NewProgram(model, teaOptions...).Run()
	castedModel := finalModel.(bubbleBathModel)
	castedAppComponent := castedModel.app.(T)
	return castedAppComponent, err
}
