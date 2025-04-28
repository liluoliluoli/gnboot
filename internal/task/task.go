package task

import (
	"github.com/google/wire"
	"github.com/liluoliluoli/gnboot/internal/task/i4k"
	"github.com/liluoliluoli/gnboot/internal/task/user"
)

var ProviderSet = wire.NewSet(i4k.NewI4kSyncTask, user.NewUserPackageCheckTask)
