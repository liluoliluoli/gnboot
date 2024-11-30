package cache_util

import (
	"github.com/go-cinch/common/utils"
	"runtime"
)

func GetCacheActionName(condition ...any) string {
	md5 := utils.StructMd5(condition)
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return "n/a"
	}
	method := runtime.FuncForPC(pc[0]).Name()
	return method + "_" + md5
}
