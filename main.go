package main

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"mir-cat/pkg/mircat"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	// 创建一个App结构体实例
	app := mircat.NewApp()
	cfg := mircat.NewConfig()
	connManager := mircat.NewConnManager(app, cfg)

	// Create application with options
	// 使用选项创建应用
	err := wails.Run(&options.App{
		Title:     "MirCat",
		Width:     900,
		Height:    600,
		MinWidth:  900,
		MinHeight: 600,
		//MaxWidth:          1200,
		//MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.Startup,
		OnDomReady:        app.DomReady,
		OnBeforeClose:     app.BeforeClose,
		OnShutdown:        app.Shutdown,
		WindowStartState:  options.Normal,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    nil,
			Middleware: nil,
		},
		Bind: []interface{}{
			cfg, connManager,
		},
		// Windows platform specific options
		// Windows平台特定选项
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
		// Mac platform specific options
		// Mac平台特定选项
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Wails Template React",
				Message: "A Wails template based on React",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
