<template>
  <BasicModal
    v-model:visible="visible"
    title="分配权限"
    :width="600"
    :confirm-loading="loading"
    @ok="handleSubmit"
  >
    <a-spin :spinning="loading">
      <a-tree
        v-model:checkedKeys="checkedKeys"
        :tree-data="menuTree"
        :field-names="{ title: 'menuName', key: 'menuId', children: 'children' }"
        checkable
        :selectable="false"
        default-expand-all
      />
    </a-spin>
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { message } from 'ant-design-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import { roleApi } from '@/api/role';
import { menuApi } from '@/api/menu';

const props = defineProps<{
  visible: boolean;
  roleId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const loading = ref(false);
const menuTree = ref<any[]>([]);
const checkedKeys = ref<number[]>([]);

// 加载菜单树
const loadMenuTree = async () => {
  try {
    loading.value = true;
    const data = await menuApi.getMenuTree();
    menuTree.value = data;
  } catch (error) {
    console.error('加载菜单树失败:', error);
  } finally {
    loading.value = false;
  }
};

// 加载角色权限
const loadRoleMenus = async () => {
  if (!props.roleId) return;
  
  try {
    loading.value = true;
    const menuIds = await roleApi.getMenus(props.roleId);
    checkedKeys.value = menuIds;
  } catch (error) {
    console.error('加载角色权限失败:', error);
  } finally {
    loading.value = false;
  }
};

// 提交
const handleSubmit = async () => {
  if (!props.roleId) return;
  
  try {
    loading.value = true;
    await roleApi.assignMenus(props.roleId, checkedKeys.value);
    message.success('权限分配成功');
    emit('success');
  } catch (error) {
    console.error('权限分配失败:', error);
  } finally {
    loading.value = false;
  }
};

// 监听弹窗显示
watch(
  () => props.visible,
  (val) => {
    if (val) {
      loadMenuTree();
      loadRoleMenus();
    } else {
      checkedKeys.value = [];
    }
  }
);
</script>
