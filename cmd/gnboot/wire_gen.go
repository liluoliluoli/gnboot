// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"gnboot/internal/biz"
	"gnboot/internal/conf"
	"gnboot/internal/data"
	"gnboot/internal/pkg/task"
	"gnboot/internal/server"
	"gnboot/internal/service"
)

import (
	_ "github.com/go-cinch/common/plugins/gorm/filter"
	_ "github.com/go-cinch/common/plugins/kratos/encoding/yml"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(c *conf.Bootstrap) (*kratos.App, func(), error) {
	worker, err := task.New(c)
	if err != nil {
		return nil, nil, err
	}
	universalClient, err := data.NewRedis(c)
	if err != nil {
		return nil, nil, err
	}
	tenant, err := data.NewDB(c)
	if err != nil {
		return nil, nil, err
	}
	sonyflake, err := data.NewSonyflake(c)
	if err != nil {
		return nil, nil, err
	}
	tracerProvider, err := data.NewTracer(c)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup := data.NewData(universalClient, tenant, sonyflake, tracerProvider)
	movieRepo := data.NewMovieRepo(dataData)
	transaction := data.NewTransaction(dataData)
	cache := data.NewCache(c, universalClient)
	movieUseCase := biz.NewMovieUseCase(c, movieRepo, transaction, cache)
	gnbootService := service.NewGnbootService(worker, movieUseCase)
	grpcServer := server.NewGRPCServer(c, gnbootService)
	httpServer := server.NewHTTPServer(c, gnbootService)
	app := newApp(grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}