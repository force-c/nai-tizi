<template>
  <BasicModal
    v-model:visible="visible"
    :title="isEdit ? '编辑组织' : '新增组织'"
    :width="600"
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
import { organizationApi } from '@/api/organization';
import type { FormSchema } from '@/types/components';

const props = defineProps<{
  visible: boolean;
  orgId?: number;
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

const isEdit = computed(() => !!props.orgId);
const loading = ref(false);
const formRef = ref();
const formData = ref<Record<string, any>>({});
const orgTreeOptions = ref<any[]>([]);

// 表单配置
const formSchemas = computed<FormSchema[]>(() => [
  {
    field: 'parentId',
    label: '上级组织',
    component: 'Select',
    componentProps: {
      options: orgTreeOptions.value,
      showSearch: true,
      placeholder: '选择上级组织',
    },
    defaultValue: props.parentId || 0,
  },
  {
    field: 'orgName',
    label: '组织名称',
    component: 'Input',
    rules: [{ required: true, message: '请输入组织名称' }],
  },
  {
    field: 'orgCode',
    label: '组织编码',
    component: 'Input',
    rules: [{ required: true, message: '请输入组织编码' }],
  },
  {
    field: 'leader',
    label: '负责人',
    component: 'Input',
  },
  {
    field: 'phone',
    label: '联系电话',
    component: 'Input',
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'Input',
    rules: [{ type: 'email', message: '请输入正确的邮箱格式' }],
  },
  {
    field: 'sort',
    label: '显示顺序',
    component: 'InputNumber',
    defaultValue: 0,
    componentProps: {
      min: 0,
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
  },
  {
    field: 'remark',
    label: '备注',
    component: 'Textarea',
    componentProps: {
      rows: 3,
    },
  },
]);

// 加载组织树选项
const loadOrgTreeOptions = async () => {
  try {
    const data = await organizationApi.tree();
    orgTreeOptions.value = [
      { label: '顶级组织', value: 0 },
      ...convertToTreeSelect(data),
    ];
  } catch (error) {
    console.error('加载组织树失败:', error);
  }
};

// 转换为树形选择器格式
const convertToTreeSelect = (data: any[]): any[] => {
  return data.map(item => ({
    label: item.orgName,
    value: item.id,
    children: item.children && item.children.length > 0 ? convertToTreeSelect(item.children) : undefined,
  }));
};

// 加载组织详情
const loadOrgDetail = async () => {
  if (!props.orgId) return;
  
  try {
    loading.value = true;
    const data = await organizationApi.detail(props.orgId);
    formData.value = data;
    formRef.value?.setFieldsValue(data);
  } catch (error) {
    console.error('加载组织详情失败:', error);
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
      await organizationApi.update(props.orgId!, values);
      message.success('更新成功');
    } else {
      await organizationApi.create(values);
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
      loadOrgTreeOptions();
      if (props.orgId) {
        loadOrgDetail();
      } else {
        formRef.value?.resetFields();
        formData.value = { parentId: props.parentId || 0 };
      }
    }
  }
);
</script>
