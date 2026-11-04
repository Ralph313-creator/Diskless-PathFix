// RapTech Portable manager: makes and looks after portable packages.
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
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
	app := NewApp()
	err := wails.Run(&options.App{
		Title:            "RapTech Portable",
		Width:            980,
		Height:           680,
		MinWidth:         720,
		MinHeight:        480,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 17, G: 20, B: 24, A: 1},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
