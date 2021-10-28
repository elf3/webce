package admins

import "webce/apis/repositories/models"

type DetectType int64

const (
	_             DetectType = iota
	TypeHTTP                       // HTTP访问
	TypeGfw                        // GFW屏蔽
	TypeDnsHijack                  // 污染检测
	TypeHijack                     // 劫持检测
	TypePing                       // Ping检测
	TypeRoute                      // 路由跟踪
	Type301       DetectType = 301 // 301砌墙

)

type Domains struct {
	models.BaseModel
	DomainName string `gorm:"size:255; unique_index;not null;" validate:"min=3,max=32" json:"domain_name"` // 节点名称，唯一索引
	DetectType int    `gorm:"size:1;not null;" form:"host" json:"host" `                                   // 检测类型
	Status     int    `gorm:"type:tinyint(1);not null" json:"status"`                                      // 域名状态
}
