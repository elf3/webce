package business

import (
	"webce/internal/repositories/models"
)

// Node 检测主机节点
type Node struct {
	models.BaseModel
	NodeName string `gorm:"size:255; unique_index;not null;" validate:"min=3,max=32" json:"node_name" form:"node_name"` // 节点名称，唯一索引
	HostAddr string `gorm:"size:255;not null;" json:"host" form:"host"`                                                 // 节点地址
	AreaId   int64  `gorm:"bigint(20); not null" json:"area_id"`                                                        // 区域ID
	NodeType int    `gorm:"int(11); not null" json:"node_type"`                                                         // 节点类型
	Status   int    `gorm:"type:tinyint(1);not null" json:"status" form:"status"`                                       // 节点状态
}
