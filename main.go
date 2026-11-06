package main

import (
	"embed"
	"log"
	"os"

	"proxy_ui/internal/config"
	"proxy_ui/internal/killswitch"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Run by the uninstaller so firewall rules never outlive the app.
	if len(os.Args) > 1 && os.Args[1] == "--disable-killswitch" {
		if err := killswitch.New(config.KillSwitchBackupFile()).Disable(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Proxy0x",
		Width:     420,
		Height:    680,
		MinWidth:  380,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Frameless:         true,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			Theme:                windows.Dark,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
