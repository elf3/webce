package request

type ReqLogin struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}
