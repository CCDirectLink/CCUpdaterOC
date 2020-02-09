package main

import (
	"path/filepath"
	"sort"

	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterOC/middle"
)

func (app *upApplication) ShowGameFinder(back framework.ButtonBehavior, vfsPath string) {
	var vfsList []middle.GameLocation

	vfsList = middle.GameFinderVFSList(vfsPath)
	items := []design.ListItemDetails{}

	for _, v := range vfsList {
		thisLocation := v.Location
		ild := design.ListItemDetails{
			Icon: design.DirectoryIconID,
			Text: filepath.Base(thisLocation),
		}
		ild.Click = func() {
			app.GSRightwards()
			app.ShowGameFinder(func() {
				app.GSInstant()
				app.ShowGameFinder(back, vfsPath)
			}, thisLocation)
		}
		if v.Valid {
			ild.Click = func() {
				app.GSRightwards()
				app.ResetWithGameLocation(true, thisLocation)
			}
			ild.Text = "CrossCode " + v.Version
			ild.Subtext = thisLocation
			ild.Icon = design.GameIconID
		} else if v.Drive != "" {
			ild.Text = v.Drive
			ild.Subtext = v.Location
			ild.Icon = design.DriveIconID
		}
		items = append(items, ild)
	}

	sort.Sort(design.SortListItemDetails(items))
	primary := design.LayoutDocument(design.Header{
		Back:  back,
		Title: "Enter CrossCode's location",
	}, design.NewUISearchBoxPtr("Directory name...", items), true)
	app.GSInstant()
	app.Teleport(primary)
}
