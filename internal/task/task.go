package task

import "github.com/google/wire"

var TaskProviderSet = wire.NewSet(NewI4kSyncTask)
