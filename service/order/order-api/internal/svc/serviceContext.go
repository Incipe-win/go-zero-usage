package svc

import (
	"order-api/internal/config"
	"order-api/internal/interceptor"
	"order/model"
	"user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   userclient.User
	UserModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySql.DataSource)
	return &ServiceContext{
		Config: c,
		UserRpc: userclient.NewUser(
			zrpc.MustNewClient(
				c.UserRpc,
				zrpc.WithUnaryClientInterceptor(interceptor.HInterceptor),
			),
		),
		UserModel: model.NewOrderModel(conn, c.CacheRedis),
	}
}
