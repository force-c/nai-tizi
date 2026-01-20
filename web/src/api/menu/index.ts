import { request } from '@/utils/request';
import type { MenuInfo } from '@/types/menu';

export const menuApi = {
  // 获取当前用户的菜单树（用于生成路由和侧边栏）
  getUserMenuTree: () => request.get<MenuInfo[]>('/api/v1/menu/user/tree'),

  // 获取所有菜单列表（用于菜单管理页面）
  getMenuList: (params?: any) => request.get<MenuInfo[]>('/api/v1/menu', { params }),

  // 获取菜单树（用于菜单管理页面）
  getMenuTree: () => request.get<MenuInfo[]>('/api/v1/menu/tree'),

  // 获取菜单详情
  detail: (menuId: number) => request.get<MenuInfo>(`/api/v1/menu/${menuId}`),

  // 创建菜单
  create: (data: Partial<MenuInfo>) => request.post('/api/v1/menu', data),

  // 更新菜单
  update: (menuId: number, data: Partial<MenuInfo>) =>
    request.put(`/api/v1/menu/${menuId}`, data),

  // 删除菜单
  delete: (menuId: number) => request.delete(`/api/v1/menu/${menuId}`),
};
