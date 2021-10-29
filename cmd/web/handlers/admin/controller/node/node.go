package node

import (
	admin "webce/cmd/web/handlers/admin/controller"
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
	h.ApiJson(200, "", nil)
}
