package task

import (
	"github.com/google/wire"
	"github.com/liluoliluoli/gnboot/internal/task/user"
	"github.com/liluoliluoli/gnboot/internal/task/video"
)

var ProviderSet = wire.NewSet(
	user.NewUserPackageCheckTask,
	video.NewJfVideoTask,
)
