<template>
  <div class="basic-table">
    <!-- 搜索表单 -->
    <div v-if="useSearchForm && formConfig" class="table-search">
      <BasicForm
        ref="searchFormRef"
        v-bind="formConfig"
        :show-action-buttons="true"
        submit-button-text="搜索"
        reset-button-text="重置"
        @submit="handleSearch"
        @reset="handleReset"
      />
    </div>

    <!-- 工具栏 -->
    <div v-if="$slots.toolbar" class="table-toolbar">
      <slot name="toolbar"></slot>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="paginationConfig"
      :row-key="rowKey"
      :row-selection="rowSelection"
      v-bind="$attrs"
      @change="handleTableChange"
    >
      <!-- 自定义列插槽 -->
      <template #bodyCell="{ column, record, index }">
        <slot
          name="bodyCell"
          :column="column"
          :record="record"
          :index="index"
        ></slot>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue';
import BasicForm from '../Form/BasicForm.vue';

interface Props {
  columns?: any[];
  dataSource?: any[];
  loading?: boolean;
  pagination?: boolean | Record<string, any>;
  api?: (params: any) => Promise<any>;
  useSearchForm?: boolean;
  formConfig?: any;
  showActionColumn?: boolean;
  immediate?: boolean;
  rowKey?: string;
}

const props = withDefaults(defineProps<Props>(), {
  columns: () => [],
  dataSource: () => [],
  loading: false,
  pagination: true,
  useSearchForm: false,
  showActionColumn: false,
  immediate: true,
  rowKey: 'id',
});

const emit = defineEmits<{
  (e: 'register', methods: any): void;
}>();

// 表格数据
const dataSource = ref<any[]>(props.dataSource || []);
const loading = ref(false);
const selectedRowKeys = ref<any[]>([]);

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条`,
});

// 搜索表单引用
const searchFormRef = ref();

// 搜索参数
const searchParams = ref<Record<string, any>>({});

// 计算分页配置
const paginationConfig = computed(() => {
  if (props.pagination === false) {
    return false;
  }
  if (typeof props.pagination === 'object') {
    return { ...pagination, ...props.pagination };
  }
  return pagination;
});

// 行选择配置
const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys: any[]) => {
    selectedRowKeys.value = keys;
  },
}));

// 加载数据
const loadData = async (params?: any) => {
  if (!props.api) {
    dataSource.value = props.dataSource || [];
    return;
  }

  loading.value = true;
  try {
    const requestParams = {
      pageNum: pagination.current,
      pageSize: pagination.pageSize,
      ...searchParams.value,
      ...params,
    };

    const res = await props.api(requestParams);
    
    // 标准格式：{ data: { records: [], total: 0, pages: 0 } }
    if (res.data && res.data.records) {
      dataSource.value = res.data.records;
      pagination.total = res.data.total || 0;
    } else if (res.records) {
      // 直接返回格式：{ records: [], total: 0, pages: 0 }
      dataSource.value = res.records;
      pagination.total = res.total || 0;
    } else {
      console.error('Invalid pagination response format. Expected { records: [], total: 0 }');
    }
  } catch (error) {
    console.error('加载表格数据失败:', error);
  } finally {
    loading.value = false;
  }
};

// 处理表格变化（分页、排序、筛选）
const handleTableChange = (pag: any, filters: any, sorter: any) => {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  loadData();
};

// 处理搜索
const handleSearch = (values: Record<string, any>) => {
  searchParams.value = values;
  pagination.current = 1; // 重置到第一页
  loadData();
};

// 处理重置
const handleReset = () => {
  searchParams.value = {};
  pagination.current = 1;
  loadData();
};

// 刷新表格
const reload = () => {
  loadData();
};

// 重置表格（清空搜索条件）
const reset = () => {
  searchParams.value = {};
  pagination.current = 1;
  if (searchFormRef.value) {
    searchFormRef.value.resetFields();
  }
  loadData();
};

// 获取选中行
const getSelectRows = () => {
  return dataSource.value.filter((item) => 
    selectedRowKeys.value.includes(item[props.rowKey])
  );
};

// 暴露方法给父组件
defineExpose({
  reload,
  reset,
  getSelectRows,
  loadData,
});

// 注册方法
emit('register', {
  reload,
  reset,
  getSelectRows,
  loadData,
});

// 监听dataSource变化
watch(
  () => props.dataSource,
  (newVal) => {
    if (newVal && !props.api) {
      dataSource.value = newVal;
    }
  },
  { deep: true }
);

// 组件挂载时加载数据
onMounted(() => {
  if (props.immediate && props.api) {
    loadData();
  }
});
</script>

<style scoped lang="less">
.basic-table {
  .table-toolbar {
    margin-bottom: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .table-search {
    margin-bottom: 16px;
    padding: 16px;
    background-color: #fafafa;
    border-radius: 4px;
  }
}
</style>
