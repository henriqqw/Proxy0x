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

	mDisconnect := systray.AddMenuItem("Disconnect", "Disconnect the active proxy")
	mRotate := systray.AddMenuItem("Rotate", "Rotate the active proxy's IP")
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

	mDisconnect.Click(func() {
		if id := a.activeProxyID(); id != "" {
			a.DisconnectProxy(id)
		}
	})

	mRotate.Click(func() {
		if id := a.activeProxyID(); id != "" {
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

func (a *App) activeProxyID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for id, p := range a.proxies {
		if p.Active {
			return id
		}
	}
	return ""
}
