package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
	"user/model"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var secret = []byte("hahahahah")

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func md5Password(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	h.Write(secret)
	return hex.EncodeToString(h.Sum(nil))
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	if req.RePassword != req.Password {
		return nil, errors.New("passwords do not match")
	}
	logx.Debugv(req)
	// 查询 username 是否已经存在
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, errors.New("internal error")
	}
	if u != nil {
		return nil, errors.New("username already exists")
	}

	// 加密密码(加盐)
	passwordStr := md5Password(req.Password)
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		UserId:   time.Now().Unix(),
		Username: req.Username,
		Password: passwordStr,
		Gender:   req.Gender,
	})
	if err != nil {
		return nil, err
	}
	return &types.SignupResponse{
		Message: "hello " + req.Username + ", congratulations on signing up.",
	}, nil
}
