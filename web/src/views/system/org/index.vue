<template>
  <div class="p-4">
    <!-- 搜索区域 -->
    <a-card class="mb-4" :bordered="false">
      <a-form ref="queryFormRef" :model="queryParams" layout="inline">
        <a-form-item label="组织名称" name="orgName">
          <a-input
            v-model:value="queryParams.orgName"
            placeholder="请输入组织名称"
            allow-clear
            @press-enter="handleQuery"
            style="width: 200px"
          />
        </a-form-item>
        <a-form-item label="组织编码" name="orgCode">
          <a-input
            v-model:value="queryParams.orgCode"
            placeholder="请输入组织编码"
            allow-clear
            @press-enter="handleQuery"
            style="width: 200px"
          />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-select
            v-model:value="queryParams.status"
            placeholder="组织状态"
            allow-clear
            style="width: 120px"
          >
            <a-select-option :value="0">正常</a-select-option>
            <a-select-option :value="1">停用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleQuery">
              <template #icon><SearchOutlined /></template>
              搜索
            </a-button>
            <a-button @click="resetQuery">
              <template #icon><ReloadOutlined /></template>
              重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 表格区域 -->
    <a-card :bordered="false">
      <template #title>
        <a-space>
          <a-button v-permission="'org.create'" type="primary" @click="handleAdd()">
            <template #icon><PlusOutlined /></template>
            新增
          </a-button>
          <a-button @click="handleToggleExpandAll">
            <template #icon><MenuUnfoldOutlined v-if="!isExpandAll" /><MenuFoldOutlined v-else /></template>
            {{ isExpandAll ? '折叠' : '展开' }}
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data-source="orgList"
        :loading="loading"
        :pagination="false"
        :row-key="(record) => record.id"
        :default-expand-all-rows="isExpandAll"
        :expanded-row-keys="expandedRowKeys"
        @expand="handleExpand"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 0 ? 'success' : 'error'">
              {{ record.status === 0 ? '正常' : '停用' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'createdTime'">
            {{ record.createdTime }}
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-tooltip title="修改">
                <a-button
                  v-permission="'org.update'"
                  type="link"
                  size="small"
                  @click="handleEdit(record)"
                >
                  <template #icon><EditOutlined /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip title="新增">
                <a-button
                  v-permission="'org.create'"
                  type="link"
                  size="small"
                  @click="handleAdd(record)"
                >
                  <template #icon><PlusOutlined /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip title="删除">
                <a-popconfirm
                  title="确定删除该组织吗？"
                  @confirm="handleDelete(record)"
                >
                  <a-button
                    v-permission="'org.delete'"
                    type="link"
                    danger
                    size="small"
                  >
                    <template #icon><DeleteOutlined /></template>
                  </a-button>
                </a-popconfirm>
              </a-tooltip>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 编辑弹窗 -->
    <OrgModal
      v-model:visible="modalVisible"
      :org-id="currentOrgId"
      :parent-id="parentId"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import {
  SearchOutlined,
  ReloadOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
} from '@ant-design/icons-vue';
import OrgModal from './OrgModal.vue';
import { organizationApi } from '@/api/organization';
import { usePermission } from '@/composables/usePermission';

const { hasPermission } = usePermission();

const loading = ref(false);
const orgList = ref<any[]>([]);
const isExpandAll = ref(true);
const expandedRowKeys = ref<number[]>([]);

// 查询参数
const queryParams = reactive({
  orgName: '',
  orgCode: '',
  status: undefined as number | undefined,
});

const queryFormRef = ref();

// 表格列定义
const columns = [
  {
    title: '组织名称',
    dataIndex: 'orgName',
    key: 'orgName',
    width: 260,
  },
  {
    title: '组织编码',
    dataIndex: 'orgCode',
    key: 'orgCode',
    width: 200,
    align: 'center',
  },
  {
    title: '显示顺序',
    dataIndex: 'sort',
    key: 'sort',
    width: 120,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    align: 'center',
  },
  {
    title: '创建时间',
    key: 'createdTime',
    width: 180,
    align: 'center',
  },
  {
    title: '操作',
    key: 'action',
    fixed: 'right',
    width: 180,
    align: 'center',
  },
];

// 加载组织树
const loadOrgTree = async () => {
  try {
    loading.value = true;
    const data = await organizationApi.tree();
    orgList.value = data;
    
    // 默认展开所有节点
    if (isExpandAll.value) {
      expandedRowKeys.value = getAllKeys(data);
    }
  } catch (error) {
    console.error('加载组织树失败:', error);
    message.error('加载组织树失败');
  } finally {
    loading.value = false;
  }
};

// 获取所有节点的key
const getAllKeys = (data: any[]): number[] => {
  const keys: number[] = [];
  const traverse = (nodes: any[]) => {
    nodes.forEach((node) => {
      keys.push(node.id);
      if (node.children && node.children.length > 0) {
        traverse(node.children);
      }
    });
  };
  traverse(data);
  return keys;
};

// 搜索按钮操作
const handleQuery = () => {
  loadOrgTree();
};

// 重置按钮操作
const resetQuery = () => {
  queryFormRef.value?.resetFields();
  queryParams.orgName = '';
  queryParams.orgCode = '';
  queryParams.status = undefined;
  handleQuery();
};

// 展开/折叠操作
const handleToggleExpandAll = () => {
  isExpandAll.value = !isExpandAll.value;
  if (isExpandAll.value) {
    expandedRowKeys.value = getAllKeys(orgList.value);
  } else {
    expandedRowKeys.value = [];
  }
};

// 处理展开事件
const handleExpand = (expanded: boolean, record: any) => {
  if (expanded) {
    if (!expandedRowKeys.value.includes(record.id)) {
      expandedRowKeys.value = [...expandedRowKeys.value, record.id];
    }
  } else {
    expandedRowKeys.value = expandedRowKeys.value.filter(key => key !== record.id);
  }
};

// 弹窗状态
const modalVisible = ref(false);
const currentOrgId = ref<number>();
const parentId = ref<number>();

// 新增按钮操作
const handleAdd = (record?: any) => {
  currentOrgId.value = undefined;
  parentId.value = record?.id;
  modalVisible.value = true;
};

// 编辑按钮操作
const handleEdit = (record: any) => {
  currentOrgId.value = record.id;
  parentId.value = undefined;
  modalVisible.value = true;
};

// 删除按钮操作
const handleDelete = async (record: any) => {
  try {
    await organizationApi.delete(record.id);
    message.success('删除成功');
    loadOrgTree();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 操作成功回调
const handleSuccess = () => {
  modalVisible.value = false;
  loadOrgTree();
};

onMounted(() => {
  loadOrgTree();
});
</script>

<style scoped lang="less">
:deep(.ant-table) {
  .ant-table-tbody > tr > td {
    padding: 12px 8px;
  }
}
</style>
