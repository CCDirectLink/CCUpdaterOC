package main

import (
	"github.com/CCDirectLink/CCUpdaterUI/design"
	"github.com/CCDirectLink/CCUpdaterUI/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterUI/middle"
	"path/filepath"
	"sort"
)

func (app *upApplication) ShowGameFinder(back framework.ButtonBehavior, vfsPath string) {
	var vfsList []middle.GameLocation
	
	app.ShowWaiter("Reading", func (progress func(string)) {
		progress("Scanning to find all of the contents of in:\n" + vfsPath + "\nIf this includes CD/DVD drives or network partitions, this may take a while.")
		vfsList = middle.GameFinderVFSList(vfsPath)
	}, func () {
		items := []design.ListItemDetails{}
		
		for _, v := range vfsList {
			thisLocation := v.Location
			ild := design.ListItemDetails{
				Icon: design.DirectoryIconID,
				Text: filepath.Base(thisLocation),
			}
			ild.Click = func () {
				app.GSRightwards()
				app.ShowGameFinder(func () {
					app.GSLeftwards()
					app.ShowGameFinder(back, vfsPath)
				}, thisLocation)
			}
			if v.Valid {
				ild.Click = func () {
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
			Back: back,
			Title: "Enter CrossCode's location",
		}, design.NewUISearchBoxPtr("Directory name...", items), true)
		app.Teleport(primary)
	})
}
