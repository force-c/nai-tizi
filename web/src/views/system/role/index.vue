<template>
  <div class="role-page">
    <BasicTable
      :columns="columns"
      :api="loadData"
      :use-search-form="true"
      :form-config="searchFormConfig"
      :show-action-column="true"
      @register="registerTable"
    >
      <template #toolbar>
        <a-space>
          <a-button v-permission="'role.create'" type="primary" size="middle" @click="handleCreate">
            <template #icon><PlusOutlined /></template>
            新增
          </a-button>
          <a-button v-permission="'role.delete'" danger size="middle" @click="handleBatchDelete">
            <template #icon><DeleteOutlined /></template>
            批量删除
          </a-button>
        </a-space>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status'">
          <a-tag :color="record.status === 0 ? 'success' : 'error'">
            {{ record.status === 0 ? '正常' : '停用' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <TableAction
            :actions="[
              {
                label: '编辑',
                onClick: () => handleEdit(record),
                ifShow: hasPermission('role.update'),
              },
              {
                label: '权限',
                onClick: () => handlePermission(record),
                ifShow: hasPermission('role.permission'),
              },
              {
                label: '删除',
                color: 'error',
                popConfirm: {
                  title: '确定删除该角色吗？',
                  onConfirm: () => handleDelete(record.id),
                },
                ifShow: hasPermission('role.delete'),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>

    <RoleModal
      v-model:visible="modalVisible"
      :role-id="currentRoleId"
      @success="handleSuccess"
    />

    <PermissionModal
      v-model:visible="permissionModalVisible"
      :role-id="currentRoleId"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue';
import BasicTable from '@/components/Table/BasicTable.vue';
import TableAction from '@/components/Table/TableAction.vue';
import RoleModal from './RoleModal.vue';
import PermissionModal from './PermissionModal.vue';
import { roleApi } from '@/api/role';
import { usePermission } from '@/composables/usePermission';
import type { FormConfig } from '@/types/form';

const { hasPermission } = usePermission();

// 表格列配置
const columns = [
  {
    title: '角色ID',
    dataIndex: 'id',
    width: 80,
  },
  {
    title: '角色名称',
    dataIndex: 'roleName',
    width: 150,
  },
  {
    title: '角色标识',
    dataIndex: 'roleKey',
    width: 150,
  },
  {
    title: '显示顺序',
    dataIndex: 'sort',
    width: 100,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
  },
  {
    title: '创建时间',
    dataIndex: 'createdTime',
    width: 180,
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 200,
    fixed: 'right',
  },
];

// 搜索表单配置
const searchFormConfig: FormConfig = {
  schemas: [
    {
      field: 'roleName',
      label: '角色名称',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'roleKey',
      label: '角色标识',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'status',
      label: '状态',
      component: 'Select',
      componentProps: {
        options: [
          { label: '正常', value: 0 },
          { label: '停用', value: 1 },
        ],
      },
      colProps: { span: 6 },
    },
  ],
};

// 表格实例
const tableRef = ref();
const registerTable = (methods: any) => {
  tableRef.value = methods;
};

// 加载数据
const loadData = async (params: any) => {
  return await roleApi.list(params);
};

// 弹窗状态
const modalVisible = ref(false);
const permissionModalVisible = ref(false);
const currentRoleId = ref<number>();

// 新增角色
const handleCreate = () => {
  currentRoleId.value = undefined;
  modalVisible.value = true;
};

// 编辑角色
const handleEdit = (record: any) => {
  currentRoleId.value = record.id;
  modalVisible.value = true;
};

// 权限分配
const handlePermission = (record: any) => {
  currentRoleId.value = record.id;
  permissionModalVisible.value = true;
};

// 删除角色
const handleDelete = async (roleId: number) => {
  try {
    await roleApi.delete(roleId);
    message.success('删除成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 批量删除
const handleBatchDelete = async () => {
  const selectedRows = tableRef.value?.getSelectRows();
  if (!selectedRows || selectedRows.length === 0) {
    message.warning('请选择要删除的角色');
    return;
  }
  try {
    const roleIds = selectedRows.map((row: any) => row.id);
    await roleApi.batchDelete(roleIds);
    message.success('批量删除成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('批量删除失败:', error);
  }
};

// 操作成功回调
const handleSuccess = () => {
  modalVisible.value = false;
  permissionModalVisible.value = false;
  tableRef.value?.reload();
};
</script>

<style scoped lang="less">
.role-page {
  background: #fff;
  padding: 24px;
  border-radius: 4px;
}
</style>
