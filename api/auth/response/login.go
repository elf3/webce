package response

import "webce/internal/repositories/models/admins"

type LoginResponse struct {
	*admins.Admin
	Token string `json:"token"`
}
