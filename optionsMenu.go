package main

import (
	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/framework"
)

// ShowOptionsMenu shows the options menu (run game, credits)
func (app *upApplication) ShowOptionsMenu(back framework.ButtonBehavior) {
	backHere := func() {
		app.GSLeftwards()
		app.ShowOptionsMenu(back)
	}
	listSlots := []framework.FlexboxSlot{
		{
			Element: design.ListItem(design.ListItemDetails{
				Text:    "Credits",
				Subtext: "See the various components that make up CCUpdaterOC",
				Click: func() {
					app.GSRightwards()
					app.ShowCredits(backHere)
				},
			}),
		},
		{
			Grow: 1,
		},
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Options",
		Back:  back,
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       listSlots,
	}), true))
}
