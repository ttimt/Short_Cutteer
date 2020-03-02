// +build windows

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ttimt/systray"

	_ "github.com/HouzuoGuo/tiedot/db"
	_ "github.com/lxn/walk"

	. "github.com/ttimt/Short_Cutteer/hook/windows"
)

const (
	windowsNewLine = "\r\n"
)

var (
	processInterruptSignal = make(chan os.Signal)
)

func init() {
	// Setup process interrupt signal
	signal.Notify(processInterruptSignal, os.Interrupt)

	// Initialize templates
	initializeTemplates()

	// Initialize DB
	initializeDB()

	// Create all hook keys
	createAllHookKeys()
}

func main() {
	// Call systray GUI
	systray.Run(onReady, nil)
}

func onReady() {
	// Run the server
	setupHTTPServer()

	// Start low level keyboard listener
	setupWindowsHook()

	// Setup system tray icon
	setupTrayIcon()
}

// When Ctrl+C or interrupt signal received
func processInterrupted() {
	// Unhook Windows keyboard
	log.Println("Removing Windows hook ......")
	_, _ = UnhookWindowsHookEx(hhook)

	// Quit system tray
	log.Println("Removing sytem tray ......")
	systray.Quit()

	// Close db
	err := myDB.Close()
	if err != nil {
		panic(err)
	}

	// Exit
	os.Exit(1)
}
