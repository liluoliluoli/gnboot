//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"gnboot/internal/biz"
	"gnboot/internal/conf"
	"gnboot/internal/data"
	"gnboot/internal/pkg/task"
	"gnboot/internal/server"
	"gnboot/internal/service"
)

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, task.ProviderSet, service.ProviderSet, newApp))
}
