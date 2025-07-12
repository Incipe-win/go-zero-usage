package logic

import (
	"api/internal/svc"
	"api/internal/types"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {

	logx.Debugf("JWT UserID:%v\n", l.ctx.Value("userId"))

	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserID)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("internal error")
	}
	return &types.DetailResponse{
		Username: user.Username,
		Gender:   user.Gender,
	}, nil
}
