package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/force-c/nai-tizi/internal/domain/model"
	"github.com/force-c/nai-tizi/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CasbinService Casbin 权限管理服务接口
type CasbinService interface {
	// CheckPermission 检查用户权限
	// userId: 用户ID
	// orgId: 组织ID
	// resource: 资源路径（例如: "user.create", "device.read"）
	// action: 操作类型（例如: "read", "write", "delete"）
	CheckPermission(ctx context.Context, userId, orgId int64, resource, action string) (bool, error)

	// AddRoleForUser 为用户分配角色
	// userId: 用户ID
	// roleKey: 角色标识（例如: "admin", "user_manager"）
	// orgId: 组织ID
	AddRoleForUser(ctx context.Context, userId int64, roleKey string, orgId int64) error

	// DeleteRoleForUser 删除用户的角色
	DeleteRoleForUser(ctx context.Context, userId int64, roleKey string, orgId int64) error

	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(ctx context.Context, userId, orgId int64) ([]string, error)

	// GetUsersForRole 获取拥有指定角色的所有用户
	GetUsersForRole(ctx context.Context, roleKey string, orgId int64) ([]int64, error)

	// AddPermissionForRole 为角色添加权限
	// roleKey: 角色标识
	// orgId: 组织ID
	// resource: 资源路径（支持通配符，例如: "user.*", "*.read", "*"）
	// action: 操作类型（支持通配符，例如: "write", "*"）
	AddPermissionForRole(ctx context.Context, roleKey string, orgId int64, resource, action string) error

	// DeletePermissionForRole 删除角色的权限
	DeletePermissionForRole(ctx context.Context, roleKey string, orgId int64, resource, action string) error

	// GetPermissionsForRole 获取角色的所有权限
	GetPermissionsForRole(ctx context.Context, roleKey string, orgId int64) ([][]string, error)

	// SetRoleInherit 设置角色继承关系（子角色继承父角色的所有权限）
	// childRoleKey: 子角色标识
	// parentRoleKey: 父角色标识
	// orgId: 组织ID
	SetRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) error

	// DeleteRoleInherit 删除角色继承关系
	DeleteRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) error

	// GetRoleInherits 获取角色的所有父角色
	GetRoleInherits(ctx context.Context, roleKey string, orgId int64) ([]string, error)

	// CheckCircularInherit 检查是否存在循环继承
	CheckCircularInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) (bool, error)

	// GetInheritDepth 获取角色继承深度
	GetInheritDepth(ctx context.Context, roleKey string, orgId int64) (int, error)

	// ReloadPolicy 重新加载策略（从数据库）
	ReloadPolicy(ctx context.Context) error
}

type casbinService struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	logger   logger.Logger
}

// NewCasbinService 创建 Casbin 服务实例
func NewCasbinService(enforcer *casbin.Enforcer, db *gorm.DB, logger logger.Logger) CasbinService {
	return &casbinService{
		enforcer: enforcer,
		db:       db,
		logger:   logger,
	}
}

// CheckPermission 检查用户权限
// 示例: CheckPermission(ctx, 1001, 1, "user.create", "write")
// 含义: 检查用户1001在组织1中是否有权限对 user.create 资源执行 write 操作
func (s *casbinService) CheckPermission(ctx context.Context, userId, orgId int64, resource, action string) (bool, error) {
	// 构造 Casbin 请求参数
	// sub: "user::1001"
	// dom: "org::1"
	// obj: "user.create"
	// act: "write"
	sub := fmt.Sprintf("user::%d", userId)
	dom := fmt.Sprintf("org::%d", orgId)

	// 执行权限检查
	ok, err := s.enforcer.Enforce(sub, dom, resource, action)
	if err != nil {
		s.logger.Error("权限检查失败",
			zap.Int64("userId", userId),
			zap.Int64("orgId", orgId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err))
		return false, fmt.Errorf("权限检查失败: %w", err)
	}

	s.logger.Debug("权限检查",
		zap.Int64("userId", userId),
		zap.Int64("orgId", orgId),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.Bool("allowed", ok))

	return ok, nil
}

