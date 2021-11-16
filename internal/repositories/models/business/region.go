package business

// Region 省市区
type Region struct {
	RegionId        int64  `json:"region_id" form:"region_id"`                 // 地区ID
	RegionName      string `json:"region_name" form:"region_name"`             // 地区名称
	RegionShortName string `json:"region_short_name" form:"region_short_name"` // 地区缩写
	RegionCode      int    `json:"region_code" form:"region_code"`             // 邮政代码
	RegionParentId  int64  `json:"region_parent_id" form:"region_parent_id"`   // 父级ID
	RegionLevel     int    `json:"region_level" form:"region_level"`           // 注册等级
}
