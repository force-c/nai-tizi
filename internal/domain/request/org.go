package request

import "github.com/force-c/nai-tizi/internal/utils/pagination"

// CreateOrgRequest 创建组织请求
type CreateOrgRequest struct {
	ParentId int64  `json:"parentId"`                                                   // 父组织ID，0表示根组织
	OrgName  string `json:"orgName" binding:"required,min=2,max=50"`                    // 组织名称
	OrgCode  string `json:"orgCode" binding:"required,min=2,max=30"`                    // 组织编码
	OrgType  string `json:"orgType" binding:"omitempty,oneof=company department group"` // 组织类型
	Leader   string `json:"leader" binding:"omitempty,max=50"`                          // 负责人
	Phone    string `json:"phone" binding:"omitempty,min=11,max=11"`                    // 联系电话
	Email    string `json:"email" binding:"omitempty,email"`                            // 邮箱
	Status   int32  `json:"status" binding:"omitempty,oneof=0 1"`                       // 状态：0正常 1停用
	Sort     int64  `json:"sort"`                                                       // 显示顺序
	Remark   string `json:"remark" binding:"omitempty,max=500"`                         // 备注
	CreateBy int64  `json:"-"`                                                          // 从上下文获取
	UpdateBy int64  `json:"-"`                                                          // 从上下文获取
}

// UpdateOrgRequest 更新组织请求
type UpdateOrgRequest struct {
	OrgId    int64  `json:"-"`                                                          // 从路径参数获取
	ParentId int64  `json:"parentId"`                                                   // 父组织ID
	OrgName  string `json:"orgName" binding:"omitempty,min=2,max=50"`                   // 组织名称
	OrgCode  string `json:"orgCode" binding:"omitempty,min=2,max=30"`                   // 组织编码
	OrgType  string `json:"orgType" binding:"omitempty,oneof=company department group"` // 组织类型
	Leader   string `json:"leader" binding:"omitempty,max=50"`                          // 负责人
	Phone    string `json:"phone" binding:"omitempty,min=11,max=11"`                    // 联系电话
	Email    string `json:"email" binding:"omitempty,email"`                            // 邮箱
	Status   int32  `json:"status" binding:"omitempty,oneof=0 1"`                       // 状态：0正常 1停用
	Sort     int64  `json:"sort"`                                                       // 显示顺序
	Remark   string `json:"remark" binding:"omitempty,max=500"`                         // 备注
	UpdateBy int64  `json:"-"`                                                          // 从上下文获取
}

// BatchDeleteOrgsRequest 批量删除组织请求
type BatchDeleteOrgsRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1"` // 组织ID列表
}

// PageOrgsRequest 查询组织列表请求
type PageOrgsRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	OrgName              string `json:"orgName"`                              // 组织名称（模糊查询）
	OrgCode              string `json:"orgCode"`                              // 组织编码
	Status               int32  `json:"status" binding:"omitempty,oneof=0 1"` // 状态：0正常 1停用
	ParentId             *int64 `json:"parentId"`                             // 父组织ID（可选）
}
