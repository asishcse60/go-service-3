package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
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
		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()
		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, key, &v)
		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
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
