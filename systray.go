package main

import (
	"log"
	"os/exec"

	icon "github.com/ttimt/Short_Cutteer/icons"
	"github.com/ttimt/systray"
)

var (
	menuLaunchUI *systray.MenuItem
	menuQuit     *systray.MenuItem
)

// Initialize and start the tray icon
func setupTrayIcon() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Short Cutteer")

	// Add default menu items in sequence
	menuLaunchUI = systray.AddMenuItem("Launch UI", "", true)
	systray.AddSeparator()
	menuQuit = systray.AddMenuItem("Quit", "", false)

	go TrayIconListener()
}

// Listen for events from the system tray iconn
func TrayIconListener() {
	for {
		select {
		case <-menuLaunchUI.ClickedCh:
			err := exec.Command("explorer", httpURL).Start()
			if err != nil {
				log.Println(err)
			}

		case <-menuQuit.ClickedCh:
			processInterrupted()

		case <-processInterruptSignal:
			processInterrupted()
		}
	}
}
