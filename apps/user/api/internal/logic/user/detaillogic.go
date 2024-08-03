package user

import (
	"context"
	"github.com/jinzhu/copier"
	"liteChat/apps/user/rpc/user"
	"liteChat/pkg/ctxdata"

	"liteChat/apps/user/api/internal/svc"
	"liteChat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your types here and delete this line
	uid := ctxdata.GetUId(l.ctx)
	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}
	var res types.User
	_ = copier.Copy(&res, userInfoResp.User)
	return &types.UserInfoResp{
		Info: res,
	}, nil
}
