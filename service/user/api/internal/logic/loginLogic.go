package logic

import (
	"context"
	"errors"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 1. 验证参数
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password cannot be empty")
	}

	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err == sqlx.ErrNotFound {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("internal error")
	}
	if user.Password != md5Password(req.Password) {
		return nil, errors.New("invalid username or password")
	}

	// 生成 JWT
	now := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, user.UserId)
	if err != nil {
		return nil, errors.New("internal error")
	}

	// 2. 登录
	return &types.LoginResponse{
		Message:      "Login successful",
		AccessToken:  token,
		AccessExpire: int(now + expire),
		RefreshAfter: int(now + expire/2),
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
