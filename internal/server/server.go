package server

import (
	"embed"
	"github.com/google/wire"
	task2 "gnboot/internal/adaptor/task"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, task2.NewWorker)

//go:embed middleware/locales
var locales embed.FS
