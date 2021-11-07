package node

import (
	admin "webce/cmd/web/handlers/admin/controller"
	"webce/pkg/library/resp"
)

type HandlerNode struct {
	admin.BaseHandler
}

func NewNode() *HandlerNode {
	return &HandlerNode{}
}
func (h *HandlerNode) GetNodeLists() {
	//node := business.Node{}
	//node.CreateNode()
	resp.ApiJson(h.Ctx, 200, "", nil)
}
