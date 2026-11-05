// Diskless PathFix manager: makes and looks after portable packages, and
// (in its GameFix tab) copies game folders into a main folder for clients.
package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"raptech-portable/internal/gamefix"
)

//go:embed all:frontend/dist
var assets embed.FS

// launcherFS holds assets/launcher.exe, built by build.cmd before the
// manager. A checkout without it still builds; making a package then says
// to run build.cmd.
//
//go:embed all:assets
var launcherFS embed.FS

func main() {
	// Started as the GameFix tab's Windows service, or to install or
	// remove it: no window.
	if gamefix.RunCommandLine() {
		return
	}
	app := NewApp()
	gf := gamefix.NewApp()
	app.gameFixMain = gf.GetMainFolder // hunt never suggests it
	err := wails.Run(&options.App{
		Title:            "Diskless PathFix",
		Width:            1040,
		Height:           700,
		MinWidth:         760,
		MinHeight:        480,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 17, G: 20, B: 24, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			gf.Startup(ctx)
		},
		Bind: []interface{}{app, gf},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
