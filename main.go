package main

import (
	"context"
	"embed"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	secondInstanceArgs := secondInstanceData.Args

	println("user opened second instance", strings.Join(secondInstanceData.Args, ","))
	println("user opened second from", secondInstanceData.WorkingDirectory)
	runtime.WindowUnminimise(*&a.ctx)
	runtime.Show(*&a.ctx)
	go runtime.EventsEmit(*&a.ctx, "launchArgs", secondInstanceArgs)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	onStartup := func(ctx context.Context) {
		// patcher()
		app.startup(ctx)
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "Backyard Monster Refitted",
		Width:         344,
		Height:        560,
		DisableResize: true,
		Frameless:     true,
		//CSSDragProperty: "widows",
		//CSSDragValue:    "drag",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 255},
		OnStartup:        onStartup,
		Bind: []interface{}{
			app,
		},

		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "e3984e09-28dc-4e3d-b70a-45e961589cdc",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
