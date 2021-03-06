package main

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/integration"
	"github.com/CCDirectLink/CCUpdaterOC/middle"
	"github.com/CCDirectLink/ccmu/game"
)

func (app *upApplication) ResetWithGameLocation(save bool, location string) {
	app.gameInstance = game.At(location)
	app.gameSelected = false
	app.config.GamePath = location
	if save {
		middle.WriteUpdaterConfig(app.config)
	}
	// Re-kick
	app.ShowGameFinderPreface()
}

func (app *upApplication) ShowGameFinderPreface() {
	var gameLocations []middle.GameLocation
	app.ShowWaiter("Starting...", func(progress func(string)) {
		progress("Scanning local installation...")
		gi := game.At(app.config.GamePath)
		fmt.Printf("Doing preliminary check of %s\n", app.config.GamePath)
		_, err := gi.BasePath()
		if err == nil {
			app.gameInstance = gi
			app.gameSelected = true
			return
		}
		fmt.Printf("Game not present?\n")
		progress("Not configured ; Autodetecting game locations...")
	}, func() {
		if !app.gameSelected {
			app.ShowGameFinderPrefaceInternal(gameLocations)
		} else {
			app.ShowPrimaryView()
		}
	})
}

func (app *upApplication) ShowGameFinderPrefaceInternal(locations []middle.GameLocation) {

	suggestSlots := []framework.FlexboxSlot{}
	for _, location := range locations {
		suggestSlots = append(suggestSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Icon:    design.GameIconID,
				Text:    "CrossCode " + location.Version,
				Subtext: location.Location,
				Click: func() {
					app.GSRightwards()
					app.ResetWithGameLocation(true, location.Location)
				},
			}),
			RespectMinimumSize: true,
		})
	}
	// Space-taker to prevent wrongly scaled list items
	suggestSlots = append(suggestSlots, framework.FlexboxSlot{
		Grow:   1,
		Shrink: 0,
	})

	foundInstallsScroller := design.ScrollboxV(framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		WrapMode:    framework.FlexboxWrapModeNone,
		Slots:       suggestSlots,
	}))

	content := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Welcome to the official 1-Click CrossCode Mod Installer.\nBefore we begin, you need to indicate which CrossCode installation you want to install mods to.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element:            foundInstallsScroller,
				Grow:               1,
				Shrink:             1,
				RespectMinimumSize: true,
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("If the installation you'd like to install mods to isn't shown here, you can locate it manually.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: design.ButtonBar([]framework.UILayoutElement{
					design.ButtonAction(design.ThemeOkActionButton, "LOCATE MANUALLY", func() {
						app.GSDownwards()
						app.ShowGameFinder(func() {
							app.GSUpwards()
							app.ShowGameFinderPrefaceInternal(locations)
						}, middle.GameFinderVFSPathDefault)
					}),
				}),
			},
		},
	})
	primary := design.LayoutDocument(design.Header{
		BackIcon: design.WarningIconID,
		Back: func() {
			app.GSLeftwards()
			app.ShowCredits(func() {
				app.GSRightwards()
				app.ShowGameFinderPrefaceInternal(locations)
			})
		},
		Title: "Welcome",
	}, content, true)
	app.Teleport(primary)
}
