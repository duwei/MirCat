package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var CFG = Config{}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx
	CFG.load()
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

func (a *App) GetConfig() Config {
	return CFG.GetData()
}

func (a *App) SetConfig(val Config) {
	CFG.SetData(&val)
	runtime.EventsEmit(a.ctx, "config-saved", "OK")
}

var clientArray []*TcpClient

func (a *App) ClientTcpOpen() int {
	tcpClient, err := NewTcpClient(CFG.Client.ServerIp + ":" + CFG.Client.ServerPort)
	if err != nil {
		runtime.EventsEmit(a.ctx, "client-tcp-error", []interface{}{-1, fmt.Sprintf("%v", err)})
		fmt.Printf("Failed to connect: %v\n", err)
		return -1
	}
	if clientArray == nil {
		clientArray = []*TcpClient{}
	}
	clientArray = append(clientArray, tcpClient)
	tcpClient.index = len(clientArray) - 1
	tcpClient.ctx = a.ctx
	runtime.EventsEmit(a.ctx, "client-tcp-info", []interface{}{tcpClient.index, "connection opened"})
	return tcpClient.index
}

func (a *App) ClientTcpSend(index int, base64Data string) {
	if clientArray == nil || index >= len(clientArray) || index < 0 {
		runtime.EventsEmit(a.ctx, "client-tcp-error", []interface{}{index, "invalid client index"})
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		runtime.EventsEmit(a.ctx, "client-tcp-error", []interface{}{index, base64Data + " decode failed"})
		return
	}
	clientArray[index].Send(decodedBytes)
}

func (a *App) ClientTcpClose(index int) {
	if clientArray == nil || index >= len(clientArray) || index < 0 {
		runtime.EventsEmit(a.ctx, "client-tcp-error", []interface{}{index, "invalid client index"})
		return
	}
	clientArray[index].Shutdown()
}

func (a *App) ClientTcpCloseAll() {
	for _, client := range clientArray {
		client.Shutdown()
	}
	clientArray = []*TcpClient{}
}
