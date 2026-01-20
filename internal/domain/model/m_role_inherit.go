package model

import (
	"github.com/force-c/nai-tizi/internal/utils"
	"gorm.io/gorm"
)

// MRoleInherit 角色继承关系表（映射表，支持多继承）
type MRoleInherit struct {
	Id          int64           `gorm:"column:id;primaryKey" autogen:"int64" json:"id"`                  // 使用分布式ID
	RoleId      int64           `gorm:"column:role_id;not null;index:idx_role_parent" json:"roleId"`     // 子角色ID
	ParentId    int64           `gorm:"column:parent_id;not null;index:idx_role_parent" json:"parentId"` // 父角色ID
	TenantId    int64           `gorm:"column:tenant_id;not null;default:1" json:"tenantId"`             // 租户ID（预留多租户）
	CreateBy    int64           `gorm:"column:create_by" json:"createBy"`                                // 创建人
	UpdateBy    int64           `gorm:"column:update_by" json:"updateBy"`                                // 更新人
	CreatedTime utils.LocalTime `gorm:"column:created_time;autoCreateTime" json:"createdTime"`
	UpdatedTime utils.LocalTime `gorm:"column:updated_time;autoUpdateTime" json:"updatedTime"`
	DeletedAt   gorm.DeletedAt  `gorm:"column:deleted_at;index" json:"-"`
}

func (*MRoleInherit) TableName() string { return "m_role_inherit" }

// FindByRoleId 根据子角色ID查询父角色关联
func (m *MRoleInherit) FindByRoleId(db *gorm.DB, roleId int64) ([]MRoleInherit, error) {
	var inherits []MRoleInherit
	err := db.Where("role_id = ?", roleId).Find(&inherits).Error
	return inherits, err
}

// FindByParentId 根据父角色ID查询子角色关联
func (m *MRoleInherit) FindByParentId(db *gorm.DB, parentId int64) ([]MRoleInherit, error) {
	var inherits []MRoleInherit
	err := db.Where("parent_id = ?", parentId).Find(&inherits).Error
	return inherits, err
}

// Exists 检查角色继承关系是否存在
func (m *MRoleInherit) Exists(db *gorm.DB, roleId, parentId int64) (bool, error) {
	var count int64
	err := db.Model(&MRoleInherit{}).
		Where("role_id = ? AND parent_id = ?", roleId, parentId).
		Count(&count).Error
	return count > 0, err
}

// Create 创建角色继承关系
func (m *MRoleInherit) Create(db *gorm.DB, inherit *MRoleInherit) error {
	return db.Create(inherit).Error
}

// Delete 删除角色继承关系
func (m *MRoleInherit) Delete(db *gorm.DB, roleId, parentId int64) error {
	return db.Where("role_id = ? AND parent_id = ?", roleId, parentId).Delete(&MRoleInherit{}).Error
}

// DeleteByRoleId 删除子角色的所有继承关系
func (m *MRoleInherit) DeleteByRoleId(db *gorm.DB, roleId int64) error {
	return db.Where("role_id = ?", roleId).Delete(&MRoleInherit{}).Error
}

// DeleteByParentId 删除父角色的所有继承关系
func (m *MRoleInherit) DeleteByParentId(db *gorm.DB, parentId int64) error {
	return db.Where("parent_id = ?", parentId).Delete(&MRoleInherit{}).Error
}

// DeleteByRoleIdOrParentId 删除涉及指定角色的所有继承关系（作为子角色或父角色）
func (m *MRoleInherit) DeleteByRoleIdOrParentId(db *gorm.DB, roleId int64) error {
	return db.Where("role_id = ? OR parent_id = ?", roleId, roleId).Delete(&MRoleInherit{}).Error
}

// CountByParentId 统计继承该父角色的子角色数量
func (m *MRoleInherit) CountByParentId(db *gorm.DB, parentId int64) (int64, error) {
	var count int64
	err := db.Model(&MRoleInherit{}).Where("parent_id = ?", parentId).Count(&count).Error
	return count, err
}

// GetParentRoleIds 获取角色的所有父角色ID列表
func (m *MRoleInherit) GetParentRoleIds(db *gorm.DB, roleId int64) ([]int64, error) {
	var inherits []MRoleInherit
	err := db.Where("role_id = ?", roleId).Find(&inherits).Error
	if err != nil {
		return nil, err
	}

	parentIds := make([]int64, 0, len(inherits))
	for _, inherit := range inherits {
		parentIds = append(parentIds, inherit.ParentId)
	}
	return parentIds, nil
}

// GetChildRoleIds 获取角色的所有子角色ID列表
func (m *MRoleInherit) GetChildRoleIds(db *gorm.DB, parentId int64) ([]int64, error) {
	var inherits []MRoleInherit
	err := db.Where("parent_id = ?", parentId).Find(&inherits).Error
	if err != nil {
		return nil, err
	}

	childIds := make([]int64, 0, len(inherits))
	for _, inherit := range inherits {
		childIds = append(childIds, inherit.RoleId)
	}
	return childIds, nil
}
