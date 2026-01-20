package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/force-c/nai-tizi/internal/config"
	"github.com/force-c/nai-tizi/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CasbinServiceV2 Casbin 权限管理服务接口（支持单一企业和多租户）
type CasbinServiceV2 interface {
	// CheckPermission 检查用户权限（自动适配单租户/多租户）
	// 单一企业模式：CheckPermission(ctx, userId, resource, action)
	// 多租户模式：从 ctx 中获取 tenantId
	CheckPermission(ctx context.Context, userId int64, resource, action string) (bool, error)

	// AddRoleForUser 为用户分配角色（自动适配）
	AddRoleForUser(ctx context.Context, userId int64, roleKey string) error

	// DeleteRoleForUser 删除用户的角色（自动适配）
	DeleteRoleForUser(ctx context.Context, userId int64, roleKey string) error

	// GetRolesForUser 获取用户的所有角色（自动适配）
	GetRolesForUser(ctx context.Context, userId int64) ([]string, error)

	// GetUsersForRole 获取拥有指定角色的所有用户（自动适配）
	GetUsersForRole(ctx context.Context, roleKey string) ([]int64, error)

	// AddPermissionForRole 为角色添加权限（自动适配）
	// resource: 资源路径（支持通配符，例如: "user.*", "*.read", "*"）
	// action: 操作类型（支持通配符，例如: "write", "*"）
	AddPermissionForRole(ctx context.Context, roleKey string, resource, action string) error

	// DeletePermissionForRole 删除角色的权限（自动适配）
	DeletePermissionForRole(ctx context.Context, roleKey string, resource, action string) error

	// GetPermissionsForRole 获取角色的所有权限（自动适配）
	GetPermissionsForRole(ctx context.Context, roleKey string) ([][]string, error)

	// SetRoleInherit 设置角色继承关系（自动适配）
	SetRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string) error

	// DeleteRoleInherit 删除角色继承关系（自动适配）
	DeleteRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string) error

	// GetRoleInherits 获取角色的所有父角色（自动适配）
	GetRoleInherits(ctx context.Context, roleKey string) ([]string, error)

	// CheckCircularInherit 检查是否存在循环继承（自动适配）
	CheckCircularInherit(ctx context.Context, childRoleKey, parentRoleKey string) (bool, error)

	// GetInheritDepth 获取角色继承深度（自动适配）
	GetInheritDepth(ctx context.Context, roleKey string) (int, error)

	// ReloadPolicy 重新加载策略（从数据库）
	ReloadPolicy(ctx context.Context) error
}

type casbinServiceV2 struct {
	enforcer           *casbin.Enforcer
	db                 *gorm.DB
	logger             logger.Logger
	multiTenantEnabled bool // 多租户开关
}

// NewCasbinServiceV2 创建 Casbin 服务实例（支持单一企业和多租户）
func NewCasbinServiceV2(enforcer *casbin.Enforcer, db *gorm.DB, logger logger.Logger, cfg *config.Config) CasbinServiceV2 {
	return &casbinServiceV2{
		enforcer:           enforcer,
		db:                 db,
		logger:             logger,
		multiTenantEnabled: cfg.MultiTenant.Enabled,
	}
}

// getTenantId 从 context 获取租户ID（多租户模式）
func (s *casbinServiceV2) getTenantId(ctx context.Context) int64 {
	if !s.multiTenantEnabled {
		return 1 // 单一企业模式，默认租户ID为1
	}

	// 从 context 获取租户ID
	if tenantId, ok := ctx.Value("tenantId").(int64); ok {
		return tenantId
	}

	// 如果没有租户ID，返回默认值1
	s.logger.Warn("租户ID未设置，使用默认值1")
	return 1
}

// CheckPermission 检查用户权限（自动适配单租户/多租户）
func (s *casbinServiceV2) CheckPermission(ctx context.Context, userId int64, resource, action string) (bool, error) {
	sub := fmt.Sprintf("user::%d", userId)

	var ok bool
	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		ok, err = s.enforcer.Enforce(sub, dom, resource, action)

		s.logger.Debug("权限检查（多租户）",
			zap.Int64("userId", userId),
			zap.Int64("tenantId", tenantId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Bool("allowed", ok))
	} else {
		// 单一企业模式
		ok, err = s.enforcer.Enforce(sub, resource, action)

		s.logger.Debug("权限检查（单一企业）",
			zap.Int64("userId", userId),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Bool("allowed", ok))
	}

	if err != nil {
		s.logger.Error("权限检查失败", zap.Error(err))
		return false, fmt.Errorf("权限检查失败: %w", err)
	}

	return ok, nil
}

