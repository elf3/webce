package permission

import (
	"github.com/kataras/iris/v12"
	admin "webce/cmd/web/handlers/admin/controller"
	"webce/cmd/web/handlers/admin/controller/permission/request"
	"webce/internal/repositories/models/admins/permissions"
	"webce/pkg/library/log"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerPermission struct {
	admin.BaseHandler
}

func NewPermissionHandler() *HandlerPermission {
	return &HandlerPermission{}
}

// GetList @Tags 权限管理
// @Router /admin/permission/list [get]
func (p *HandlerPermission) GetList() {
	model := permissions.Permissions{}
	where := iris.Map{}
	build, args, err := sql.WhereBuild(where)
	if err != nil {
		p.Error(303, "无法获取正确的参数")
		return
	}
	count := model.GetByCount(build, args)
	pages := page.NewPagination(p.Ctx.Request(), count)
	lists, err := model.Lists(build, args, pages.GetPage(), pages.Perineum)
	if err != nil {
		p.Error(303, err.Error())
		return
	}
	p.Page(lists, pages.GetPageResp())
}

// GetDetail @Tags 权限管理
// @Router /admin/permission/detail [get]
func (p *HandlerPermission) GetDetail() {
	id := p.Ctx.URLParamUint64("id")

	if id <= 0 {
		p.Error(300, "please check id ")
		return
	}
	p2 := permissions.Permissions{}
	build, args, err := sql.WhereBuild(iris.Map{
		"id": id,
	})
	if err != nil {
		p.Error(300, "please check search condition ")
		return
	}
	data, err := p2.Get(build, args)
	if err != nil {
		p.Error(400, "get detail error")
		return
	}
	p.Success(data)
}

// PostAdd @Tags 权限管理
// @Router /admin/permission/add [post]
func (p *HandlerPermission) PostAdd() {
	req := request.ReqAddPermission{}
	err := p.Ctx.ReadForm(&req)
	if err != nil {
		p.Error(500, "invalid request")
		return
	}
	err = req.Permissions.Create()
	if err != nil {
		p.Error(300, err.Error())
		return
	}

	p.Success(req)
}

// PostEdit @Tags 权限管理
// @Router /admin/permission/edit [post]
func (p *HandlerPermission) PostEdit() {
	req := request.ReqEditPermission{}
	err := p.Ctx.ReadForm(&req)
	if err != nil {
		log.Log.Error(" invalid request ", err)
		p.Error(500, "invalid request")
		return
	}

	update, err := req.Permissions.Update(req.Id)
	if err != nil {
		p.Error(300, err.Error())
		return
	}
	p.Success(update)
}

// PostDelete @Tags 权限管理
// @Router /admin/permission/delete [post]
func (p *HandlerPermission) PostDelete() {
	id, err := p.Ctx.PostValueInt64("id")

	if id <= 0 || err != nil {
		p.Error(300, "please check id ")
		return
	}
	p2 := permissions.Permissions{}
	err = p2.Delete(uint64(id))
	if err != nil {
		p.Error(400, "delete error")
		return
	}
	p.Success(p2)
}
