package business

import (
	"webce/internal/repositories/models"
)

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

// Domains 域名表
type Domains struct {
	models.BaseModel
	DetectType   int    `gorm:"size:1;not null;" form:"host" json:"host" `                   // 检测类型
	Protocol     int    `gorm:"size:1; not null" form:"protocol" json:"protocol"`            // 协议
	Method       int    `gorm:"size:1; not null" form:"method" json:"method"`                // 请求方法
	RequestBody  string `gorm:"type:text; not null" form:"request_body" json:"request_body"` // 请求主体
	Url          string `gorm:"size:255; not null" form:"url" json:"url"`                    // 检测的连接
	ResponseTime int    `gorm:"size:1; not null" form:"response_time" json:"response_time"`  // 响应时间
	Status       int    `gorm:"type:tinyint(1);not null" json:"status"`                      // 域名状态
	DetectName   string `gorm:"size:255; not null" form:"detect_name" json:"detect_name"`    // 检测的连接
	Nodes        []Node `gorm:"foreignKey:id;references:Nodes"`
	//DomainName string `gorm:"size:255; unique_index;not null;" validate:"min=3,max=32" json:"domain_name"` // 域名 唯一索引

}
