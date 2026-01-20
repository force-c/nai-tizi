import { usePermissionStore } from '@/stores/permission';

export function usePermission() {
  const permissionStore = usePermissionStore();

  const hasPermission = (permission: string | string[]): boolean => {
    if (!permission) return true;

    const permissions = Array.isArray(permission) ? permission : [permission];
    
    return permissions.some((perm) => permissionStore.hasPermission(perm));
  };

  return {
    hasPermission,
  };
}
