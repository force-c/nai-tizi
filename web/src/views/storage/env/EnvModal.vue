<template>
  <BasicModal
    v-model:visible="visible"
    :title="isEdit ? '编辑存储环境' : '新增存储环境'"
    :width="700"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <BasicForm
      ref="formRef"
      :schemas="formSchemas"
      :model="formData"
      :label-width="120"
      :show-action-buttons="false"
    />
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { message } from 'ant-design-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import BasicForm from '@/components/Form/BasicForm.vue';
import MonacoEditor from '@/components/MonacoEditor/index.vue';
import { storageEnvApi } from '@/api/storage';
import type { FormSchema } from '@/types/form';

const props = defineProps<{
  visible: boolean;
  envId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const isEdit = computed(() => !!props.envId);
const loading = ref(false);
const formRef = ref();
const formData = ref<Record<string, any>>({});
const editorTheme = ref<'vs-dark' | 'vs'>('vs-dark');

// JSON 验证规则
const validateJson = (_: any, value: string) => {
  if (!value) return Promise.resolve();
  try {
    JSON.parse(value);
    return Promise.resolve();
  } catch (e) {
    return Promise.reject('请输入有效的 JSON 格式');
  }
};

// 表单配置
const formSchemas = computed<FormSchema[]>(() => [
  {
    field: 'name',
    label: '环境名称',
    component: 'Input',
    rules: [{ required: true, message: '请输入环境名称' }],
  },
  {
    field: 'code',
    label: '环境编码',
    component: 'Input',
    rules: [{ required: true, message: '请输入环境编码' }],
  },
  {
    field: 'storageType',
    label: '存储类型',
    component: 'Select',
    defaultValue: 'local',
    rules: [{ required: true, message: '请选择存储类型' }],
    componentProps: {
      options: [
        { label: '本地存储', value: 'local' },
        { label: 'MinIO', value: 'minio' },
        { label: 'S3', value: 's3' },
        { label: '阿里云OSS', value: 'oss' },
      ],
    },
  },
  {
    field: 'isDefault',
    label: '默认环境',
    component: 'Select',
    defaultValue: false,
    componentProps: {
      options: [
        { label: '否', value: true },
        { label: '是', value: false },
      ],
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
    field: 'config',
    label: '配置信息',
    component: 'MonacoEditor',
    componentProps: {
      height: '300px',
      language: 'json',
      theme: editorTheme.value,
    },
    rules: [{ validator: validateJson }],
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

const toggleTheme = () => {
  editorTheme.value = editorTheme.value === 'vs-dark' ? 'vs' : 'vs-dark';
};

// 加载环境详情
const loadEnvDetail = async () => {
  if (!props.envId) return;
  
  try {
    loading.value = true;
    const data = await storageEnvApi.detail(props.envId);
    // 将config字段转换为字符串用于显示
    const formattedData = {
      ...data,
      config: data.config ? JSON.stringify(data.config, null, 2) : '',
    };
    formData.value = formattedData;
    formRef.value?.setFieldsValue(formattedData);
  } catch (error) {
    console.error('加载环境详情失败:', error);
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
    
    // 解析config字段为JSON对象
    let parsedConfig = null;
    if (values.config) {
      parsedConfig = JSON.parse(values.config);
    }
    
    if (isEdit.value) {
      await storageEnvApi.update({ 
        ...values, 
        config: parsedConfig,
        id: props.envId 
      });
      message.success('更新成功');
    } else {
      await storageEnvApi.create({
        ...values,
        config: parsedConfig,
      });
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
      if (props.envId) {
        loadEnvDetail();
      } else {
        formRef.value?.resetFields();
        formData.value = {};
      }
    }
  }
);
</script>
