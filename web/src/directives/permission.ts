import type { Directive, DirectiveBinding } from 'vue';
import { usePermissionStore } from '@/stores/permission';

export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding;
    
    if (!value) return;

    const permissionStore = usePermissionStore();
    const permissions = Array.isArray(value) ? value : [value];
    
    const hasPermission = permissions.some((perm) => permissionStore.hasPermission(perm));
    
    if (!hasPermission) {
      el.style.display = 'none';
      // 或者直接移除元素
      // el.parentNode?.removeChild(el);
    }
  },
};
