package main

import (
	"flag"
	"net/url"
	"os"
	"strings"

	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard"
	"github.com/CCDirectLink/CCUpdaterOC/frenyard/framework"
	"github.com/CCDirectLink/CCUpdaterOC/middle"
	"github.com/CCDirectLink/ccmu/game"
	"github.com/CCDirectLink/ccmu/pkg"
)

type upApplication struct {
	gameInstance     game.Game
	gameSelected     bool
	config           middle.UpdaterConfig
	mainContainer    *framework.UISlideTransitionContainer
	window           frenyard.Window
	upQueued         chan func()
	teleportSettings framework.SlideTransition
}

const upTeleportLen float64 = 0.25

// GSLeftwards sets the teleportation affinity to LEFT.
func (app *upApplication) GSLeftwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = false
	app.teleportSettings.Length = upTeleportLen
}

// GSRightwards sets the teleportation affinity to RIGHT.
func (app *upApplication) GSRightwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = false
	app.teleportSettings.Length = upTeleportLen
}

// GSLeftwards sets the teleportation affinity to UP.
func (app *upApplication) GSUpwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = true
	app.teleportSettings.Length = upTeleportLen
}

// GSRightwards sets the teleportation affinity to DOWN.
func (app *upApplication) GSDownwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = true
	app.teleportSettings.Length = upTeleportLen
}

// GSInstant sets the teleportation affinity to INSTANT.
func (app *upApplication) GSInstant() {
	// direction doesn't matter
	app.teleportSettings.Length = 0
}

// Teleport starts a transition with the cached affinity settings.
func (app *upApplication) Teleport(target framework.UILayoutElement) {
	forkTD := app.teleportSettings
	forkTD.Element = target
	app.mainContainer.TransitionTo(forkTD)
}

func main() {
	frenyard.TargetFrameTime = 0.016
	slideContainer := framework.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindowInit)
	wnd, err := framework.CreateBoundWindow("CCUpdaterOC", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}
	design.Setup(frenyard.InferScale(wnd))
	wnd.SetSize(design.SizeWindow)
	// Ok, now get it ready.
	app := (&upApplication{
		config:           middle.ReadUpdaterConfig(),
		mainContainer:    slideContainer,
		window:           wnd,
		upQueued:         make(chan func(), 16),
		teleportSettings: framework.SlideTransition{},
	})

	pkg := app.parseArgs()

	if pkg != nil {
		app.ShowPackageView(func() {}, func() {}, pkg)
	} else {
		app.ShowGameFinderPreface()
	}

	// Started!
	frenyard.GlobalBackend.Run(func(frameTime float64) {
		select {
		case fn := <-app.upQueued:
			fn()
		default:
		}
	})
}

func (app *upApplication) parseArgs() pkg.Package {
	if len(os.Args) <= 1 {
		return nil
	}

	flag.String("game", "", "if set it overrides the path of the game")
	uri := flag.String("url", "", "the url that executed ccmu")
	flag.Parse()

	if uri == nil || *uri == "" {
		return nil
	}

	raw := *uri
	arg := strings.Split(raw[7:len(raw)-1], "/")[0]
	arg, err := url.PathUnescape(arg)
	if err != nil {
		return nil
	}

	app.gameInstance = game.At(flag.Lookup("game").Value.String())
	app.gameSelected = true

	result := app.gameInstance.Find(arg)
	if len(result) == 0 {
		return nil
	}

	return result[0]
}
