package main

import (
	_ "embed"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/tray.ico
var trayIcon []byte

type langDef struct {
	code string
	name string
}

var supportedLanguages = []langDef{
	{"en", "English"},
	{"pt-BR", "Português (Brasil)"},
	{"pt-PT", "Português (Portugal)"},
	{"es", "Español"},
	{"fr", "Français"},
	{"de", "Deutsch"},
	{"it", "Italiano"},
	{"ru", "Русский"},
	{"ja", "日本語"},
	{"ko", "한국어"},
	{"zh-CN", "中文(简体)"},
	{"ar", "العربية"},
	{"hi", "हिन्दी"},
	{"tr", "Türkçe"},
	{"th", "ไทย"},
	{"pl", "Polski"},
}

func (a *App) setupTray() {
	systray.Run(a.onTrayReady, a.onTrayExit)
}

func (a *App) onTrayReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("Proxy0x")
	systray.SetTooltip("Proxy0x — Proxy Manager")

	mShow := systray.AddMenuItem("Show Proxy0x", "Show the main window")
	systray.AddSeparator()

	mConnectAll := systray.AddMenuItem("Connect All", "Connect all proxies")
	mDisconnectAll := systray.AddMenuItem("Disconnect All", "Disconnect all proxies")
	mRotateAll := systray.AddMenuItem("Rotate All", "Rotate all active proxies")
	systray.AddSeparator()

	mLang := systray.AddMenuItem("Language", "Change language")
	for _, l := range supportedLanguages {
		lang := l
		item := mLang.AddSubMenuItem(lang.name, lang.name)
		item.Click(func() {
			a.SetLanguage(lang.code)
		})
	}

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit Proxy0x")

	mShow.Click(func() {
		runtime.WindowShow(a.ctx)
	})

	mConnectAll.Click(func() {
		a.mu.RLock()
		var offlineIDs []string
		for id, p := range a.proxies {
			if p.Status == StatusOffline {
				offlineIDs = append(offlineIDs, id)
			}
		}
		a.mu.RUnlock()
		for _, id := range offlineIDs {
			go a.ConnectProxy(id)
		}
	})

	mDisconnectAll.Click(func() {
		a.mu.RLock()
		var activeIDs []string
		for id, p := range a.proxies {
			if p.Active {
				activeIDs = append(activeIDs, id)
			}
		}
		a.mu.RUnlock()
		for _, id := range activeIDs {
			a.DisconnectProxy(id)
		}
	})

	mRotateAll.Click(func() {
		a.mu.RLock()
		var activeIDs []string
		for id, p := range a.proxies {
			if p.Active {
				activeIDs = append(activeIDs, id)
			}
		}
		a.mu.RUnlock()
		for _, id := range activeIDs {
			go a.RotateProxy(id)
		}
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
