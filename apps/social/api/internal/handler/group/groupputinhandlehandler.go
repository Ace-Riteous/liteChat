package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"liteChat/apps/social/api/internal/logic/group"
	"liteChat/apps/social/api/internal/svc"
	"liteChat/apps/social/api/internal/types"
)

func GroupPutInHandleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPutInHandleRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := group.NewGroupPutInHandleLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutInHandle(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
