package business

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/internal/repositories/models"
	"webce/pkg/library/databases"
)

// Node 检测主机节点
type Node struct {
	models.BaseModel
	NodeName string `gorm:"size:255; unique_index;not null;" validate:"min=3,max=32" json:"node_name"` // 节点名称，唯一索引
	HostAddr string `gorm:"size:255;not null;" form:"host" json:"host" `                               // 节点地址
	Status   int    `gorm:"type:tinyint(1);not null" json:"status"`                                    // 节点状态
}

// CreateNode 创建节点
func (n *Node) CreateNode() (*Node, error) {
	create := databases.DB.Create(n)
	if create.Error != nil {
		return nil, create.Error
	}
	return n, nil
}

// SyncNode 同步节点
func (n *Node) SyncNode() (*Node, error) {
	first := databases.DB.Model(n).First(n)
	if first.Error != nil {
		return nil, first.Error
	}
	return n, nil
}

// Delete 删除节点
func (n *Node) Delete(id int64) error {
	first := databases.DB.Model(n).First(n)
	if first.Error != nil {
		return errors.New("not found node")
	}
	db := databases.DB.Model(&n).Where("id = ?", id).Delete(&n)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

// GetNodePage 获取节点列表
func (n *Node) GetNodePage(where map[string]interface{}, offset, limit int) ([]*Node, error) {
	list := make([]*Node, 0)
	find := databases.DB.Model(n).Where(where).Offset(offset).Limit(limit).Find(list)
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, find.Error
	}
	return list, nil
}
