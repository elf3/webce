package node

import (
	"github.com/kataras/iris/v12"
	"webce/api/node"
	"webce/cmd/web/handlers/admin"
	"webce/cmd/web/handlers/admin/node/request"
	"webce/pkg/library/log"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerNode struct {
	admin.BaseHandler
	api node.ApiNode
}

func NewNode() *HandlerNode {
	return &HandlerNode{}
}

// GetLists 获取节点列表
// @Router /admin/node/lists [get]
func (w HandlerNode) GetLists() {
	where := iris.Map{}
	build, args, err := sql.WhereBuild(where)
	if err != nil {
		log.Log.Error("get by search err", err)
		w.Error(400, "get by search err")
		return
	}
	count := w.api.GetByCount(build, args)
	pages := page.NewPagination(w.Ctx.Request(), count)
	getPage, err := w.api.GetPage(build, args, pages.GetPage(), pages.Perineum)
	if err != nil {
		log.Log.Error("get node list err", err)
		w.Error(403, "get node list err")
	}
	w.Page(getPage, pages.GetPageResp())
}

// PostAdd 添加节点
// @Router /admin/node/add [post]
func (w HandlerNode) PostAdd() {
	req := request.ReqAddNode{}
	err := w.Ctx.ReadForm(&req)
	if err != nil {
		w.Error(400, "Invalid request")
		return
	}
	err = w.Validate(req)
	if err != nil {
		w.Error(402, "Invalid request")
		return
	}
	createNode, err := w.api.CreateNode(req.Node)
	if err != nil {
		w.Error(403, "err to add node")
		return
	}
	w.Success(createNode)
}

// PostEdit 添加节点
// @Router /admin/node/edit [post]
func (w HandlerNode) PostEdit() {
	req := request.ReqEditNode{}
	err := w.Ctx.ReadForm(&req)
	if err != nil {
		w.Error(400, "Invalid request")
		return
	}
	err = w.Validate(req)
	if err != nil {
		w.Error(402, "Invalid request")
		return
	}
	createNode, err := w.api.UpdateNode(req.Id, req.Node)
	if err != nil {
		w.Error(403, "err to add node")
		return
	}
	w.Success(createNode)
}

// PostDelete 添加节点
// @Router /admin/node/delete [post]
func (w HandlerNode) PostDelete() {
	id, err := w.Ctx.PostValueInt64("id")

	if id <= 0 || err != nil {
		w.Error(300, "please check id ")
		return
	}

	err = w.api.Delete(id)
	if err != nil {
		log.Log.Error("delete error : ", err)
		w.Error(400, "delete error")
		return
	}

	w.Msg("success")
}
