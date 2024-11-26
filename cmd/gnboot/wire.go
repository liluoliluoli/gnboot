//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"github.com/liluoliluoli/gnboot/internal/adaptor"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/server"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/task"
)

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet, task.ProviderSet, adaptor.ProviderSet, service.ProviderSet, repo.ProviderSet, newApp))
}
