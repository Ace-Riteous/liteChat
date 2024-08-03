package logic

import (
	"context"

	"liteChat/apps/user/rpc/internal/svc"
	"liteChat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *user.Request) (*user.Respond, error) {
	// todo: add your types here and delete this line

	return &user.Respond{
		Pong: "riteous",
	}, nil
}
