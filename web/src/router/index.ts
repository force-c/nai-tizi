import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { setupRouterGuard } from './guards';

/**
 * 基础路由（静态路由）
 * 这些路由不需要权限控制，所有用户都可以访问
 */
export const basicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login/index.vue'),
    meta: {
      title: '登录',
      requiresAuth: false,
    },
  },
  {
    path: '/',
    name: 'Root',
    component: () => import('@/layouts/default/index.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true,
      hidden: true,
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          icon: 'dashboard',
          requiresAuth: true,
        },
      },
    ],
  },
  // 注意：404 路由不要放在这里，会在动态路由加载后再添加
];

/**
 * 默认路由（用于开发环境或动态路由加载失败时的后备）
 * 这些路由会在动态路由加载失败时使用
 */
export const defaultRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layouts/default/index.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表盘',
          icon: 'dashboard',
          requiresAuth: true,
        },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: basicRoutes,
  scrollBehavior: () => ({ top: 0 }),
});

// 设置路由守卫
setupRouterGuard(router);

export default router;
