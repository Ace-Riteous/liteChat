package logic

import (
	"context"
	"liteChat/apps/im/api/internal/svc"
	"liteChat/apps/im/api/internal/types"
	"liteChat/apps/im/rpc/im"
	"liteChat/apps/social/rpc/socialclient"
	"liteChat/apps/user/rpc/user"
	"liteChat/pkg/bitmap"
	"liteChat/pkg/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// todo: add your logic here and delete this line
	chatLogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		return nil, err
	}
	var (
		chatLog = chatLogs.List[0]
		reads   = []string{chatLog.SendId}
		unreads []string
		ids     []string
	)
	switch constants.ChatType(chatLog.ChatType) {
	case constants.SingleChatType:
		if len(chatLog.ReadRecords) == 0 || chatLog.ReadRecords[0] == 0 {
			unreads = []string{chatLog.RecvId}
		} else {
			reads = append(reads, chatLog.RecvId)
		}
		ids = []string{chatLog.SendId, chatLog.RecvId}
	case constants.GroupChatType:
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatLog.RecvId,
		})
		if err != nil {
			return nil, err
		}
		bitmaps := bitmap.Load(chatLog.ReadRecords)
		for _, member := range groupUsers.List {
			ids = append(ids, member.UserId)
			if member.UserId == chatLog.SendId {
				continue
			}
			if bitmaps.IsSet(member.UserId) {
				reads = append(reads, member.UserId)
			} else {
				unreads = append(unreads, member.UserId)
			}
		}
		userEntitys, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
			Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		userEntitysSet := make(map[string]*user.UserEntity, len(userEntitys.User))
		for i, entity := range userEntitys.User {
			userEntitysSet[entity.Id] = userEntitys.User[i]
		}
		for i, read := range reads {
			if u := userEntitysSet[read]; u != nil {
				reads[i] = u.Phone
			}
		}
		for i, unread := range unreads {
			if u := userEntitysSet[unread]; u != nil {
				unreads[i] = u.Phone
			}
		}
	}
	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
