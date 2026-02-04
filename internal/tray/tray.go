package tray

import (
	"fyne.io/systray"
)

// Callbacks holds the function pointers the tray menu items call.
type Callbacks struct {
	OnToggleCore func()
	OnShowWindow func()
	OnQuit       func()
}

// icon is the tray icon data (PNG).
var icon []byte

// callbacks stores the current callbacks for use by onReady.
var callbacks Callbacks

// SetIcon sets the tray icon bytes.
func SetIcon(data []byte) {
	icon = data
}

// Run starts the system tray in a background goroutine.
func Run(cb Callbacks) {
	callbacks = cb
	go systray.Run(onReady, onExit)
}

// Quit stops the system tray.
func Quit() {
	systray.Quit()
}

func onReady() {
	if len(icon) > 0 {
		systray.SetIcon(icon)
	}
	// Don't set title - just show the icon in menu bar
	systray.SetTooltip("RayUI - Proxy Client")

	mToggle := systray.AddMenuItem("Toggle Core", "Start or stop the proxy core")
	mShow := systray.AddMenuItem("Dashboard", "Open RayUI window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit RayUI")

	// Use channel-based event handling
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if callbacks.OnToggleCore != nil {
					callbacks.OnToggleCore()
				}
			case <-mShow.ClickedCh:
				if callbacks.OnShowWindow != nil {
					callbacks.OnShowWindow()
				}
			case <-mQuit.ClickedCh:
				// Call the quit callback which will trigger wailsruntime.Quit
				// This will eventually call tray.Quit() from OnShutdown
				if callbacks.OnQuit != nil {
					callbacks.OnQuit()
				}
				return
			}
		}
	}()
}

func onExit() {
	// Cleanup if needed
}