// AddRoleForUser 为用户分配角色
// 示例: AddRoleForUser(ctx, 1001, "admin", 1)
// 含义: 为用户1001在组织1中分配 admin 角色
func (s *casbinService) AddRoleForUser(ctx context.Context, userId int64, roleKey string, orgId int64) error {
	// 检查角色是否存在
	var role model.Role
	if err := s.db.Where("role_key = ?", roleKey).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在: %s", roleKey)
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 构造 Casbin 参数
	sub := fmt.Sprintf("user::%d", userId)
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 添加角色关系到 Casbin
	// g, user::1001, role::admin, org::1
	if _, err := s.enforcer.AddGroupingPolicy(sub, roleStr, dom); err != nil {
		s.logger.Error("添加用户角色失败",
			zap.Int64("userId", userId),
			zap.String("roleKey", roleKey),
			zap.Int64("orgId", orgId),
			zap.Error(err))
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	// 同步到数据库
	userRole := &model.MUserRole{
		UserId: userId,
		RoleId: role.ID,
	}
	if err := s.db.Create(userRole).Error; err != nil {
		// 如果数据库插入失败，回滚 Casbin 策略
		s.enforcer.RemoveGroupingPolicy(sub, roleStr, dom)
		return fmt.Errorf("保存用户角色关系失败: %w", err)
	}

	s.logger.Info("为用户分配角色成功",
		zap.Int64("userId", userId),
		zap.String("roleKey", roleKey),
		zap.Int64("orgId", orgId))

	return nil
}

// DeleteRoleForUser 删除用户的角色
func (s *casbinService) DeleteRoleForUser(ctx context.Context, userId int64, roleKey string, orgId int64) error {
	// 构造 Casbin 参数
	sub := fmt.Sprintf("user::%d", userId)
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 从 Casbin 删除角色关系
	if _, err := s.enforcer.RemoveGroupingPolicy(sub, roleStr, dom); err != nil {
		s.logger.Error("删除用户角色失败",
			zap.Int64("userId", userId),
			zap.String("roleKey", roleKey),
			zap.Int64("orgId", orgId),
			zap.Error(err))
		return fmt.Errorf("删除用户角色失败: %w", err)
	}

	// 从数据库删除
	var role model.Role
	if err := s.db.Where("role_key = ?", roleKey).First(&role).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询角色失败: %w", err)
		}
	} else {
		if err := s.db.Where("user_id = ? AND role_id = ? AND org_id = ?", userId, role.ID, orgId).
			Delete(&model.MUserRole{}).Error; err != nil {
			return fmt.Errorf("删除用户角色关系失败: %w", err)
		}
	}

	s.logger.Info("删除用户角色成功",
		zap.Int64("userId", userId),
		zap.String("roleKey", roleKey),
		zap.Int64("orgId", orgId))

	return nil
}

// GetRolesForUser 获取用户的所有角色
func (s *casbinService) GetRolesForUser(ctx context.Context, userId, orgId int64) ([]string, error) {
	sub := fmt.Sprintf("user::%d", userId)
	dom := fmt.Sprintf("org::%d", orgId)

	// 获取用户的所有角色（包括继承的角色）
	roles, err := s.enforcer.GetImplicitRolesForUser(sub, dom)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 去除 "role::" 前缀
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		if len(role) > 6 && role[:6] == "role::" {
			result = append(result, role[6:])
		}
	}

	return result, nil
}

// GetUsersForRole 获取拥有指定角色的所有用户
func (s *casbinService) GetUsersForRole(ctx context.Context, roleKey string, orgId int64) ([]int64, error) {
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 获取所有拥有该角色的用户
	users := s.enforcer.GetUsersForRoleInDomain(roleStr, dom)

	// 解析用户ID
	result := make([]int64, 0, len(users))
	for _, user := range users {
		var userId int64
		if _, err := fmt.Sscanf(user, "user::%d", &userId); err == nil {
			result = append(result, userId)
		}
	}

	return result, nil
}

// AddPermissionForRole 为角色添加权限
// 示例: AddPermissionForRole(ctx, "user_manager", 1, "user.*", "write")
// 含义: 为 user_manager 角色在组织1中添加对 user.* 资源的 write 权限
// 支持通配符: "user.*" 表示用户模块所有操作, "*.read" 表示所有模块的读操作, "*" 表示所有权限
func (s *casbinService) AddPermissionForRole(ctx context.Context, roleKey string, orgId int64, resource, action string) error {
	// 检查角色是否存在
	var role model.Role
	if err := s.db.Where("role_key = ?", roleKey).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在: %s", roleKey)
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 构造 Casbin 参数
	// p, role::user_manager, org::1, user.*, write
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 添加权限策略
	if _, err := s.enforcer.AddPolicy(roleStr, dom, resource, action); err != nil {
		s.logger.Error("添加角色权限失败",
			zap.String("roleKey", roleKey),
			zap.Int64("orgId", orgId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err))
		return fmt.Errorf("添加角色权限失败: %w", err)
	}

	s.logger.Info("为角色添加权限成功",
		zap.String("roleKey", roleKey),
		zap.Int64("orgId", orgId),
		zap.String("resource", resource),
		zap.String("action", action))

	return nil
}

