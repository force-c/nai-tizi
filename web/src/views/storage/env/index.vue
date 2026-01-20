<template>
  <div class="storage-env-page">
    <BasicTable
      :columns="columns"
      :api="loadData"
      :use-search-form="true"
      :form-config="searchFormConfig"
      :show-action-column="true"
      @register="registerTable"
    >
      <template #toolbar>
        <a-button v-permission="'storage_env.create'" type="primary" @click="handleCreate">
          <template #icon><PlusOutlined /></template>
          新增
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'storageType'">
          <a-tag v-if="record.storageType === 'local'" color="blue">本地存储</a-tag>
          <a-tag v-else-if="record.storageType === 'minio'" color="green">MinIO</a-tag>
          <a-tag v-else-if="record.storageType === 's3'" color="orange">S3</a-tag>
          <a-tag v-else-if="record.storageType === 'oss'" color="purple">阿里云OSS</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'config'">
          <div 
            class="json-content clickable" 
            @click="handleViewConfig(record.config)"
          >
            <span v-if="record.config" class="json-preview">
              {{ getJsonPreview(record.config) }}
            </span>
            <span v-else class="text-gray">-</span>
          </div>
        </template>
        <template v-else-if="column.dataIndex === 'isDefault'">
          <a-tag :color="record.isDefault === '1' ? 'success' : 'default'">
            {{ record.isDefault === '1' ? '是' : '否' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
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
                ifShow: hasPermission('storage_env.update'),
              },
              {
                label: '设为默认',
                onClick: () => handleSetDefault(getEnvId(record)),
                ifShow: hasPermission('storage_env.update') && record.isDefault !== '1',
              },
              {
                label: '测试连接',
                onClick: () => handleTest(getEnvId(record)),
                ifShow: hasPermission('storage_env.update'),
              },
              {
                label: '删除',
                color: 'error',
                popConfirm: {
                  title: '确定删除该存储环境吗？',
                  confirm: () => handleDelete(getEnvId(record)),
                },
                ifShow: hasPermission('storage_env.delete'),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>

    <EnvModal
      v-model:visible="modalVisible"
      :env-id="currentEnvId"
      @success="handleSuccess"
    />

    <EnvViewModal
      v-model:visible="viewModalVisible"
      :config="viewConfig"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { PlusOutlined } from '@ant-design/icons-vue';
import BasicTable from '@/components/Table/BasicTable.vue';
import TableAction from '@/components/Table/TableAction.vue';
import EnvModal from './EnvModal.vue';
import EnvViewModal from './EnvViewModal.vue';
import { storageEnvApi } from '@/api/storage';
import { usePermission } from '@/composables/usePermission';

const { hasPermission } = usePermission();

const getEnvId = (record: any) => record.id;

// 表格列配置
const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    width: 150,
  },
  {
    title: '编码',
    dataIndex: 'code',
    width: 120,
  },
  {
    title: '存储类型',
    dataIndex: 'storageType',
    width: 120,
  },
  {
    title: '配置信息',
    dataIndex: 'config',
    width: 300,
    ellipsis: true,
  },
  {
    title: '默认',
    dataIndex: 'isDefault',
    width: 100,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 250,
    fixed: 'right',
  },
];

// 搜索表单配置
const searchFormConfig = {
  schemas: [
    {
      field: 'name',
      label: '名称',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'storageType',
      label: '存储类型',
      component: 'Select',
      componentProps: {
        options: [
          { label: '本地存储', value: 'local' },
          { label: 'MinIO', value: 'minio' },
          { label: 'S3', value: 's3' },
          { label: '阿里云OSS', value: 'oss' },
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
  return await storageEnvApi.list(params);
};

// 弹窗状态
const modalVisible = ref(false);
const currentEnvId = ref<number>();
const viewModalVisible = ref(false);
const viewConfig = ref<any>(null);

// 获取JSON预览（截断显示）
const getJsonPreview = (config: any) => {
  try {
    let jsonStr = '';
    if (typeof config === 'string') {
      jsonStr = config;
    } else {
      jsonStr = JSON.stringify(config);
    }
    
    // 显示前80个字符
    if (jsonStr.length > 80) {
      return jsonStr.substring(0, 80) + '...';
    }
    return jsonStr;
  } catch (e) {
    return String(config);
  }
};

// 查看完整配置
const handleViewConfig = (config: any) => {
  if (!config) return;
  viewConfig.value = config;
  viewModalVisible.value = true;
};

// 新增环境
const handleCreate = () => {
  currentEnvId.value = undefined;
  modalVisible.value = true;
};

// 编辑环境
const handleEdit = (record: any) => {
  currentEnvId.value = getEnvId(record);
  modalVisible.value = true;
};

// 设为默认
const handleSetDefault = async (envId: number) => {
  try {
    await storageEnvApi.setDefault(envId);
    message.success('设置成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('设置失败:', error);
  }
};

// 测试连接
const handleTest = async (envId: number) => {
  try {
    await storageEnvApi.testConnection(envId);
    message.success('连接测试成功');
  } catch (error) {
    console.error('连接测试失败:', error);
  }
};

// 删除环境
const handleDelete = async (envId: number) => {
  try {
    await storageEnvApi.delete(envId);
    message.success('删除成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 操作成功回调
const handleSuccess = () => {
  modalVisible.value = false;
  tableRef.value?.reload();
};
</script>

<style scoped lang="less">
.storage-env-page {
  padding: 0;

  .json-content {
    cursor: pointer;
    transition: all 0.2s;

    &.clickable:hover {
      color: #1890ff;
      background-color: #f0f5ff;
      padding: 4px 8px;
      border-radius: 4px;
    }

    .json-preview {
      font-family: 'Courier New', monospace;
      font-size: 12px;
      color: #666;
    }
  }

  .text-gray {
    color: #999;
  }
}
</style>
