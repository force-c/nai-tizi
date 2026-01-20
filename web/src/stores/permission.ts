import { defineStore } from 'pinia';
import type { RouteRecordRaw } from 'vue-router';
import { menuApi } from '@/api/menu';
import type { MenuInfo } from '@/types/menu';

interface PermissionState {
  routes: RouteRecordRaw[];
  dynamicRoutes: RouteRecordRaw[];
  menuTree: MenuInfo[];
  permissions: string[];
  isRoutesGenerated: boolean;
}

// 布局组件映射
const LayoutMap: Record<string, () => Promise<any>> = {
  Layout: () => import('@/layouts/default/index.vue'),
  BlankLayout: () => import('@/layouts/blank/index.vue'),
};

// 视图组件映射（懒加载）
const viewModules = import.meta.glob('@/views/**/*.vue');

export const usePermissionStore = defineStore('permission', {
  state: (): PermissionState => ({
    routes: [],
    dynamicRoutes: [],
    menuTree: [],
    permissions: [],
    isRoutesGenerated: false, // 会被持久化到 sessionStorage
  }),

  getters: {
    // 获取所有路由
    getRoutes: (state) => state.routes,

    // 获取动态路由
    getDynamicRoutes: (state) => state.dynamicRoutes,

    // 获取菜单树
    getMenuTree: (state) => state.menuTree,

    // 检查是否有权限
    hasPermission: (state) => (permission: string) => {
      if (!permission) return true;
      return (
        state.permissions.includes(permission) ||
        state.permissions.includes('*') ||
        state.permissions.some((p) => {
          // 支持通配符匹配，例如 user.* 匹配 user.create, user.read 等
          if (p.endsWith('.*')) {
            const prefix = p.slice(0, -2);
            return permission.startsWith(prefix + '.');
          }
          return false;
        })
      );
    },
  },

  actions: {
    // 生成路由
    async generateRoutes() {
      try {
        // 从后端获取用户菜单树
        const menuTree = await menuApi.getUserMenuTree();
        
        // 转换后端数据格式到前端格式
        const normalizedMenuTree = this.normalizeMenuData(menuTree);
        this.menuTree = normalizedMenuTree;

        // 提取所有权限标识
        this.permissions = this.extractPermissions(normalizedMenuTree);

        // 将菜单树转换为路由配置
        const dynamicRoutes = this.transformMenuToRoutes(normalizedMenuTree);
        this.dynamicRoutes = dynamicRoutes;

        // 标记路由已生成
        this.isRoutesGenerated = true;

        return dynamicRoutes;
      } catch (error) {
        console.error('生成路由失败:', error);
        throw error;
      }
    },

    // 规范化菜单数据（将后端字段映射到前端字段）
    normalizeMenuData(menus: any[]): MenuInfo[] {
      return menus.map(menu => {
        const normalized: any = {
          menuId: menu.id || menu.menuId,
          menuName: menu.menuName,
          parentId: menu.parentId,
          sortOrder: menu.sort || menu.sortOrder || 0,
          path: menu.path,
          component: menu.component,
          query: menu.query,
          isFrame: String(menu.isFrame || 0),
          isCache: String(menu.isCache || 0),
          // 转换 menuType: 0->M(目录), 1->C(菜单), 2->F(按钮)
          menuType: menu.menuType === 0 ? 'M' : menu.menuType === 1 ? 'C' : menu.menuType === 2 ? 'F' : menu.menuType,
          visible: String(menu.visible),
          status: String(menu.status),
          perms: menu.perms || '',
          icon: menu.icon || '',
          remark: menu.remark,
          createBy: menu.createBy,
          createdAt: menu.createdTime || menu.createdAt,
          updatedAt: menu.updatedTime || menu.updatedAt,
        };

        // 递归处理子菜单
        if (menu.children && menu.children.length > 0) {
          normalized.children = this.normalizeMenuData(menu.children);
        }

        return normalized;
      });
    },

    // 提取权限标识
    extractPermissions(menus: MenuInfo[]): string[] {
      const permissions: string[] = [];

      const extract = (menuList: MenuInfo[]) => {
        menuList.forEach((menu) => {
          if (menu.perms) {
            permissions.push(menu.perms);
          }
          if (menu.children && menu.children.length > 0) {
            extract(menu.children);
          }
        });
      };

      extract(menus);
      return permissions;
    },

    // 将菜单树转换为路由配置
    transformMenuToRoutes(menus: MenuInfo[]): RouteRecordRaw[] {
      const routes: RouteRecordRaw[] = [];

      menus.forEach((menu) => {
        // 只处理目录(M)和菜单(C)，按钮(F)不生成路由
        if (menu.menuType === 'M' || menu.menuType === 'C') {
          const route = this.createRoute(menu);
          if (route) {
            // 如果是顶级菜单（目录），需要作为 Root 路由的子路由
            if (menu.parentId === 0 && menu.menuType === 'M') {
              // 顶级目录，添加到 Root 路由的 children
              routes.push(route);
            } else {
              routes.push(route);
            }
          }
        }
      });

      return routes;
    },

    // 创建单个路由
    createRoute(menu: MenuInfo): RouteRecordRaw | null {
      // 如果菜单不可见，跳过
      if (menu.visible === '1') {
        return null;
      }

      const route: Partial<RouteRecordRaw> = {
        path: menu.path,
        name: this.generateRouteName(menu),
        meta: {
          title: menu.menuName,
          icon: menu.icon,
          permission: menu.perms,
          hidden: menu.visible === '1',
          requiresAuth: true,
        },
      };

      // 处理组件
      if (menu.menuType === 'M') {
        // 目录类型，使用布局组件
        // 如果 component 为空或为 Layout，使用默认布局
        if (!menu.component || menu.component === '' || menu.component === 'Layout') {
          route.component = LayoutMap.Layout;
        } else {
          route.component = LayoutMap[menu.component] || LayoutMap.Layout;
        }
        
        // 如果有子菜单，添加重定向到第一个可见子菜单
        if (menu.children && menu.children.length > 0) {
          const firstVisibleChild = menu.children.find(child => child.visible === '0');
          if (firstVisibleChild) {
            // 构建完整的重定向路径
            route.redirect = `${menu.path}/${firstVisibleChild.path}`;
          }
        }
      } else if (menu.menuType === 'C') {
        // 菜单类型，使用视图组件
        if (menu.component) {
          route.component = this.loadViewComponent(menu.component);
        }
      }

      // 处理子菜单
      if (menu.children && menu.children.length > 0) {
        route.children = this.transformMenuToRoutes(menu.children);
      }

      return route as RouteRecordRaw;
    },

    // 生成路由名称
    generateRouteName(menu: MenuInfo): string {
      // 将路径转换为驼峰命名，例如 /system/user -> SystemUser
      const pathParts = menu.path.split('/').filter(Boolean);
      return pathParts
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join('');
    },

    // 加载视图组件
    loadViewComponent(component: string) {
      // 确保组件路径以 /src/views/ 开头
      let componentPath = component;
      if (!componentPath.startsWith('/src/views/')) {
        componentPath = `/src/views/${component}`;
      }
      if (!componentPath.endsWith('.vue')) {
        componentPath = `${componentPath}.vue`;
      }

      // 使用 Vite 的 glob import
      const matchedComponent = viewModules[componentPath];

      if (matchedComponent) {
        return matchedComponent;
      }

      // 如果找不到组件，返回 404 页面
      console.warn(`组件 ${componentPath} 不存在，使用 404 页面`);
      return () => import('@/views/error/404.vue');
    },

    // 重置路由
    resetRoutes() {
      this.routes = [];
      this.dynamicRoutes = [];
      this.menuTree = [];
      this.permissions = [];
      this.isRoutesGenerated = false;
    },

    // 设置路由
    setRoutes(routes: RouteRecordRaw[]) {
      this.routes = routes;
    },
  },

  persist: {
    key: 'permission',
    storage: sessionStorage, // 使用 sessionStorage，关闭浏览器后清除
    paths: ['permissions', 'menuTree', 'isRoutesGenerated', 'routes'], // 持久化路由配置
  },
});
