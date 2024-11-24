//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"gnboot/internal/adaptor"
	"gnboot/internal/conf"
	"gnboot/internal/repo"
	"gnboot/internal/server"
	"gnboot/internal/service"
	"gnboot/internal/task"
)

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, repo.ProviderSet, service.ProviderSet, task.ProviderSet, adaptor.ProviderSet, newApp))
}
