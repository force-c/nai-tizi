<template>
  <BasicModal
    v-model:visible="visible"
    title=""
    :width="700"
    :confirm-loading="loading"
    wrap-class-name="config-modal"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <BasicForm
      ref="formRef"
      :schemas="formSchemas"
      :model="formData"
      :label-width="100"
      :show-action-buttons="false"
    >
      <template #data="{ schema, model }">
        <MonacoEditor
          v-model="model.data"
          :height="schema.componentProps?.height || '300px'"
          :language="schema.componentProps?.language || 'json'"
          :theme="editorTheme"
        >
          <template #toolbar>
            <a-button size="small" type="text" @click="toggleTheme" :title="'切换主题'">
              {{ editorTheme === 'vs-dark' ? '🌙' : '☀️' }}
            </a-button>
          </template>
        </MonacoEditor>
      </template>
    </BasicForm>
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { message } from 'ant-design-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import BasicForm from '@/components/Form/BasicForm.vue';
import MonacoEditor from '@/components/MonacoEditor/index.vue';
import { configApi } from '@/api/config';
import type { FormSchema } from '@/types/form';
import { useAuthStore } from '@/stores/auth';

const props = defineProps<{ visible: boolean; configId?: number; }>();
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

const authStore = useAuthStore();
const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const isEdit = computed(() => !!props.configId);
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

const formSchemas = computed<FormSchema[]>(() => [
  { field: 'name', label: '配置名称', component: 'Input', rules: [{ required: true, message: '请输入配置名称' }] },
  { field: 'code', label: '配置编码', component: 'Input', rules: [{ required: true, message: '请输入配置编码' }] },
  { 
    field: 'data', 
    label: '配置数据', 
    component: 'MonacoEditor', 
    componentProps: { height: '300px', language: 'json', theme: editorTheme.value },
    rules: [{ validator: validateJson }]
  },
  { field: 'remark', label: '备注', component: 'Textarea', componentProps: { rows: 3 } },
]);

const toggleTheme = () => {
  editorTheme.value = editorTheme.value === 'vs-dark' ? 'vs' : 'vs-dark';
};

const loadConfigDetail = async () => {
  if (!props.configId) return;
  try {
    loading.value = true;
    const data = await configApi.detail(props.configId);
    // 将data字段转换为字符串用于显示
    const formattedData = {
      ...data,
      data: data.data ? JSON.stringify(data.data, null, 2) : '',
    };
    formData.value = formattedData;
    formRef.value?.setFieldsValue(formattedData);
  } catch (error) {
    console.error('加载配置详情失败:', error);
    message.error('加载配置详情失败');
  } finally {
    loading.value = false;
  }
};

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    
    const userId = authStore.userInfo?.userId;
    if (!userId) {
      message.error('用户信息获取失败，请重新登录');
      return;
    }

    loading.value = true;
    const values = formRef.value?.getFieldsValue();
    
    // 解析data字段为JSON对象
    let parsedData = null;
    if (values.data) {
      parsedData = JSON.parse(values.data);
    }

    if (isEdit.value) {
      await configApi.update({
        id: props.configId!,
        name: values.name,
        code: values.code,
        data: parsedData,
        remark: values.remark,
        updateBy: userId,
      });
      message.success('更新成功');
    } else {
      await configApi.create({
        name: values.name,
        code: values.code,
        data: parsedData,
        remark: values.remark,
        createBy: userId,
        updateBy: userId,
      });
      message.success('创建成功');
    }
    emit('success');
    visible.value = false;
  } catch (error) {
    console.error('提交失败:', error);
    message.error('提交失败');
  } finally {
    loading.value = false;
  }
};

const handleCancel = () => {
  formRef.value?.resetFields();
  formData.value = {};
};

watch(() => props.visible, (val) => {
  if (val) {
    if (props.configId) {
      loadConfigDetail();
    } else {
      formRef.value?.resetFields();
      formData.value = {};
    }
  }
});
</script>

<style lang="less">
.config-modal {
  .ant-modal-header {
    display: none;
  }
}
</style>
