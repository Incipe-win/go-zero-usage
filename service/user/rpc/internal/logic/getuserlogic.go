package logic

import (
	"context"
	"errors"

	"user/rpc/internal/svc"
	"user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserID)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user.GetUserResponse{
		UserID:   u.UserId,
		Username: u.Username,
		Gender:   u.Gender,
	}, nil
}
