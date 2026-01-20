<template>
  <a-layout class="layout-container">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      class="layout-sider"
      :style="themeStore.siderStyle"
    >
      <div class="logo" :style="themeStore.logoStyle">
        <h1 v-if="!collapsed">{{ title }}</h1>
        <h1 v-else>NAI</h1>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        :theme="themeStore.isDark ? 'dark' : 'light'"
        mode="inline"
        @click="handleMenuClick"
      >
        <template v-for="menu in menuList" :key="menu.path">
          <a-sub-menu v-if="menu.children && menu.children.length > 0" :key="`sub-${menu.path}`">
            <template #title>
              <component :is="getIcon(menu.icon)" v-if="menu.icon" />
              <span>{{ menu.title }}</span>
            </template>
            <a-menu-item
              v-for="child in menu.children"
              :key="child.path"
            >
              <component :is="getIcon(child.icon)" v-if="child.icon" />
              <span>{{ child.title }}</span>
            </a-menu-item>
          </a-sub-menu>
          <a-menu-item v-else :key="`item-${menu.path}`">
            <component :is="getIcon(menu.icon)" v-if="menu.icon" />
            <span>{{ menu.title }}</span>
          </a-menu-item>
        </template>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="layout-header" :style="themeStore.headerStyle">
        <div class="header-left">
          <menu-unfold-outlined
            v-if="collapsed"
            class="trigger"
            @click="toggleCollapsed"
          />
          <menu-fold-outlined v-else class="trigger" @click="toggleCollapsed" />
        </div>

        <div class="header-right">
          <a-tooltip :title="themeStore.isDark ? '切换到日间模式' : '切换到夜间模式'">
            <a-button
              type="text"
              :icon="themeStore.isDark ? h(BulbFilled) : h(BulbOutlined)"
              @click="themeStore.toggleTheme"
              class="theme-toggle"
            />
          </a-tooltip>
          
          <a-dropdown>
            <div class="user-info">
              <a-avatar :size="32">
                <template #icon>
                  <user-outlined />
                </template>
              </a-avatar>
              <span class="username">{{ authStore.nickname || authStore.username }}</span>
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="handleLogout">
                  <logout-outlined />
                  退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="layout-content" :style="themeStore.contentStyle">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue';
import { useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
  BulbOutlined,
  BulbFilled,
} from '@ant-design/icons-vue';
import { useAuthStore } from '@/stores/auth';
import { useAppStore } from '@/stores/app';
import { usePermissionStore } from '@/stores/permission';
import { useThemeStore } from '@/stores/theme';
import { getIcon } from '@/config/icons';

const router = useRouter();
const authStore = useAuthStore();
const appStore = useAppStore();
const permissionStore = usePermissionStore();
const themeStore = useThemeStore();

const title = import.meta.env.VITE_APP_TITLE;
const collapsed = ref(appStore.collapsed);
const selectedKeys = ref([router.currentRoute.value.path]);

// 从路由生成菜单列表
const menuList = computed(() => {
  // 优先使用 menuTree，如果没有则使用 routes
  const menuTree = permissionStore.getMenuTree;
  
  if (menuTree && menuTree.length > 0) {
    return buildMenuFromTree(menuTree);
  }
  
  const routes = permissionStore.getRoutes;
  const menus: any[] = [];

  routes.forEach((route) => {
    // 跳过隐藏的路由
    if (route.meta?.hidden) return;

    const menu: any = {
      path: route.path,
      title: route.meta?.title || route.name,
      icon: route.meta?.icon,
      children: [],
    };

    // 处理子路由
    if (route.children && route.children.length > 0) {
      route.children.forEach((child) => {
        if (child.meta?.hidden) return;
        
        menu.children.push({
          path: child.path.startsWith('/') ? child.path : `${route.path}/${child.path}`,
          title: child.meta?.title || child.name,
          icon: child.meta?.icon,
        });
      });
    }

    menus.push(menu);
  });

  return menus;
});

// 从菜单树构建菜单列表
const buildMenuFromTree = (menuTree: any[]) => {
  const menus: any[] = [];
  
  menuTree.forEach((menu) => {
    // 跳过隐藏的菜单和按钮
    if (menu.visible === '1' || menu.menuType === 'F') return;
    
    const menuItem: any = {
      path: menu.path,
      title: menu.menuName,
      icon: menu.icon,
      children: [],
    };
    
    // 处理子菜单
    if (menu.children && menu.children.length > 0) {
      menu.children.forEach((child: any) => {
        // 跳过隐藏的菜单和按钮
        if (child.visible === '1' || child.menuType === 'F') return;
        
        menuItem.children.push({
          path: child.path.startsWith('/') ? child.path : `${menu.path}/${child.path}`,
          title: child.menuName,
          icon: child.icon,
        });
      });
    }
    
    menus.push(menuItem);
  });
  
  return menus;
};

const toggleCollapsed = () => {
  collapsed.value = !collapsed.value;
  appStore.setCollapsed(collapsed.value);
};

const handleMenuClick = ({ key }: any) => {
  router.push(key as string);
  selectedKeys.value = [key as string];
};

const handleLogout = async () => {
  try {
    await authStore.logout();
    message.success('退出登录成功');
    router.push('/login');
  } catch (error: any) {
    message.error(error.message || '退出登录失败');
  }
};
</script>

<style scoped>
.layout-container {
  width: 100%;
  height: 100%;
}

.layout-sider {
  overflow: auto;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 600;
  transition: all 0.3s;
}

.logo h1 {
  margin: 0;
  font-size: 18px;
  transition: color 0.3s;
}

.layout-header {
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: fixed;
  top: 0;
  right: 0;
  left: 200px;
  z-index: 9;
  height: 64px;
  transition: all 0.3s;
}

.layout-sider.ant-layout-sider-collapsed + .ant-layout .layout-header {
  left: 80px;
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.theme-toggle {
  font-size: 18px;
  transition: all 0.3s;
}

.theme-toggle:hover {
  transform: rotate(20deg);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
}

.layout-content {
  margin: 20px 0 0 90px;
  padding: 24px;
  min-height: calc(100vh - 64px);
  transition: all 0.3s;
}

.layout-sider.ant-layout-sider-collapsed + .ant-layout .layout-content {
  margin-left: 80px;
}
</style>
