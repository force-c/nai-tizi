import { defineStore } from 'pinia';
import { theme } from 'ant-design-vue';

export type ThemeMode = 'light' | 'dark';

interface ThemeState {
  mode: ThemeMode;
}

export const useThemeStore = defineStore('theme', {
  state: (): ThemeState => ({
    mode: 'light',
  }),

  getters: {
    isDark: (state) => state.mode === 'dark',
    
    // Ant Design Vue 主题配置
    antdTheme: (state) => {
      if (state.mode === 'dark') {
        return {
          algorithm: theme.darkAlgorithm,
          token: {
            colorPrimary: '#1890ff',
            colorBgContainer: '#141414',
            colorBgElevated: '#1f1f1f',
            colorBgLayout: '#000000',
            colorBorder: '#424242',
            colorText: 'rgba(255, 255, 255, 0.85)',
            colorTextSecondary: 'rgba(255, 255, 255, 0.65)',
          },
        };
      }
      
      return {
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1890ff',
          colorBgContainer: '#ffffff',
          colorBgElevated: '#ffffff',
          colorBgLayout: '#f0f2f5',
          colorBorder: '#d9d9d9',
          colorText: 'rgba(0, 0, 0, 0.88)',
          colorTextSecondary: 'rgba(0, 0, 0, 0.65)',
        },
      };
    },
    
    // 侧边栏样式
    siderStyle: (state) => {
      if (state.mode === 'dark') {
        return {
          backgroundColor: '#1f1f1f',
          borderRight: '1px solid #303030',
        };
      }
      return {
        backgroundColor: '#ffffff',
        borderRight: '1px solid #f0f0f0',
      };
    },
    
    // Logo 样式
    logoStyle: (state) => {
      if (state.mode === 'dark') {
        return {
          color: '#ffffff',
          backgroundColor: '#1f1f1f',
          borderBottom: '1px solid #303030',
        };
      }
      return {
        color: '#1890ff',
        backgroundColor: '#ffffff',
        borderBottom: '1px solid #f0f0f0',
      };
    },
    
    // Header 样式
    headerStyle: (state) => {
      if (state.mode === 'dark') {
        return {
          backgroundColor: '#141414',
          borderBottom: '1px solid #303030',
          color: '#ffffff',
        };
      }
      return {
        backgroundColor: '#ffffff',
        borderBottom: '1px solid #f0f0f0',
        color: '#000000',
      };
    },
    
    // Content 样式
    contentStyle: (state) => {
      if (state.mode === 'dark') {
        return {
          backgroundColor: '#000000',
        };
      }
      return {
        backgroundColor: '#f0f2f5',
      };
    },
  },

  actions: {
    toggleTheme() {
      this.mode = this.mode === 'light' ? 'dark' : 'light';
    },

    setTheme(mode: ThemeMode) {
      this.mode = mode;
    },

    initTheme() {
      // 从 localStorage 读取主题设置
      const savedTheme = localStorage.getItem('theme-mode') as ThemeMode;
      if (savedTheme) {
        this.mode = savedTheme;
      }
    },
  },

  persist: {
    key: 'theme-mode',
    storage: localStorage,
    paths: ['mode'],
  },
});
