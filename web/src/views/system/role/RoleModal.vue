<template>
  <BasicModal
    v-model:visible="visible"
    :title="isEdit ? '编辑角色' : '新增角色'"
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
import { roleApi } from '@/api/role';
import type { FormSchema } from '@/types/components';

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

const isEdit = computed(() => !!props.roleId);
const loading = ref(false);
const formRef = ref();
const formData = ref<Record<string, any>>({});

// 表单配置
const formSchemas: FormSchema[] = [
  {
    field: 'roleName',
    label: '角色名称',
    component: 'Input',
    rules: [{ required: true, message: '请输入角色名称' }],
  },
  {
    field: 'roleKey',
    label: '角色标识',
    component: 'Input',
    rules: [{ required: true, message: '请输入角色标识' }],
    componentProps: {
      disabled: isEdit.value,
    },
  },
  {
    field: 'roleSort',
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
];

// 加载角色详情
const loadRoleDetail = async () => {
  if (!props.roleId) return;
  
  try {
    loading.value = true;
    const data = await roleApi.detail(props.roleId);
    formData.value = data;
    formRef.value?.setFieldsValue(data);
  } catch (error) {
    console.error('加载角色详情失败:', error);
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
      await roleApi.update({ ...values, roleId: props.roleId });
      message.success('更新成功');
    } else {
      await roleApi.create(values);
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
      if (props.roleId) {
        loadRoleDetail();
      } else {
        formRef.value?.resetFields();
        formData.value = {};
      }
    }
  }
);
</script>
