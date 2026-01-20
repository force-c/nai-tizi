import { defineStore } from 'pinia';

interface UserState {
  permissions: string[];
  roles: string[];
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    permissions: [],
    roles: [],
  }),

  getters: {
    hasPermission: (state) => (permission: string) => {
      return state.permissions.includes(permission) || state.permissions.includes('*');
    },

    hasRole: (state) => (role: string) => {
      return state.roles.includes(role);
    },

    hasAnyPermission: (state) => (permissions: string[]) => {
      return permissions.some((permission) => state.permissions.includes(permission));
    },
  },

  actions: {
    setPermissions(permissions: string[]) {
      this.permissions = permissions;
    },

    setRoles(roles: string[]) {
      this.roles = roles;
    },

    clearPermissions() {
      this.permissions = [];
      this.roles = [];
    },
  },

  persist: {
    key: 'user',
    storage: localStorage,
    paths: ['permissions', 'roles'],
  },
});
