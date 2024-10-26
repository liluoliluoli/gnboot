package service

import (
	"context"
	"fmt"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/query"
	"gnboot/pkg/global"
)

type MysqlService struct {
	Q query.MySql
}

func New(ctx context.Context) MysqlService {
	ops := []func(*query.MysqlOptions){
		query.WithMysqlCtx(ctx),
		query.WithMysqlDb(global.Mysql),
		query.WithMysqlCasbinEnforcer(global.CasbinEnforcer),
		query.WithMysqlCachePrefix(fmt.Sprintf("%s_%s", global.Conf.Mysql.DSN.DBName, constant.QueryCachePrefix)),
	}
	if global.Conf.Redis.Enable {
		ops = append(ops, query.WithMysqlRedis(global.Redis))
	}
	my := MysqlService{
		Q: query.NewMySql(ops...),
	}
	return my
}
