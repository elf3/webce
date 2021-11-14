package permissions

// Permissions 菜单权限
type Permissions struct {
	ID          uint64 `gorm:"primary_key"  structs:"id" form:"id" json:"id"`                               // 主键ID
	Title       string `gorm:"type:varchar(50);unique_index" form:"title" json:"title" validate:"required"` // 权限标题
	Description string `gorm:"type:char(64);" form:"description" json:"description"`                        // 注解
	Slug        string `gorm:"type:varchar(50);" form:"slug" json:"slug"`                                   // 权限名称
	HttpPath    string `gorm:"type:text" form:"http_path" json:"http_path" validate:"required"`             // URI路径
	Method      string `gorm:"type:char(10);" form:"method" json:"method" validate:"required"`              // 请求方法
}
