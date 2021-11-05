package permission

import (
	"github.com/kataras/iris/v12"
	admin "webce/cmd/web/handlers/admin/controller"
	"webce/cmd/web/handlers/admin/controller/permission/request"
	"webce/internal/repositories/models/admins/permissions"
	"webce/pkg/lib"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerPermission struct {
	admin.BaseHandler
}

func NewPermissionHandler() *HandlerPermission {
	return &HandlerPermission{}
}

func (p *HandlerPermission) GetList() {
	model := permissions.Permissions{}
	where := iris.Map{}
	build, args, err := sql.WhereBuild(where)
	if err != nil {
		lib.ErrJson(p.Ctx, 303, "无法获取正确的参数")
		return
	}
	count := model.GetByCount(build, args)
	pages := page.NewPagination(p.Ctx.Request(), count)
	lists, err := model.Lists(build, args, pages.GetPage(), pages.Perineum)
	if err != nil {
		lib.ErrJson(p.Ctx, 303, err.Error())
		return
	}
	lib.MJson(p.Ctx, 200, "", iris.Map{
		"lists": lists,
		"page":  pages.GetPageResp(),
	})
}

func (p *HandlerPermission) PostAdd() {
	req := request.ReqPermission{}
	err := p.Ctx.ReadForm(&req)
	if err != nil {
		lib.ErrJson(p.Ctx, 500, "invalid request")
		return
	}
	lib.Json(p.Ctx, req)
}
