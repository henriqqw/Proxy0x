package main

// appVersion is set at build time: -ldflags "-X main.appVersion=1.2.3".
var appVersion = "dev"

func (a *App) GetVersion() string { return appVersion }
