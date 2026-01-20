<template>
  <BasicModal
    v-model:visible="visible"
    :title="isEdit ? '编辑菜单' : '新增菜单'"
    :width="700"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <BasicForm
      ref="formRef"
      :schemas="formSchemas"
      :model="formData"
      :label-width="100"
      :show-action-buttons="false"
    />
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { message } from 'ant-design-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import BasicForm from '@/components/Form/BasicForm.vue';
import { menuApi } from '@/api/menu';

const props = defineProps<{
  visible: boolean;
  menuId?: number;
  parentId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const isEdit = computed(() => !!props.menuId);
const loading = ref(false);
const formRef = ref();
const formData = ref<Record<string, any>>({});
const menuTreeOptions = ref<any[]>([]);

// 表单配置
const formSchemas = computed<any[]>(() => [
  {
    field: 'parentId',
    label: '上级菜单',
    component: 'TreeSelect',
    componentProps: {
      treeData: menuTreeOptions.value,
      showSearch: true,
      treeDefaultExpandAll: true,
      placeholder: '选择上级菜单',
      fieldNames: { label: 'label', value: 'value', children: 'children' },
    },
    defaultValue: props.parentId || 0,
    colProps: { span: 24 },
  },
  {
    field: 'menuType',
    label: '菜单类型',
    component: 'RadioGroup',
    defaultValue: 0,
    rules: [{ required: true, message: '请选择菜单类型' }],
    componentProps: {
      options: [
        { label: '目录', value: 0 },
        { label: '菜单', value: 1 },
        { label: '按钮', value: 2 },
      ],
    },
    colProps: { span: 24 },
  },
  {
    field: 'icon',
    label: '菜单图标',
    component: 'IconPicker',
    componentProps: {
      placeholder: '请选择图标',
    },
    show: (model) => model.menuType !== 2,
    colProps: { span: 24 },
    helpMessage: '点击选择图标',
  },
  {
    field: 'menuName',
    label: '菜单名称',
    component: 'Input',
    rules: [{ required: true, message: '请输入菜单名称' }],
    componentProps: {
      placeholder: '请输入菜单名称',
    },
    colProps: { span: 12 },
  },
  {
    field: 'sort',
    label: '显示排序',
    component: 'InputNumber',
    defaultValue: 0,
    componentProps: {
      min: 0,
      style: { width: '100%' },
    },
    colProps: { span: 12 },
  },
  {
    field: 'isFrame',
    label: '是否外链',
    component: 'RadioGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '是', value: 1 },
        { label: '否', value: 0 },
      ],
    },
    show: (model) => model.menuType !== 2,
    colProps: { span: 12 },
    helpMessage: '选择是外链则路由地址需要以http(s)://开头',
  },
  {
    field: 'path',
    label: '路由地址',
    component: 'Input',
    componentProps: {
      placeholder: '请输入路由地址',
    },
    show: (model) => model.menuType !== 2,
    colProps: { span: 12 },
    helpMessage: '访问的路由地址，如：user',
  },
  {
    field: 'component',
    label: '组件路径',
    component: 'Input',
    componentProps: {
      placeholder: '请输入组件路径',
    },
    show: (model) => model.menuType === 1,
    colProps: { span: 12 },
    helpMessage: '访问的组件路径，如：system/user/index',
  },
  {
    field: 'perms',
    label: '权限标识',
    component: 'Input',
    componentProps: {
      placeholder: '请输入权限标识',
      maxlength: 100,
    },
    show: (model) => model.menuType !== 0,
    colProps: { span: 12 },
    helpMessage: '控制器中定义的权限字符，如：user.list',
  },
  {
    field: 'query',
    label: '路由参数',
    component: 'Input',
    componentProps: {
      placeholder: '请输入路由参数',
      maxlength: 255,
    },
    show: (model) => model.menuType === 1,
    colProps: { span: 12 },
    helpMessage: '访问路由的默认传递参数，如：{"id": 1}',
  },
  {
    field: 'isCache',
    label: '是否缓存',
    component: 'RadioGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '缓存', value: 1 },
        { label: '不缓存', value: 0 },
      ],
    },
    show: (model) => model.menuType === 1,
    colProps: { span: 12 },
    helpMessage: '选择是则会被keep-alive缓存',
  },
  {
    field: 'visible',
    label: '显示状态',
    component: 'RadioGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '显示', value: 0 },
        { label: '隐藏', value: 1 },
      ],
    },
    show: (model) => model.menuType !== 2,
    colProps: { span: 12 },
    helpMessage: '选择隐藏则路由将不会出现在侧边栏',
  },
  {
    field: 'status',
    label: '菜单状态',
    component: 'RadioGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
    colProps: { span: 12 },
    helpMessage: '选择停用则路由将不会出现在侧边栏',
  },
]);

// 加载菜单树选项
const loadMenuTreeOptions = async () => {
  try {
    const data = await menuApi.getMenuTree();
    menuTreeOptions.value = [
      { label: '主类目', value: 0 },
      ...convertToTreeSelect(data),
    ];
  } catch (error) {
    console.error('加载菜单树失败:', error);
  }
};

// 转换为树形选择器格式
const convertToTreeSelect = (data: any[]): any[] => {
  return data.map(item => ({
    label: item.menuName,
    value: item.id,
    children: item.children && item.children.length > 0 ? convertToTreeSelect(item.children) : undefined,
  }));
};

// 加载菜单详情
const loadMenuDetail = async () => {
  if (!props.menuId) return;
  
  try {
    loading.value = true;
    const data = await menuApi.detail(props.menuId);
    
    console.log('加载的菜单数据:', data);
    
    // 等待下一个 tick 确保表单已渲染
    await new Promise(resolve => setTimeout(resolve, 100));
    
    // 设置表单值
    if (formRef.value) {
      formRef.value.setFieldsValue(data);
      console.log('设置表单值后:', formRef.value.getFieldsValue());
    }
  } catch (error) {
    console.error('加载菜单详情失败:', error);
  } finally {
    loading.value = false;
  }
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    loading.value = true;
    
    const values = formRef.value?.getFieldsValue();
    
    if (isEdit.value) {
      await menuApi.update(props.menuId!, values);
      message.success('更新成功');
    } else {
      await menuApi.create(values);
      message.success('创建成功');
    }
    
    emit('success');
  } catch (error) {
    console.error('提交失败:', error);
  } finally {
    loading.value = false;
  }
};

// 取消
const handleCancel = () => {
  formRef.value?.resetFields();
  formData.value = {};
};

// 监听弹窗显示
watch(
  () => props.visible,
  (val) => {
    if (val) {
      loadMenuTreeOptions();
      if (props.menuId) {
        loadMenuDetail();
      } else {
        formRef.value?.resetFields();
        formData.value = { parentId: props.parentId || 0 };
      }
    }
  }
);
</script>
