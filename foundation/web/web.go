package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// Handler - A handler is a type that handles an http request within our own little mini.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App - App is the entry point of our application.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp - creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// SignalShutdown - used for gracefully shutdown the app.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method, group, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// PRE CODE PROCCESSING.
		if err := handler(r.Context(), w, r); err != nil {

		}
		// POST CODE PROCCESSING.
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	// original call to library.
	a.ContextMux.Handle(method, finalPath, h)
}
