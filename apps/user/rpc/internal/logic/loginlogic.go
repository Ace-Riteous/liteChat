package logic

import (
	"context"
	"github.com/pkg/errors"
	"liteChat/apps/user/models"
	"liteChat/pkg/ctxdata"
	"liteChat/pkg/encrypt"
	"liteChat/pkg/xerr"
	"time"

	"liteChat/apps/user/rpc/internal/svc"
	"liteChat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister  = xerr.New(xerr.SERVER_COMMON_ERROR, "This Phone Number has not been Register! ")
	ErrUserPasswordError = xerr.New(xerr.SERVER_COMMON_ERROR, "Password is not Right! ")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your types here and delete this line
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "Find User By Phone Number Err %v , req %v ", err, in.Phone)
	}
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPasswordError)
	}
	now := time.Now().Unix()
	token, err := ctxdata.GenJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "Ctxdata Get Token Err %v", err)

	}

	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
