package permission

import (
	"github.com/kataras/iris/v12"
	"webce/api/permission"
	"webce/cmd/web/handlers/admin"
	"webce/cmd/web/handlers/admin/permission/request"
	"webce/pkg/library/log"
	"webce/pkg/library/page"
	"webce/pkg/library/sql"
)

type HandlerPermission struct {
	admin.BaseHandler
	api permission.ApiPermission
}

func NewPermissionHandler() *HandlerPermission {
	return &HandlerPermission{}
}

// GetList @Tags 权限管理
// @Router /admin/permission/list [get]
func (w *HandlerPermission) GetList() {
	where := iris.Map{}
	build, args, err := sql.WhereBuild(where)
	if err != nil {
		w.Error(303, "无法获取正确的参数")
		return
	}
	count := w.api.GetByCount(build, args)
	pages := page.NewPagination(w.Ctx.Request(), count)
	lists, err := w.api.Lists(build, args, pages.GetPage(), pages.Perineum)
	if err != nil {
		w.Error(303, err.Error())
		return
	}
	w.Page(lists, pages.GetPageResp())
}

// GetDetail @Tags 权限管理
// @Router /admin/permission/detail [get]
func (w *HandlerPermission) GetDetail() {
	id := w.Ctx.URLParamUint64("id")

	if id <= 0 {
		w.Error(300, "please check id ")
		return
	}
	build, args, err := sql.WhereBuild(iris.Map{
		"id": id,
	})
	if err != nil {
		w.Error(300, "please check search condition ")
		return
	}
	data, err := w.api.Get(build, args)
	if err != nil {
		w.Error(400, "get detail error")
		return
	}
	w.Success(data)
}

// PostAdd @Tags 权限管理
// @Router /admin/permission/add [post]
func (w *HandlerPermission) PostAdd() {
	req := request.ReqAddPermission{}
	err := w.Ctx.ReadForm(&req)
	if err != nil {
		w.Error(500, "invalid request")
		return
	}
	err = w.api.Create(req.Permissions)
	if err != nil {
		w.Error(300, err.Error())
		return
	}

	w.Success(req)
}

// PostEdit @Tags 权限管理
// @Router /admin/permission/edit [post]
func (w *HandlerPermission) PostEdit() {
	req := request.ReqEditPermission{}
	err := w.Ctx.ReadForm(&req)
	if err != nil {
		log.Log.Error(" invalid request ", err)
		w.Error(500, "invalid request")
		return
	}

	update, err := w.api.Update(req.Id, req.Permissions)
	if err != nil {
		w.Error(300, err.Error())
		return
	}
	w.Success(update)
}

// PostDelete @Tags 权限管理
// @Router /admin/permission/delete [post]
func (w *HandlerPermission) PostDelete() {
	id, err := w.Ctx.PostValueInt64("id")

	if id <= 0 || err != nil {
		w.Error(300, "please check id ")
		return
	}
	err = w.api.Delete(uint64(id))
	if err != nil {
		w.Error(400, "delete error")
		return
	}
	w.Success("success")
}
