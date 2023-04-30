package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/box-layout-test/app_components/favorite_thing"
	"github.com/mieubrisse/box-layout-test/app_components/favorite_things_list"
	"github.com/mieubrisse/box-layout-test/components"
	"github.com/mieubrisse/box-layout-test/components/stylebox"
)

type App interface {
	components.Component
}

type appImpl struct {
	root components.Component
}

func New() App {
	myFavoriteThings := []favorite_thing.FavoriteThing{
		favorite_thing.New().SetName("Pourover coffee").SetDescription("It takes so long to make though"),
		favorite_thing.New().SetName("Pizza").SetDescription("Pepperoni is the best"),
		favorite_thing.New().SetName("Jiu jitsu").SetDescription("Rolling all day"),
	}

	var root components.Component = favorite_things_list.New().SetThings(myFavoriteThings)

	root = stylebox.New(root).SetStyle(lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()))
	return &appImpl{
		root: root,
	}
}

func (a appImpl) GetContentMinMax() (minWidth, maxWidth, minHeight, maxHeight int) {
	return a.root.GetContentMinMax()
}

func (a appImpl) GetContentHeightForGivenWidth(width int) int {
	return a.root.GetContentHeightForGivenWidth(width)
}

func (a appImpl) View(width int, height int) string {
	return a.root.View(width, height)
}
