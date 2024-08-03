package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"liteChat/apps/user/models"
	"liteChat/apps/user/rpc/internal/config"
	"liteChat/pkg/constants"
	"liteChat/pkg/ctxdata"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	UsersModel models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	systemToken, err := ctxdata.GenJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(), 999_999_999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}

	return svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
