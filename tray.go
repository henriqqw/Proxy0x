package main

import (
	_ "embed"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

func (a *App) setupTray() {
	systray.Run(a.onTrayReady, a.onTrayExit)
}

func (a *App) onTrayReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("ProxyUI")
	systray.SetTooltip("ProxyUI — Proxy Manager")

	mShow := systray.AddMenuItem("Show ProxyUI", "Show the main window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit ProxyUI")

	mShow.Click(func() {
		runtime.WindowShow(a.ctx)
	})

	mQuit.Click(func() {
		systray.Quit()
		runtime.Quit(a.ctx)
	})

	systray.SetOnClick(func(menu systray.IMenu) {
		runtime.WindowShow(a.ctx)
	})

	systray.SetOnDClick(func(menu systray.IMenu) {
		runtime.WindowShow(a.ctx)
	})
}

func (a *App) onTrayExit() {}
