package requests

type Login struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required,gte=6,lte=30"  comment:"密码"`
}
