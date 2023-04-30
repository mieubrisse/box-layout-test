package stylebox

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/utilities"
)

// Stylebox is a box explicitly for controlling style
// No other elements control style
type Stylebox interface {
	components.Component

	GetStyle() lipgloss.Style
	// NOTE: all layout-affecting properties (height, width, alignment, margin, inline) are ignored
	// The only layout-affecting property left in place are border and padding
	SetStyle(style lipgloss.Style) Stylebox
}

type styleboxImpl struct {
	component components.Component

	style lipgloss.Style
}

func New(component components.Component) Stylebox {
	return &styleboxImpl{
		component: component,
		style:     lipgloss.NewStyle(),
	}
}

func (s styleboxImpl) GetStyle() lipgloss.Style {
	return s.style
}

func (s styleboxImpl) SetStyle(style lipgloss.Style) Stylebox {
	s.style = style.Copy().
		UnsetMargins().
		UnsetAlign().
		UnsetAlignHorizontal().
		UnsetAlignVertical().
		UnsetWidth().
		UnsetMaxWidth().
		UnsetHeight().
		UnsetMaxHeight().
		UnsetInline()
	return s
}

func (s styleboxImpl) GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight int) {
	// TODO cache the results?
	innerMinWidth, innerMaxWidth, innerMinHeight, innerMaxHeight := s.component.GetContentMinMax()

	minWidth = innerMinWidth + s.style.GetHorizontalFrameSize()
	maxWidth = innerMaxWidth + s.style.GetHorizontalFrameSize()

	minHeight = innerMinHeight + s.style.GetVerticalFrameSize()
	maxHeight = innerMaxHeight + s.style.GetVerticalFrameSize()
	return
}

func (s styleboxImpl) GetContentHeightForGivenWidth(width int) int {
	innerWidth := utilities.GetMaxInt(0, width-s.style.GetHorizontalFrameSize())
	return s.component.GetContentHeightForGivenWidth(innerWidth)
}

func (s styleboxImpl) View(width int, height int) string {
	if width == 0 || height == 0 {
		return ""
	}

	innerWidth := utilities.GetMaxInt(0, width-s.style.GetHorizontalFrameSize())
	innerHeight := utilities.GetMaxInt(0, height-s.style.GetVerticalFrameSize())
	innerStr := s.component.View(innerWidth, innerHeight)

	result := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		MaxWidth(innerWidth).
		MaxHeight(innerHeight).
		Render(innerStr)

	result = s.style.Render(result)

	return result
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (s styleboxImpl) getExtraWidth() int {
	return s.style.GetHorizontalFrameSize()
}

func (s styleboxImpl) getExtraHeight() int {
	return s.style.GetVerticalFrameSize()
}