// AddRoleForUser 为用户分配角色（自动适配）
func (s *casbinServiceV2) AddRoleForUser(ctx context.Context, userId int64, roleKey string) error {
	sub := fmt.Sprintf("user::%d", userId)
	role := fmt.Sprintf("role::%s", roleKey)

	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.AddGroupingPolicy(sub, role, dom)

		s.logger.Info("添加用户角色（多租户）",
			zap.Int64("userId", userId),
			zap.String("roleKey", roleKey),
			zap.Int64("tenantId", tenantId))
	} else {
		// 单一企业模式
		_, err = s.enforcer.AddGroupingPolicy(sub, role)

		s.logger.Info("添加用户角色（单一企业）",
			zap.Int64("userId", userId),
			zap.String("roleKey", roleKey))
	}

	if err != nil {
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	return nil
}

// DeleteRoleForUser 删除用户的角色（自动适配）
func (s *casbinServiceV2) DeleteRoleForUser(ctx context.Context, userId int64, roleKey string) error {
	sub := fmt.Sprintf("user::%d", userId)
	role := fmt.Sprintf("role::%s", roleKey)

	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.RemoveGroupingPolicy(sub, role, dom)
	} else {
		// 单一企业模式
		_, err = s.enforcer.RemoveGroupingPolicy(sub, role)
	}

	if err != nil {
		return fmt.Errorf("删除用户角色失败: %w", err)
	}

	return nil
}

// GetRolesForUser 获取用户的所有角色（自动适配）
func (s *casbinServiceV2) GetRolesForUser(ctx context.Context, userId int64) ([]string, error) {
	sub := fmt.Sprintf("user::%d", userId)

	var roles []string
	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		roles = s.enforcer.GetRolesForUserInDomain(sub, dom)
	} else {
		// 单一企业模式
		roles, err = s.enforcer.GetRolesForUser(sub)
		if err != nil {
			return nil, fmt.Errorf("获取用户角色失败: %w", err)
		}
	}

	// 去除 "role::" 前缀
	var result []string
	for _, role := range roles {
		if len(role) > 6 && role[:6] == "role::" {
			result = append(result, role[6:])
		}
	}

	return result, nil
}

// GetUsersForRole 获取拥有指定角色的所有用户（自动适配）
func (s *casbinServiceV2) GetUsersForRole(ctx context.Context, roleKey string) ([]int64, error) {
	role := fmt.Sprintf("role::%s", roleKey)

	var users []string
	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		users = s.enforcer.GetUsersForRoleInDomain(role, dom)
	} else {
		// 单一企业模式
		users, err = s.enforcer.GetUsersForRole(role)
		if err != nil {
			return nil, fmt.Errorf("获取角色用户失败: %w", err)
		}
	}

	// 解析用户ID
	var result []int64
	for _, user := range users {
		if len(user) > 6 && user[:6] == "user::" {
			var userId int64
			if _, err := fmt.Sscanf(user[6:], "%d", &userId); err == nil {
				result = append(result, userId)
			}
		}
	}

	return result, nil
}

// AddPermissionForRole 为角色添加权限（自动适配）
func (s *casbinServiceV2) AddPermissionForRole(ctx context.Context, roleKey string, resource, action string) error {
	sub := fmt.Sprintf("role::%s", roleKey)

	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.AddPolicy(sub, dom, resource, action)
	} else {
		// 单一企业模式
		_, err = s.enforcer.AddPolicy(sub, resource, action)
	}

	if err != nil {
		return fmt.Errorf("添加角色权限失败: %w", err)
	}

	s.logger.Info("添加角色权限",
		zap.String("roleKey", roleKey),
		zap.String("resource", resource),
		zap.String("action", action))

	return nil
}

// DeletePermissionForRole 删除角色的权限（自动适配）
func (s *casbinServiceV2) DeletePermissionForRole(ctx context.Context, roleKey string, resource, action string) error {
	sub := fmt.Sprintf("role::%s", roleKey)

	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.RemovePolicy(sub, dom, resource, action)
	} else {
		// 单一企业模式
		_, err = s.enforcer.RemovePolicy(sub, resource, action)
	}

	if err != nil {
		return fmt.Errorf("删除角色权限失败: %w", err)
	}

	return nil
}

// GetPermissionsForRole 获取角色的所有权限（自动适配）
func (s *casbinServiceV2) GetPermissionsForRole(ctx context.Context, roleKey string) ([][]string, error) {
	sub := fmt.Sprintf("role::%s", roleKey)

	var permissions [][]string
	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		permissions = s.enforcer.GetPermissionsForUserInDomain(sub, dom)
	} else {
		// 单一企业模式
		permissions, err = s.enforcer.GetPermissionsForUser(sub)
		if err != nil {
			return nil, fmt.Errorf("获取角色权限失败: %w", err)
		}
	}

	return permissions, nil
}

