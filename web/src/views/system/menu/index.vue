<template>
  <div class="menu-page">
    <div class="toolbar">
      <a-space>
        <a-button v-permission="'menu.create'" type="primary" @click="handleCreate">
          <template #icon><PlusOutlined /></template>
            新增
        </a-button>
        <a-button @click="expandAll">
          <template #icon><DownOutlined /></template>
          展开全部
        </a-button>
        <a-button @click="collapseAll">
          <template #icon><UpOutlined /></template>
          折叠全部
        </a-button>
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="menuTree"
      :loading="loading"
      :pagination="false"
      :expanded-row-keys="expandedKeys"
      row-key="id"
      @expand="handleExpand"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'icon'">
          <component :is="getIcon(record.icon)" v-if="record.icon" />
          <span v-else class="text-gray-400">-</span>
        </template>
        <template v-else-if="column.dataIndex === 'menuType'">
          <a-tag v-if="record.menuType === 0" color="blue">目录</a-tag>
          <a-tag v-else-if="record.menuType === 1" color="green">菜单</a-tag>
          <a-tag v-else-if="record.menuType === 2" color="orange">按钮</a-tag>
          <a-tag v-else color="default">{{ record.menuType }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'isFrame'">
          <a-tag :color="record.isFrame === 1 ? 'success' : 'default'">
            {{ record.isFrame === 1 ? '是' : '否' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'isCache'">
          <a-tag :color="record.isCache === 1 ? 'success' : 'default'">
            {{ record.isCache === 1 ? '是' : '否' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'visible'">
          <a-tag :color="record.visible === 0 ? 'success' : 'error'">
            {{ record.visible === 0 ? '显示' : '隐藏' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <a-tag :color="record.status === 0 ? 'success' : 'error'">
            {{ record.status === 0 ? '正常' : '停用' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-tooltip title="查看">
              <a-button
                type="link"
                size="small"
                @click="handleView(record)"
              >
                <template #icon><EyeOutlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="修改">
              <a-button
                v-permission="'menu.update'"
                type="link"
                size="small"
                @click="handleEdit(record)"
              >
                <template #icon><EditOutlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="新增" v-if="record.menuType !== 2">
              <a-button
                v-permission="'menu.create'"
                type="link"
                size="small"
                @click="handleCreateChild(record)"
              >
                <template #icon><PlusOutlined /></template>
              </a-button>
            </a-tooltip>
            <a-tooltip title="删除">
              <a-popconfirm
                title="确定删除该菜单吗？"
                @confirm="handleDelete(record.id)"
              >
                <a-button
                  v-permission="'menu.delete'"
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

    <MenuModal
      v-model:visible="modalVisible"
      :menu-id="currentMenuId"
      :parent-id="parentId"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { 
  PlusOutlined, 
  DownOutlined, 
  UpOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
} from '@ant-design/icons-vue';
import MenuModal from './MenuModal.vue';
import { menuApi } from '@/api/menu';
import { usePermission } from '@/composables/usePermission';
import { getIcon } from '@/config/icons';

const { hasPermission } = usePermission();

// 表格列配置
const columns = [
  {
    title: '菜单名称',
    dataIndex: 'menuName',
    width: 200,
    ellipsis: true,
  },
  {
    title: '图标',
    dataIndex: 'icon',
    width: 80,
  },
  {
    title: '排序',
    dataIndex: 'sort',
    width: 80,
  },
  {
    title: '组件路径',
    dataIndex: 'component',
    width: 220,
    ellipsis: true,
  },
  {
    title: '权限标识',
    dataIndex: 'perms',
    width: 180,
    ellipsis: true,
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
    key: 'action',
    width: 200,
    fixed: 'right',
  },
];

const loading = ref(false);
const menuTree = ref<any[]>([]);
const expandedKeys = ref<number[]>([]);

// 加载菜单树
const loadMenuTree = async () => {
  try {
    loading.value = true;
    const data = await menuApi.getMenuTree();
    menuTree.value = data;
    
    // 默认不展开，或者只展开第一层（根据需求选择）
    // expandedKeys.value = data.map((item: any) => item.id);
    expandedKeys.value = [];
  } catch (error) {
    console.error('加载菜单树失败:', error);
  } finally {
    loading.value = false;
  }
};

// 处理单个节点展开/收缩
const handleExpand = (expanded: boolean, record: any) => {
  if (expanded) {
    // 展开：添加到 expandedKeys
    if (!expandedKeys.value.includes(record.id)) {
      expandedKeys.value = [...expandedKeys.value, record.id];
    }
  } else {
    // 收缩：从 expandedKeys 中移除
    expandedKeys.value = expandedKeys.value.filter(key => key !== record.id);
  }
};

// 展开全部
const expandAll = () => {
  const getAllKeys = (data: any[]): number[] => {
    let keys: number[] = [];
    data.forEach(item => {
      keys.push(item.id);
      if (item.children && item.children.length > 0) {
        keys = keys.concat(getAllKeys(item.children));
      }
    });
    return keys;
  };
  expandedKeys.value = getAllKeys(menuTree.value);
};

// 折叠全部
const collapseAll = () => {
  expandedKeys.value = [];
};

// 弹窗状态
const modalVisible = ref(false);
const currentMenuId = ref<number>();
const parentId = ref<number>();

// 新增菜单
const handleCreate = () => {
  currentMenuId.value = undefined;
  parentId.value = undefined;
  modalVisible.value = true;
};

// 新增子菜单
const handleCreateChild = (record: any) => {
  currentMenuId.value = undefined;
  parentId.value = record.id;
  modalVisible.value = true;
};

// 查看菜单（只读模式）
const handleView = (record: any) => {
  currentMenuId.value = record.id;
  parentId.value = undefined;
  modalVisible.value = true;
};

// 编辑菜单
const handleEdit = (record: any) => {
  currentMenuId.value = record.id;
  parentId.value = undefined;
  modalVisible.value = true;
};

// 删除菜单
const handleDelete = async (menuId: number) => {
  try {
    await menuApi.delete(menuId);
    message.success('删除成功');
    loadMenuTree();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 操作成功回调
const handleSuccess = () => {
  modalVisible.value = false;
  loadMenuTree();
};

onMounted(() => {
  loadMenuTree();
});
</script>

<style scoped lang="less">
.menu-page {
  background: #fff;
  padding: 24px;
  border-radius: 4px;

  .toolbar {
    margin-bottom: 16px;
  }
}
</style>
