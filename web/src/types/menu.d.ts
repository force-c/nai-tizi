// 菜单类型定义

// 菜单类型枚举
export type MenuType = 'M' | 'C' | 'F';

// 菜单信息
export interface MenuInfo {
  menuId: number;
  menuName: string;
  parentId: number;
  sortOrder: number;
  path: string;
  component: string;
  query?: string;          // 路由参数
  isFrame?: number;        // 是否外链：0否 1是
  isCache?: number;        // 是否缓存：0不缓存 1缓存
  menuType: MenuType;      // 菜单类型：M目录 C菜单 F按钮
  visible: number;         // 显示状态：0显示 1隐藏
  status: number;          // 状态：0正常 1停用
  perms: string;           // 权限标识
  icon: string;            // 菜单图标
  remark?: string;         // 备注
  createBy?: number;
  createdAt?: string;
  updatedAt?: string;
  children?: MenuInfo[];
}

// 路由元信息
export interface RouteMeta {
  title: string;
  icon?: string;
  permission?: string;
  hidden?: boolean;
  keepAlive?: boolean;
  requiresAuth?: boolean;
}

// 动态路由配置
export interface DynamicRoute {
  path: string;
  name: string;
  component: string;
  redirect?: string;
  meta: RouteMeta;
  children?: DynamicRoute[];
}