// SetRoleInherit 设置角色继承关系（自动适配）
func (s *casbinServiceV2) SetRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string) error {
	// 检查循环继承
	circular, err := s.CheckCircularInherit(ctx, childRoleKey, parentRoleKey)
	if err != nil {
		return err
	}
	if circular {
		return errors.New("检测到循环继承，无法设置")
	}

	// 检查继承深度
	depth, err := s.GetInheritDepth(ctx, childRoleKey)
	if err != nil {
		return err
	}
	if depth >= 3 {
		return errors.New("继承深度超过限制（最多3层）")
	}

	child := fmt.Sprintf("role::%s", childRoleKey)
	parent := fmt.Sprintf("role::%s", parentRoleKey)

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.AddGroupingPolicy(child, parent, dom)
	} else {
		// 单一企业模式
		_, err = s.enforcer.AddGroupingPolicy(child, parent)
	}

	if err != nil {
		return fmt.Errorf("设置角色继承失败: %w", err)
	}

	s.logger.Info("设置角色继承",
		zap.String("childRole", childRoleKey),
		zap.String("parentRole", parentRoleKey))

	return nil
}

// DeleteRoleInherit 删除角色继承关系（自动适配）
func (s *casbinServiceV2) DeleteRoleInherit(ctx context.Context, childRoleKey, parentRoleKey string) error {
	child := fmt.Sprintf("role::%s", childRoleKey)
	parent := fmt.Sprintf("role::%s", parentRoleKey)

	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		_, err = s.enforcer.RemoveGroupingPolicy(child, parent, dom)
	} else {
		// 单一企业模式
		_, err = s.enforcer.RemoveGroupingPolicy(child, parent)
	}

	if err != nil {
		return fmt.Errorf("删除角色继承失败: %w", err)
	}

	return nil
}

// GetRoleInherits 获取角色的所有父角色（自动适配）
func (s *casbinServiceV2) GetRoleInherits(ctx context.Context, roleKey string) ([]string, error) {
	role := fmt.Sprintf("role::%s", roleKey)

	var parents []string
	var err error

	if s.multiTenantEnabled {
		// 多租户模式
		tenantId := s.getTenantId(ctx)
		dom := fmt.Sprintf("tenant::%d", tenantId)
		parents = s.enforcer.GetRolesForUserInDomain(role, dom)
	} else {
		// 单一企业模式
		parents, err = s.enforcer.GetRolesForUser(role)
		if err != nil {
			return nil, fmt.Errorf("获取角色继承失败: %w", err)
		}
	}

	// 去除 "role::" 前缀
	var result []string
	for _, parent := range parents {
		if len(parent) > 6 && parent[:6] == "role::" {
			result = append(result, parent[6:])
		}
	}

	return result, nil
}

// CheckCircularInherit 检查是否存在循环继承（自动适配）
func (s *casbinServiceV2) CheckCircularInherit(ctx context.Context, childRoleKey, parentRoleKey string) (bool, error) {
	// 如果父角色的祖先中包含子角色，则存在循环
	ancestors, err := s.GetRoleInherits(ctx, parentRoleKey)
	if err != nil {
		return false, err
	}

	for _, ancestor := range ancestors {
		if ancestor == childRoleKey {
			return true, nil
		}
		// 递归检查
		circular, err := s.CheckCircularInherit(ctx, childRoleKey, ancestor)
		if err != nil {
			return false, err
		}
		if circular {
			return true, nil
		}
	}

	return false, nil
}

// GetInheritDepth 获取角色继承深度（自动适配）
func (s *casbinServiceV2) GetInheritDepth(ctx context.Context, roleKey string) (int, error) {
	parents, err := s.GetRoleInherits(ctx, roleKey)
	if err != nil {
		return 0, err
	}

	if len(parents) == 0 {
		return 0, nil
	}

	maxDepth := 0
	for _, parent := range parents {
		depth, err := s.GetInheritDepth(ctx, parent)
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
func (s *casbinServiceV2) ReloadPolicy(ctx context.Context) error {
	if err := s.enforcer.LoadPolicy(); err != nil {
		s.logger.Error("重新加载策略失败", zap.Error(err))
		return fmt.Errorf("重新加载策略失败: %w", err)
	}

	s.logger.Info("重新加载策略成功")
	return nil
}
