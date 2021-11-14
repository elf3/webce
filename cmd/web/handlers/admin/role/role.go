package role

import (
	"github.com/kataras/iris/v12"
	"webce/api/role"
	base "webce/cmd/web/handlers/admin"
	"webce/cmd/web/handlers/admin/role/request"
	"webce/pkg/library/log"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerRole struct {
	base.BaseHandler
	api role.ApiRole
}

func NewRoleHandler() *HandlerRole {
	return &HandlerRole{}
}

// GetList @Tags 角色管理
// @Router /admin/role/list [get]
func (p *HandlerRole) GetList() {
	where := iris.Map{}
	build, args, err := sql.WhereBuild(where)
	if err != nil {
		p.Error(303, "无法获取正确的参数")
		return
	}
	count, err := p.api.GetByCount(build, args)
	if err != nil {
		p.Error(304, "分页参数错误")
		return
	}
	pages := page.NewPagination(p.Ctx.Request(), count)
	lists, err := p.api.GetRolesPage(build, args, pages.GetPage(), pages.Perineum)
	if err != nil {
		p.Error(303, err.Error())
		return
	}
	p.Page(lists, pages.GetPageResp())
}

// GetDetail @Tags 角色管理
// @Router /admin/role/detail [get]
func (p *HandlerRole) GetDetail() {
	id := p.Ctx.URLParamUint64("id")

	if id <= 0 {
		p.Error(300, "please check id ")
		return
	}
	build, args, err := sql.WhereBuild(iris.Map{
		"id": id,
	})
	if err != nil {
		p.Error(300, "please check search condition ")
		return
	}
	data, err := p.api.Get(build, args)
	if err != nil {
		p.Error(400, "get detail error")
		return
	}
	p.Success(data)
}

// PostAdd @Tags 角色管理
// @Router /admin/role/add [post]
func (p *HandlerRole) PostAdd() {
	req := request.ReqAddRole{}
	err := p.Ctx.ReadForm(&req)
	if err != nil {
		log.Log.Error("read form request:", err)
		p.Error(500, "invalid request")
		return
	}
	err = p.Validate(req)
	if err != nil {
		log.Log.Error("invalid request:", err)
		p.Error(500, "invalid request")
		return
	}
	create, err := p.api.AddRole(req.PermissionIds, req.Roles)
	if err != nil {
		p.Error(300, err.Error())
		return
	}

	p.Success(create)
}

// PostEdit @Tags 角色管理
// @Router /admin/role/edit [post]
func (p *HandlerRole) PostEdit() {
	req := request.ReqEditRole{}
	err := p.Ctx.ReadForm(&req)
	if err != nil {
		log.Log.Error(" read form request ", err)
		p.Error(500, "invalid request")
		return
	}
	err = p.Validate(req)
	if err != nil {
		log.Log.Error(" invalid request ", err)
		p.Error(500, "invalid request")
		return
	}
	update, err := p.api.EditRole(req.Id, req.PermissionIds, req.Roles)
	if err != nil {
		p.Error(300, err.Error())
		return
	}
	p.Success(update)
}

// PostDelete @Tags 角色管理
// @Router /admin/role/delete [post]
func (p *HandlerRole) PostDelete() {
	id, err := p.Ctx.PostValueInt("id")

	if id <= 0 || err != nil {
		p.Error(300, "please check id ")
		return
	}

	err = p.api.DeleteRole(id)
	if err != nil {
		log.Log.Error("delete error : ", err)
		p.Error(400, "delete error")
		return
	}

	p.Msg("success")
}
