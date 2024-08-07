package group

import (
	"context"
	"liteChat/apps/social/rpc/social"
	"liteChat/pkg/constants"

	"liteChat/apps/social/api/internal/svc"
	"liteChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// todo: add your logic here and delete this line
	groupUser, err := l.svcCtx.Social.GroupUsers(l.ctx, &social.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, err
	}
	uids := make([]string, 0, len(groupUser.List))
	for _, user := range groupUser.List {
		uids = append(uids, user.UserId)
	}

	online, err := l.svcCtx.HgetallCtx(l.ctx, constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, err
	}

	resOnlineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		if _, ok := online[s]; ok {
			resOnlineList[s] = true
		} else {
			resOnlineList[s] = false
		}
	}
	return &types.GroupUserOnlineResp{
		OnlineList: resOnlineList,
	}, nil

}
