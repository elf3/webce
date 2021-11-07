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
