<template>
  <div class="config-container">
    <a-card :bordered="false">
      <!-- 搜索表单 -->
      <a-form layout="inline" :model="searchForm" class="search-form">
      <a-form-item label="名称">
        <a-input
          v-model:value="searchForm.name"
          placeholder="请输入配置名称"
          allow-clear
          style="width: 200px"
        />
      </a-form-item>
      <a-form-item label="编码">
        <a-input
          v-model:value="searchForm.code"
          placeholder="请输入配置编码"
          allow-clear
          style="width: 200px"
        />
      </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSearch">
              <template #icon><SearchOutlined /></template>
              查询
            </a-button>
            <a-button @click="handleReset">
              <template #icon><ReloadOutlined /></template>
              重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>

      <!-- 操作按钮 -->
      <div class="table-operations">
        <a-space>
          <a-button type="primary" @click="handleCreate">
            <template #icon><PlusOutlined /></template>
            新增
          </a-button>
          <a-button danger :disabled="!hasSelected" @click="handleBatchDelete">
            <template #icon><DeleteOutlined /></template>
            批量删除
          </a-button>
        </a-space>
      </div>

      <!-- 数据表格 -->
      <a-table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="pagination"
        :row-selection="rowSelection"
        :row-key="(record) => record.id"
        @change="handleTableChange"
      >
        <!-- 编码列 -->
        <template #code="{ text }">
          <a-tag color="blue">{{ text }}</a-tag>
        </template>

        <!-- 内容列 - JSON格式展示 -->
        <template #data="{ record }">
          <div 
            class="json-content clickable" 
            @click="handleViewData(record.data)"
          >
            <span v-if="record.data" class="json-preview">
              {{ getJsonPreview(record.data) }}
            </span>
            <span v-else class="text-gray">-</span>
          </div>
        </template>

        <!-- 创建时间列 -->
        <template #createdTime="{ text }">
          {{ text || '-' }}
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <a-space>
            <a-button type="link" size="small" @click="handleEdit(record)">
              编辑
            </a-button>
            <a-popconfirm
              title="确定删除该配置吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" danger size="small">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <ConfigModal
      v-model:visible="modalVisible"
      :config-id="currentConfigId"
      @success="handleSuccess"
    />

    <!-- 查看数据弹窗 -->
    <ConfigViewModal
      v-model:visible="viewModalVisible"
      :data="viewData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import {
  SearchOutlined,
  ReloadOutlined,
  PlusOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue';
import { configApi, type Config } from '@/api/config';
import ConfigModal from './ConfigModal.vue';
import ConfigViewModal from './ConfigViewModal.vue';

// 搜索表单
const searchForm = reactive({
  name: '',
  code: '',
});

// 表格数据
const dataSource = ref<Config[]>([]);
const loading = ref(false);
const selectedRowKeys = ref<number[]>([]);

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条`,
});

// 表格列配置
const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 150,
  },
  {
    title: '编码',
    dataIndex: 'code',
    key: 'code',
    width: 120,
    slots: { customRender: 'code' },
  },
  {
    title: '内容',
    dataIndex: 'data',
    key: 'data',
    width: 300,
    ellipsis: true,
    slots: { customRender: 'data' },
  },
  {
    title: '备注',
    dataIndex: 'remark',
    key: 'remark',
    width: 200,
    ellipsis: true,
  },
  {
    title: '创建时间',
    dataIndex: 'createdTime',
    key: 'createdTime',
    width: 180,
    slots: { customRender: 'createdTime' },
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// 行选择配置
const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys: number[]) => {
    selectedRowKeys.value = keys;
  },
}));

// 是否有选中项
const hasSelected = computed(() => selectedRowKeys.value.length > 0);

// 弹窗相关
const modalVisible = ref(false);
const currentConfigId = ref<number>();
const viewModalVisible = ref(false);
const viewData = ref<any>(null);

// 获取JSON预览（截断显示）
const getJsonPreview = (data: any) => {
  try {
    let jsonStr = '';
    if (typeof data === 'string') {
      jsonStr = data;
    } else {
      jsonStr = JSON.stringify(data);
    }
    
    // 显示前80个字符
    if (jsonStr.length > 80) {
      return jsonStr.substring(0, 80) + '...';
    }
    return jsonStr;
  } catch (e) {
    return String(data);
  }
};

// 查看完整数据
const handleViewData = (data: any) => {
  if (!data) return;
  viewData.value = data;
  viewModalVisible.value = true;
};

// 加载数据
const loadData = async () => {
  try {
    loading.value = true;
    const res = await configApi.page({
      pageNum: pagination.current,
      pageSize: pagination.pageSize,
      name: searchForm.name || undefined,
      code: searchForm.code || undefined,
    });

    dataSource.value = res.records || [];
    pagination.total = res.total || 0;
  } catch (error) {
    console.error('加载配置列表失败:', error);
    message.error('加载配置列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  pagination.current = 1;
  loadData();
};

// 重置
const handleReset = () => {
  searchForm.name = '';
  searchForm.code = '';
  pagination.current = 1;
  loadData();
};

// 表格变化
const handleTableChange = (pag: any) => {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  loadData();
};

// 新增
const handleCreate = () => {
  currentConfigId.value = undefined;
  modalVisible.value = true;
};

// 编辑
const handleEdit = (record: Config) => {
  currentConfigId.value = record.id;
  modalVisible.value = true;
};

// 删除
const handleDelete = async (id: number) => {
  try {
    await configApi.delete(id);
    message.success('删除成功');
    loadData();
  } catch (error) {
    console.error('删除配置失败:', error);
    message.error('删除配置失败');
  }
};

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的配置');
    return;
  }

  try {
    await configApi.batchDelete(selectedRowKeys.value);
    message.success('批量删除成功');
    selectedRowKeys.value = [];
    loadData();
  } catch (error) {
    console.error('批量删除配置失败:', error);
    message.error('批量删除配置失败');
  }
};

// 操作成功回调
const handleSuccess = () => {
  modalVisible.value = false;
  loadData();
};

// 初始化
onMounted(() => {
  loadData();
});
</script>

<style scoped lang="less">
.config-container {
  padding: 16px;

  .search-form {
    margin-bottom: 16px;
  }

  .table-operations {
    margin-bottom: 16px;
  }

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
