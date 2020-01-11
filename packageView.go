package main

import (
	"github.com/CCDirectLink/CCUpdaterUI/design"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard/integration"
	"github.com/CCDirectLink/CCUpdaterUI/middle"
	"github.com/CCDirectLink/ccmu/pkg"
)

// ShowPackageView shows a dialog for a package.
// backHeavy is used if nothing actually happened.
// backLight is used if something happened (the ShowPackageView return call has both set to the same value).
// This allows preserving state in the PrimaryView.
func (app *upApplication) ShowPackageView(backHeavy framework.ButtonBehavior, backLight framework.ButtonBehavior, pkg pkg.Package) {
	info, _ := pkg.Info()
	// No latest package = no information.
	if pkg == nil {
		if middle.InternetConnectionWarning {
			app.MessageBox("Package not available", "The package '"+info.NiceName+"' could not be found.\n\nAs you have ended up here, the package probably had to exist in some form.\nThis error is probably because CCUpdaterUI was unable to retrieve remote packages.\n\n1. Check your internet connection\n2. Try restarting CCUpdaterUI\n3. Contact us", backLight)
		} else {
			app.MessageBox("I just don't know what went wrong...", "The package '"+info.NiceName+"' could not be found.\nYou should never be able to see this dialog in normal operation.", backLight)
		}
		return
	}

	// Ok, now let's actually start work on the UI
	showInstallButton := false
	annotations := "\n    ID: " + info.NiceName + "\n    Latest Version: " + info.CurrentVersion

	if pkg.Installed() {
		outdated, _ := info.Outdated()
		if outdated {
			annotations += "\n    Installed: " + info.CurrentVersion
			showInstallButton = true
		} else {
			annotations += "\n    Installed"
		}
	} else {
		showInstallButton = true
	}

	chunks := []integration.TypeChunk{
		integration.NewColouredTextTypeChunk(info.NiceName, design.GlobalFont, design.ThemeText),
		integration.NewColouredTextTypeChunk(annotations, design.ListItemSubTextFont, design.ThemeSubText),
	}
	buttons := []framework.UILayoutElement{}
	if pkg.Installed() && info.Name != "Simplify" {
		//TODO
		/*
			removeTx := ccmodupdater.PackageTX{
				pkg: ccmodupdater.PackageTXOperationRemove,
			}
			_, removeErr := txCtx.Solve(removeTx)
			if removeErr != nil {
				buttonText = "NOT REMOVABLE"
				removeTheme = design.ThemeImpossibleActionButton
			}
		*/
		buttons = append(buttons, design.ButtonAction(design.ThemeRemoveActionButton, "REMOVE", func() {
			app.GSDownwards()
			app.PerformTransaction(func() {
				app.GSUpwards()
				app.ShowPackageView(backHeavy, backHeavy, pkg)
			}, transaction{pkg: opUninstall})
		}))
	}
	if showInstallButton {
		//TODO
		/*
			installTx := ccmodupdater.PackageTX{
				pkg: ccmodupdater.PackageTXOperationInstall,
			}
			_, removeErr := txCtx.Solve(installTx)
			if removeErr != nil {
				buttonText = "NOT INSTALLABLE"
				buttonColour = design.ThemeImpossibleActionButton
			}
		*/
		buttonText := "INSTALL"
		buttonColour := design.ThemeOkActionButton
		buttonTx := transaction{pkg: opInstall}
		outdated, _ := info.Outdated()
		if outdated {
			buttonText = "UPDATE"
			buttonColour = design.ThemeUpdateActionButton
			buttonTx = transaction{pkg: opUpdate}
		}

		buttons = append(buttons, design.ButtonAction(buttonColour, buttonText, func() {
			app.GSDownwards()
			app.PerformTransaction(func() {
				app.GSUpwards()
				app.ShowPackageView(backHeavy, backHeavy, pkg)
			}, buttonTx)
		}))
	}

	detail := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: design.NewIconPtr(0xFFFFFFFF, middle.PackageIcon(pkg), 36),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewCompoundTypeChunk(chunks), 0xFFFFFFFF, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart}),
			},
		},
	})

	fullPanel := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: detail,
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(info.Description, design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart}),
				Shrink:  1,
			},
			framework.FlexboxSlot{
				Grow:   1,
				Shrink: 1,
			},
			framework.FlexboxSlot{
				Element: design.ButtonBar(buttons),
			},
		},
	})

	app.Teleport(design.LayoutDocument(design.Header{
		Title: info.NiceName,
		Back:  backLight,
	}, fullPanel, true))
}
