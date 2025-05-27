package task

import (
	"github.com/google/wire"
	"github.com/liluoliluoli/gnboot/internal/task/i4k"
	"github.com/liluoliluoli/gnboot/internal/task/user"
	"github.com/liluoliluoli/gnboot/internal/task/xiaoya/video"  // 新增导入
)

var ProviderSet = wire.NewSet(
	i4k.NewI4kSyncTask, 
	user.NewUserPackageCheckTask,
	video.NewXiaoyaVideoTask,  // 注册新任务
)
