package response

import (
	"webce/internal/repositories/models/admins/admin"
)

type LoginResponse struct {
	*admin.Admin
	Token string `json:"token"`
}
