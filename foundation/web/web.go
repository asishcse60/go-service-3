package web

import (
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// App - App is the entry point of our application.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp - creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App{
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
}

// SignalShutdown - used for gracefully shutdown the app.
func (a *App) SignalShutdown(){
	a.shutdown <- syscall.SIGTERM
}