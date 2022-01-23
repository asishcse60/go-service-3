// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	v1ProductGrp "github.com/asishcse60/service/app/services/sales-api/handlers/v1/productgrp"
	v1TestGrp "github.com/asishcse60/service/app/services/sales-api/handlers/v1/testgrp"
	v1UserGrp "github.com/asishcse60/service/app/services/sales-api/handlers/v1/usergrp"
	"github.com/asishcse60/service/business/core/product"
	"github.com/asishcse60/service/business/core/user"
	"github.com/asishcse60/service/business/sys/auth"
	"github.com/asishcse60/service/business/web/v1/mid"
	"github.com/asishcse60/service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(auth.RoleAdmin)

	// test endpoints.
	tgh := v1TestGrp.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize("ADMIN"))

	// Register user management and authentication endpoints.
	ugh := v1UserGrp.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
		Auth: cfg.Auth,
	}

	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen, admin)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, authen, admin)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, authen, admin)

	// Register product and sale endpoints.
	pgh := v1ProductGrp.Handlers{
		Product: product.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/products/:page/:rows", pgh.Query, authen)
	app.Handle(http.MethodGet, version, "/products/:id", pgh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/products", pgh.Create, authen)
	app.Handle(http.MethodPut, version, "/products/:id", pgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/products/:id", pgh.Delete, authen)
}