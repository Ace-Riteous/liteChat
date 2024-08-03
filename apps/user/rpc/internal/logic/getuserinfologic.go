package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"liteChat/apps/user/models"
	"liteChat/apps/user/rpc/internal/svc"
	"liteChat/apps/user/rpc/user"
	"liteChat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotFound = xerr.New(xerr.DB_ERROR, "Not Found This User!")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your types here and delete this line
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, errors.WithStack(ErrUserNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "No Such User' s Id Like %s ", in.Id)
	}
	var resp user.UserEntity
	_ = copier.Copy(&resp, userEntity)
	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
