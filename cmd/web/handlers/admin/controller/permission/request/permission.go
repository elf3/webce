package request

import (
	"webce/internal/repositories/models/admins/permissions"
)

type ReqAddPermission struct {
	permissions.Permissions
}
type ReqEditPermission struct {
	Id uint64 `form:"id" validate:"required" `
	permissions.Permissions
}
