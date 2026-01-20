import type { Router, RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { usePermissionStore } from '@/stores/permission';
import { defaultRoutes } from '@/router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

NProgress.configure({ showSpinner: false });

// 白名单路由（不需要登录即可访问）
const whiteList = ['/login', '/404'];

export function setupRouterGuard(router: Router) {
  // 前置守卫
  router.beforeEach(async (to, _from, next) => {
    NProgress.start();

    // 设置页面标题
    document.title = to.meta.title
      ? `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE}`
      : import.meta.env.VITE_APP_TITLE;

    const authStore = useAuthStore();
    const permissionStore = usePermissionStore();

    // 白名单路由（不需要登录即可访问）
    if (whiteList.includes(to.path)) {
      // 如果已登录，访问登录页时跳转到首页
      if (to.path === '/login' && authStore.isLoggedIn) {
        // 确保动态路由已加载
        if (!permissionStore.isRoutesGenerated) {
          try {
            const dynamicRoutes = await permissionStore.generateRoutes();
            const routesToAdd = dynamicRoutes.length > 0 ? dynamicRoutes : defaultRoutes;
            routesToAdd.forEach((route: RouteRecordRaw) => {
              router.addRoute(route);
            });
            permissionStore.setRoutes(routesToAdd);
          } catch (error) {
            console.error('加载动态路由失败:', error);
            defaultRoutes.forEach((route: RouteRecordRaw) => {
              router.addRoute(route);
            });
            permissionStore.setRoutes(defaultRoutes);
            permissionStore.isRoutesGenerated = true;
          }
        }
        next('/dashboard');
        return;
      }
      next();
      return;
    }

    // 需要认证的页面
    if (!authStore.isLoggedIn) {
      next({
        path: '/login',
        query: { redirect: to.fullPath },
      });
      return;
    }

    // 动态路由处理
    if (!permissionStore.isRoutesGenerated) {
      try {
        // 生成动态路由（只在首次或刷新后执行）
        const dynamicRoutes = await permissionStore.generateRoutes();

        // 将动态路由添加为 Root 路由的子路由
        if (dynamicRoutes.length > 0) {
          dynamicRoutes.forEach((route: RouteRecordRaw) => {
            router.addRoute('Root', route);
          });
          permissionStore.setRoutes(dynamicRoutes);
        }

        // 添加 404 路由（放在最后，确保动态路由优先匹配）
        router.addRoute({
          path: '/:pathMatch(.*)*',
          name: 'NotFound',
          component: () => import('@/views/error/404.vue'),
          meta: {
            title: '404',
            requiresAuth: false,
          },
        });

        // 重新导航到目标路由
        next({ ...to, replace: true });
        return;
      } catch (error) {
        console.error('加载动态路由失败:', error);
        
        // 添加 404 路由
        router.addRoute({
          path: '/:pathMatch(.*)*',
          name: 'NotFound',
          component: () => import('@/views/error/404.vue'),
          meta: {
            title: '404',
            requiresAuth: false,
          },
        });
        
        // 重新导航到目标路由
        next({ ...to, replace: true });
        return;
      }
    } else {
      // 路由已生成，但可能需要重新注册（刷新后路由会丢失）
      const currentRoutes = router.getRoutes();
      const hasSystemRoute = currentRoutes.some(r => r.path === '/system');
      
      if (permissionStore.routes.length > 0 && !hasSystemRoute) {
        // 重新注册路由（不调用 API）
        permissionStore.routes.forEach((route: RouteRecordRaw) => {
          router.addRoute('Root', route);
        });
        
        // 添加 404 路由
        router.addRoute({
          path: '/:pathMatch(.*)*',
          name: 'NotFound',
          component: () => import('@/views/error/404.vue'),
          meta: {
            title: '404',
            requiresAuth: false,
          },
        });
        
        // 重新导航到目标路由
        next({ ...to, replace: true });
        return;
      }
    }

    // 检查权限
    if (to.meta.permission) {
      const permission = to.meta.permission as string;
      if (!permissionStore.hasPermission(permission)) {
        console.warn(`没有权限访问: ${to.path}, 需要权限: ${permission}`);
        next({ name: 'NotFound' });
        return;
      }
    }

    next();
  });

  // 后置守卫
  router.afterEach(() => {
    NProgress.done();
  });

  // 错误处理
  router.onError((error) => {
    console.error('Router error:', error);
    NProgress.done();
  });
}
