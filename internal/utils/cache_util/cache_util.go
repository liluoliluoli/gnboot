package cache_util

import (
	"github.com/go-cinch/common/utils"
	"runtime"
)

func GetCacheActionName(condition ...any) string {
	md5 := utils.StructMd5(condition)
	_, file, _, _ := runtime.Caller(1)
	return file + "_" + md5
}
