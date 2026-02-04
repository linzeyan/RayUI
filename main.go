package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"time"

	rayapp "github.com/RayUI/RayUI/internal/app"
	"github.com/RayUI/RayUI/internal/tray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// AppVersion is set via ldflags at build time.
var AppVersion string

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	rayapp.SetAppVersion(AppVersion)
	application := rayapp.NewApp()

	// Set tray icon (will be started after Wails is ready)
	tray.SetIcon(appIcon)

	// Create sub-filesystem for frontend/dist
	// This ensures assets are served from the correct root
	frontendFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to create frontend sub-filesystem: %v", err)
	}

	err = wails.Run(&options.App{
		Title:     "RayUI",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: frontendFS,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		StartHidden:      false, // Ensure window is visible on start
		OnStartup: func(ctx context.Context) {
			application.Startup(ctx)

			// Ensure window is visible after a short delay
			go func() {
				time.Sleep(300 * time.Millisecond)
				wailsruntime.WindowShow(ctx)
				wailsruntime.WindowCenter(ctx)
				// Bring window to front
				wailsruntime.WindowSetAlwaysOnTop(ctx, true)
				time.Sleep(100 * time.Millisecond)
				wailsruntime.WindowSetAlwaysOnTop(ctx, false)
			}()

			// Start system tray (uses external loop, non-blocking)
			tray.Run(tray.Callbacks{
				OnToggleCore: func() {
					application.ToggleCore()
				},
				OnShowWindow: func() {
					wailsruntime.WindowShow(ctx)
					wailsruntime.WindowUnminimise(ctx)
				},
				OnQuit: func() {
					// Set quit flag so OnBeforeClose allows the close
					application.RequestQuit()
					wailsruntime.Quit(ctx)
				},
			})
		},
		OnShutdown: func(ctx context.Context) {
			tray.Quit()
			application.Shutdown(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			// OnBeforeClose is called for window close button click.
			// Only hide to tray if user clicks the close button (X).
			if application.ShouldCloseToTray() {
				wailsruntime.WindowHide(ctx)
				return true
			}
			return false
		},
		Bind: []interface{}{
			application,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
			},
			About: &mac.AboutInfo{
				Title:   "RayUI",
				Message: "A modern cross-platform proxy client",
				Icon:    appIcon,
			},
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
