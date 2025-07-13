package logic

import (
	"context"
	"errors"
	"fmt"
	"user/rpc/userclient"

	"order-api/internal/interceptor"
	"order-api/internal/svc"
	"order-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	order, err := l.svcCtx.UserModel.FindOneByOrderId(l.ctx, req.OrderID)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("order id not found")
	}
	if err != nil {
		return nil, err
	}
	fmt.Printf("order: %+v\n", order)
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "33")
	user, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userclient.GetUserRequest{UserID: int64(order.UserId)})
	if err != nil {
		return nil, err
	}
	return &types.SearchResponse{
		OrderID:  order.OrderId,
		Status:   int(order.Status),
		Username: user.GetUsername(),
	}, nil
}