// DeletePermissionForRole 删除角色的权限
func (s *casbinService) DeletePermissionForRole(ctx context.Context, roleKey string, orgId int64, resource, action string) error {
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 删除权限策略
	if _, err := s.enforcer.RemovePolicy(roleStr, dom, resource, action); err != nil {
		s.logger.Error("删除角色权限失败",
			zap.String("roleKey", roleKey),
			zap.Int64("orgId", orgId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err))
		return fmt.Errorf("删除角色权限失败: %w", err)
	}

	s.logger.Info("删除角色权限成功",
		zap.String("roleKey", roleKey),
		zap.Int64("orgId", orgId),
		zap.String("resource", resource),
		zap.String("action", action))

	return nil
}

// GetPermissionsForRole 获取角色的所有权限
func (s *casbinService) GetPermissionsForRole(ctx context.Context, roleKey string, orgId int64) ([][]string, error) {
	roleStr := fmt.Sprintf("role::%s", roleKey)

	// 获取角色的所有权限（包括继承的权限）
	permissions, err := s.enforcer.GetImplicitPermissionsForUser(roleStr, fmt.Sprintf("org::%d", orgId))
	if err != nil {
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}

	return permissions, nil
}

// SetRoleInherit 设置角色继承关系
// 示例: SetRoleInherit(ctx, "manager", "viewer", 1)
// 含义: manager 角色继承 viewer 角色的所有权限
// 注意: 会检查循环继承和继承深度限制（最多3层）
func (s *casbinService) SetRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) error {
	// 检查角色是否存在
	var childRole, parentRole model.Role
	if err := s.db.Where("role_key = ?", childRoleKey).First(&childRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("子角色不存在: %s", childRoleKey)
		}
		return fmt.Errorf("查询子角色失败: %w", err)
	}
	if err := s.db.Where("role_key = ?", parentRoleKey).First(&parentRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("父角色不存在: %s", parentRoleKey)
		}
		return fmt.Errorf("查询父角色失败: %w", err)
	}

	// 检查是否会造成循环继承
	hasCircular, err := s.CheckCircularInherit(ctx, childRoleKey, parentRoleKey, orgId)
	if err != nil {
		return fmt.Errorf("检查循环继承失败: %w", err)
	}
	if hasCircular {
		return fmt.Errorf("不能设置循环继承: %s -> %s", childRoleKey, parentRoleKey)
	}

	// 检查继承深度（最多3层）
	depth, err := s.GetInheritDepth(ctx, parentRoleKey, orgId)
	if err != nil {
		return fmt.Errorf("获取继承深度失败: %w", err)
	}
	if depth >= 3 {
		return fmt.Errorf("角色继承深度超过限制（最多3层）")
	}

	// 构造 Casbin 参数
	// g, role::manager, role::viewer, org::1
	childRoleStr := fmt.Sprintf("role::%s", childRoleKey)
	parentRoleStr := fmt.Sprintf("role::%s", parentRoleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 添加角色继承关系到 Casbin
	if _, err := s.enforcer.AddGroupingPolicy(childRoleStr, parentRoleStr, dom); err != nil {
		s.logger.Error("设置角色继承失败",
			zap.String("childRole", childRoleKey),
			zap.String("parentRole", parentRoleKey),
			zap.Int64("orgId", orgId),
			zap.Error(err))
		return fmt.Errorf("设置角色继承失败: %w", err)
	}

	// 同步到数据库
	roleInherit := &model.MRoleInherit{
		RoleId:   childRole.ID,
		ParentId: parentRole.ID,
		TenantId: 1, // 默认租户ID为1（单一企业模式）
	}
	if err := s.db.Create(roleInherit).Error; err != nil {
		// 如果数据库插入失败，回滚 Casbin 策略
		s.enforcer.RemoveGroupingPolicy(childRoleStr, parentRoleStr, dom)
		return fmt.Errorf("保存角色继承关系失败: %w", err)
	}

	s.logger.Info("设置角色继承成功",
		zap.String("childRole", childRoleKey),
		zap.String("parentRole", parentRoleKey),
		zap.Int64("orgId", orgId))

	return nil
}

