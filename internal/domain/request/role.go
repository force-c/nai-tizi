package request

import "github.com/force-c/nai-tizi/internal/utils/pagination"

// PageRoleRequest 查询角色列表请求
type PageRoleRequest struct {
	pagination.PageQuery        // 嵌入分页参数
	RoleName             string `json:"roleName" example:"管理员"` // 角色名称（模糊查询）
	Status               int32  `json:"status" example:"0"`     // 状态：0正常 1停用，-1表示不过滤
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleKey   string `json:"roleKey" binding:"required" example:"user_manager"` // 角色标识（唯一）
	RoleName  string `json:"roleName" binding:"required" example:"用户管理员"`       // 角色名称
	Sort      int64  `json:"sort" example:"1"`                                  // 显示顺序
	Status    int32  `json:"status" example:"0"`                                // 状态：0正常 1停用
	DataScope int32  `json:"dataScope" example:"2"`                             // 数据范围：1全部 2自定义 3本组织 4本组织及以下 5仅本人
	Remark    string `json:"remark" example:"负责用户管理"`                           // 备注
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleId    int64  `json:"roleId" binding:"required" example:"1"`       // 角色ID
	RoleName  string `json:"roleName" binding:"required" example:"用户管理员"` // 角色名称
	Sort      int64  `json:"sort" example:"1"`                            // 显示顺序
	Status    int32  `json:"status" example:"0"`                          // 状态：0正常 1停用
	DataScope int32  `json:"dataScope" example:"2"`                       // 数据范围
	Remark    string `json:"remark" example:"负责用户管理"`                     // 备注
}

// AssignRoleToUserRequest 为用户分配角色请求
type AssignRoleToUserRequest struct {
	UserId int64 `json:"userId" binding:"required" example:"1001"` // 用户ID
	RoleId int64 `json:"roleId" binding:"required" example:"1"`    // 角色ID
	OrgId  int64 `json:"orgId" binding:"required" example:"1"`     // 组织ID
}

// AddRolePermissionRequest 为角色添加权限请求
type AddRolePermissionRequest struct {
	RoleKey  string `json:"roleKey" binding:"required" example:"user_manager"` // 角色标识
	OrgId    int64  `json:"orgId" binding:"required" example:"1"`              // 组织ID
	Resource string `json:"resource" binding:"required" example:"user.*"`      // 资源路径（支持通配符）
	Action   string `json:"action" binding:"required" example:"write"`         // 操作类型（支持通配符）
}

// SetRoleInheritRequest 设置角色继承请求
type SetRoleInheritRequest struct {
	ChildRoleKey  string `json:"childRoleKey" binding:"required" example:"manager"` // 子角色标识
	ParentRoleKey string `json:"parentRoleKey" binding:"required" example:"viewer"` // 父角色标识
	OrgId         int64  `json:"orgId" binding:"required" example:"1"`              // 组织ID
}
