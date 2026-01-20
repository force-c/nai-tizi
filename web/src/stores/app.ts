import { defineStore } from 'pinia';

interface AppState {
  collapsed: boolean;
  theme: 'light' | 'dark';
  locale: string;
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    collapsed: false,
    theme: 'light',
    locale: 'zh-CN',
  }),

  actions: {
    toggleCollapsed() {
      this.collapsed = !this.collapsed;
    },

    setCollapsed(collapsed: boolean) {
      this.collapsed = collapsed;
    },

    setTheme(theme: 'light' | 'dark') {
      this.theme = theme;
    },

    setLocale(locale: string) {
      this.locale = locale;
    },
  },

  persist: {
    key: 'app',
    storage: localStorage,
    paths: ['collapsed', 'theme', 'locale'],
  },
});