// DeleteRoleInherit 删除角色继承关系
func (s *casbinService) DeleteRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) error {
	childRoleStr := fmt.Sprintf("role::%s", childRoleKey)
	parentRoleStr := fmt.Sprintf("role::%s", parentRoleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 从 Casbin 删除角色继承关系
	if _, err := s.enforcer.RemoveGroupingPolicy(childRoleStr, parentRoleStr, dom); err != nil {
		s.logger.Error("删除角色继承失败",
			zap.String("childRole", childRoleKey),
			zap.String("parentRole", parentRoleKey),
			zap.Int64("orgId", orgId),
			zap.Error(err))
		return fmt.Errorf("删除角色继承失败: %w", err)
	}

	// 从数据库删除
	var childRole, parentRole model.Role
	if err := s.db.Where("role_key = ?", childRoleKey).First(&childRole).Error; err == nil {
		if err := s.db.Where("role_key = ?", parentRoleKey).First(&parentRole).Error; err == nil {
			if err := s.db.Where("role_id = ? AND parent_id = ? AND org_id = ?",
				childRole.ID, parentRole.ID, orgId).
				Delete(&model.MRoleInherit{}).Error; err != nil {
				return fmt.Errorf("删除角色继承关系失败: %w", err)
			}
		}
	}

	s.logger.Info("删除角色继承成功",
		zap.String("childRole", childRoleKey),
		zap.String("parentRole", parentRoleKey),
		zap.Int64("orgId", orgId))

	return nil
}

// GetRoleInherits 获取角色的所有父角色
func (s *casbinService) GetRoleInherits(ctx context.Context, roleKey string, orgId int64) ([]string, error) {
	roleStr := fmt.Sprintf("role::%s", roleKey)
	dom := fmt.Sprintf("org::%d", orgId)

	// 获取角色的所有父角色
	parents, err := s.enforcer.GetImplicitRolesForUser(roleStr, dom)
	if err != nil {
		return nil, fmt.Errorf("获取角色继承失败: %w", err)
	}

	// 去除 "role::" 前缀，并排除自己
	result := make([]string, 0)
	for _, parent := range parents {
		if len(parent) > 6 && parent[:6] == "role::" {
			parentKey := parent[6:]
			if parentKey != roleKey {
				result = append(result, parentKey)
			}
		}
	}

	return result, nil
}

// CheckCircularInherit 检查是否存在循环继承
// 例如: A -> B -> C -> A 就是循环继承
func (s *casbinService) CheckCircularInherit(ctx context.Context, childRoleKey, parentRoleKey string, orgId int64) (bool, error) {
	// 如果父角色的祖先中包含子角色，则存在循环
	ancestors, err := s.GetRoleInherits(ctx, parentRoleKey, orgId)
	if err != nil {
		return false, err
	}

	for _, ancestor := range ancestors {
		if ancestor == childRoleKey {
			return true, nil
		}
	}

	return false, nil
}

// GetInheritDepth 获取角色继承深度
// 例如: A -> B -> C，则 A 的深度为 0，B 的深度为 1，C 的深度为 2
func (s *casbinService) GetInheritDepth(ctx context.Context, roleKey string, orgId int64) (int, error) {
	parents, err := s.GetRoleInherits(ctx, roleKey, orgId)
	if err != nil {
		return 0, err
	}

	if len(parents) == 0 {
		return 0, nil
	}

	// 递归计算最大深度
	maxDepth := 0
	for _, parent := range parents {
		depth, err := s.GetInheritDepth(ctx, parent, orgId)
		if err != nil {
			return 0, err
		}
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth + 1, nil
}

// ReloadPolicy 重新加载策略（从数据库）
func (s *casbinService) ReloadPolicy(ctx context.Context) error {
	if err := s.enforcer.LoadPolicy(); err != nil {
		s.logger.Error("重新加载 Casbin 策略失败", zap.Error(err))
		return fmt.Errorf("重新加载策略失败: %w", err)
	}

	s.logger.Info("重新加载 Casbin 策略成功")
	return nil
}
