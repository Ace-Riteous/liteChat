package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"liteChat/apps/im/immodels"
	"liteChat/apps/im/rpc/im"
	"liteChat/apps/im/rpc/internal/svc"
	"liteChat/pkg/constants"
	"liteChat/pkg/wuid"
	"liteChat/pkg/xerr"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		conversationRes, err := l.svcCtx.ConversationsModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if err == immodels.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})
				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "converastion insert err %v", err)
				}
			}
		} else if conversationRes != nil {
			return nil, nil
		} else {
			return nil, errors.Wrapf(xerr.NewDBErr(), "conversation findone err %v, req %v", err, conversationId)
		}

		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.Wrapf(xerr.NewInternalErr(), "no such chat_type, rea %v", in.ChatType)
	}

	return &im.SetUpUserConversationResp{}, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string, chatType constants.ChatType, isShow bool) error {
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.ObjectID{},
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}

		}
		return errors.Wrapf(xerr.NewDBErr(), "find by user_id err %v, req %v", err, userId)
	}

	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}
	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "update conversation err %v, req %v", err, conversations)
	}
	return nil
}
