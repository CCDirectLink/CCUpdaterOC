package main

import (
	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/integration"
	"github.com/CCDirectLink/CCUpdaterOC/middle"
)

// ShowPrimaryView shows the "Primary View" (the mod list right now)
func (app *upApplication) ShowPrimaryView() {
	pkg, err := app.gameInstance.Get("browser")
	installed := err == nil && pkg.Installed()

	topLabelText := "Your selected game instance is: \"" + app.gameInstance.Path + "\""
	topLabel := framework.FlexboxSlot{
		Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(topLabelText, design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
		Grow:    1,
	}

	if installed {
		middleLabelText := "You are ready to go!"
		middleLabel := framework.FlexboxSlot{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(middleLabelText, design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			Grow:    1,
		}

		bottomLabelText := "You can now install mods through your browser."
		bottomLabel := framework.FlexboxSlot{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(bottomLabelText, design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			Grow:    1,
		}

		thePage := design.LayoutDocument(design.Header{
			Title: "Install",
			Back: func() {
				app.GSLeftwards()
				app.ResetWithGameLocation(false, middle.GameFinderVFSPathDefault)
			},
			BackIcon:    design.GameIconID,
			ForwardIcon: design.MenuIconID,
			Forward: func() {
				app.GSRightwards()
				app.ShowOptionsMenu(func() {
					app.GSLeftwards()
					app.ShowPrimaryView()
				})
			},
		}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots: []framework.FlexboxSlot{
				topLabel,
				middleLabel,
				bottomLabel,
			},
		}), true)
		app.Teleport(thePage)
	} else {
		bottomLabelText := "Your need to link your game with this tool."
		bottomLabel := framework.FlexboxSlot{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(bottomLabelText, design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
		}

		button := framework.FlexboxSlot{
			Element: design.ButtonBar([]framework.UILayoutElement{
				design.ButtonAction(design.ThemeOkActionButton, "LINK", func() {
					pkg.Install()
					app.ShowPrimaryView()
				}),
			}),
		}
		thePage := design.LayoutDocument(design.Header{
			Title: "Install",
			Back: func() {
				app.GSLeftwards()
				app.ResetWithGameLocation(false, middle.GameFinderVFSPathDefault)
			},
			BackIcon:    design.GameIconID,
			ForwardIcon: design.MenuIconID,
			Forward: func() {
				app.GSRightwards()
				app.ShowOptionsMenu(func() {
					app.GSLeftwards()
					app.ShowPrimaryView()
				})
			},
		}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots: []framework.FlexboxSlot{
				topLabel,
				bottomLabel,
				button,
			},
		}), true)
		app.Teleport(thePage)
	}
}
