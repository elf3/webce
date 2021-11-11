package admin

import (
	"github.com/kataras/iris/v12"
	base "webce/cmd/web/handlers/admin"
	"webce/cmd/web/handlers/admin/admin/request"
	admin2 "webce/internal/repositories/models/admins/admin"
	"webce/pkg/library/log"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerAdmin struct {
	base.BaseHandler
}

func NewAdminHandler() *HandlerAdmin {
	return &HandlerAdmin{}
}

// GetList @Tags 管理员
// @Router /admin/admin/list [get]
func (p *HandlerAdmin) GetList() {
	model := admin2.Admin{}
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

// GetDetail @Tags 管理员
// @Router /admin/admin/detail [get]
func (p *HandlerAdmin) GetDetail() {
	id := p.Ctx.URLParamUint64("id")

	if id <= 0 {
		p.Error(300, "please check id ")
		return
	}
	p2 := admin2.Admin{}
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

// PostAdd @Tags 管理员
// @Router /admin/admin/add [post]
func (p *HandlerAdmin) PostAdd() {
	req := request.ReqAddAdmin{}
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
	create, err := req.Admin.Create(req.RoleIds)
	if err != nil {
		p.Error(300, err.Error())
		return
	}

	p.Success(create)
}

// PostEdit @Tags 管理员
// @Router /admin/admin/edit [post]
func (p *HandlerAdmin) PostEdit() {
	req := request.ReqEditAdmin{}
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
	update, err := req.Admin.UpdateAdmin(req.ID, req.RoleIds)
	if err != nil {
		p.Error(300, err.Error())
		return
	}
	p.Success(update)
}

// PostDelete @Tags 管理员
// @Router /admin/admin/delete [post]
func (p *HandlerAdmin) PostDelete() {
	id, err := p.Ctx.PostValueInt("id")

	if id <= 0 || err != nil {
		p.Error(300, "please check id ")
		return
	}

	err = (admin2.Admin{}).Delete(id)
	if err != nil {
		log.Log.Error("delete error : ", err)
		p.Error(400, "delete error")
		return
	}

	p.Msg("success")
}
