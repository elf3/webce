package node

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"webce/internal/repositories/models/business"
	"webce/pkg/library/databases"
)

type ApiNode struct {
}

// CreateNode 创建节点
func (a ApiNode) CreateNode(n business.Node) (*business.Node, error) {

	err := databases.DB.Model(n).Create(n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// UpdateNode 创建节点
func (a ApiNode) UpdateNode(id int64, n business.Node) (*business.Node, error) {
	find := business.Node{}
	err := databases.DB.Model(&find).Where("id = ?", id).First(&find).Error
	if err != nil {
		return nil, err
	}
	err = databases.DB.Model(n).Create(n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// SyncApiNode 同步节点
func (a ApiNode) SyncApiNode() (*business.Node, error) {
	n := business.Node{}
	err := databases.DB.Model(n).First(n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// Delete 删除节点
func (a ApiNode) Delete(id int64) error {
	n := business.Node{}
	first := databases.DB.Model(n).First(n)
	if first.Error != nil {
		return errors.New("not found node")
	}
	err := databases.DB.Model(&n).Where("id = ?", id).Delete(&n).Error
	if err != nil {
		return err
	}

	return nil
}

// GetPage  获取节点列表
func (a ApiNode) GetPage(where string, vals []interface{}, offset, limit int) ([]business.Node, error) {
	list := make([]business.Node, 0)
	n := business.Node{}
	find := databases.DB.Model(&n).Where(where, vals...).Offset(offset).Limit(limit).Find(&list)
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, find.Error
	}
	return list, nil
}

// GetByCount 获取有多少条记录
func (a ApiNode) GetByCount(whereSql string, vals []interface{}) (count int64) {
	n := business.Node{}
	databases.DB.Model(&n).Where(whereSql, vals...).Count(&count)
	return
}
