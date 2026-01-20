<template>
  <a-form
    ref="formRef"
    :model="formModel"
    :label-width="labelWidth"
    :label-align="labelAlign"
    :layout="layout"
    v-bind="$attrs"
  >
    <a-row :gutter="16">
      <a-col
        v-for="schema in visibleSchemas"
        :key="schema.field"
        v-bind="schema.colProps || { span: 24 }"
      >
        <a-form-item
          :label="schema.label"
          :name="schema.field"
          :rules="schema.rules"
        >
          <template v-if="schema.helpMessage" #extra>
            <span class="form-item-help">{{ schema.helpMessage }}</span>
          </template>
          <!-- Switch 组件使用 v-model:checked -->
          <component
            v-if="schema.component === 'Switch'"
            :is="getComponent(schema.component)"
            v-model:checked="formModel[schema.field]"
            v-bind="schema.componentProps"
          />
          <!-- IconPicker 组件使用 v-model -->
          <component
            v-else-if="schema.component === 'IconPicker'"
            :is="getComponent(schema.component)"
            v-model="formModel[schema.field]"
            v-bind="schema.componentProps"
          />
          <!-- MonacoEditor 组件使用 v-model，支持插槽 -->
          <slot
            v-else-if="schema.component === 'MonacoEditor'"
            :name="schema.field"
            :schema="schema"
            :model="formModel"
          >
            <component
              :is="getComponent(schema.component)"
              v-model="formModel[schema.field]"
              v-bind="schema.componentProps"
            />
          </slot>
          <!-- 其他组件使用 v-model:value -->
          <component
            v-else
            :is="getComponent(schema.component)"
            v-model:value="formModel[schema.field]"
            v-bind="schema.componentProps"
            :placeholder="getPlaceholder(schema)"
          />
        </a-form-item>
      </a-col>
      
      <!-- 操作按钮放在同一行最右侧 -->
      <a-col v-if="showActionButtons" :flex="1" style="text-align: right;">
        <a-form-item class="form-actions">
          <a-space>
            <a-button type="primary" @click="handleSubmit">
              {{ submitButtonText }}
            </a-button>
            <a-button @click="handleReset">
              {{ resetButtonText }}
            </a-button>
          </a-space>
        </a-form-item>
      </a-col>
    </a-row>
  </a-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import type { FormSchema } from '@/types/form';
import { 
  Input, 
  InputNumber, 
  Select, 
  DatePicker, 
  Switch, 
  Textarea,
  TreeSelect,
  Radio,
} from 'ant-design-vue';
import IconPicker from '@/components/IconPicker/index.vue';
import MonacoEditor from '@/components/MonacoEditor/index.vue';

interface Props {
  schemas?: FormSchema[];
  model?: Record<string, any>;
  labelWidth?: number;
  labelAlign?: 'left' | 'right';
  showActionButtons?: boolean;
  submitButtonText?: string;
  resetButtonText?: string;
  layout?: 'horizontal' | 'vertical' | 'inline';
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  schemas: () => [],
  model: () => ({}),
  labelWidth: 100,
  labelAlign: 'right',
  showActionButtons: true,
  submitButtonText: '提交',
  resetButtonText: '重置',
  layout: 'horizontal',
  disabled: false,
});

const emit = defineEmits<{
  (e: 'submit', values: Record<string, any>): void;
  (e: 'reset'): void;
}>();

const formRef = ref();
const formModel = reactive<Record<string, any>>({});

// 组件映射
const componentMap: Record<string, any> = {
  Input: Input,
  InputNumber: InputNumber,
  Select: Select,
  DatePicker: DatePicker,
  Switch: Switch,
  Textarea: Input.TextArea,
  TreeSelect: TreeSelect,
  RadioGroup: Radio.Group,
  IconPicker: IconPicker,
  MonacoEditor: MonacoEditor,
};

// 获取组件
const getComponent = (componentName: string) => {
  return componentMap[componentName] || Input;
};

// 获取placeholder
const getPlaceholder = (schema: FormSchema) => {
  const { component, label } = schema;
  if (component === 'Select') {
    return `请选择${label}`;
  }
  return `请输入${label}`;
};

// 计算可见的schema
const visibleSchemas = computed(() => {
  return props.schemas.filter((schema) => {
    if (typeof schema.show === 'function') {
      return schema.show(formModel);
    }
    return schema.show !== false;
  });
});

// 初始化表单值
const initFormModel = () => {
  props.schemas.forEach((schema) => {
    if (schema.defaultValue !== undefined) {
      formModel[schema.field] = schema.defaultValue;
    }
  });

  // 合并传入的model
  Object.assign(formModel, props.model);
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate();
    emit('submit', { ...formModel });
  } catch (error) {
    console.error('表单验证失败:', error);
  }
};

// 重置表单
const handleReset = () => {
  formRef.value.resetFields();
  emit('reset');
};

// 设置表单值
const setFieldsValue = (values: Record<string, any>) => {
  Object.assign(formModel, values);
};

// 获取表单值
const getFieldsValue = () => {
  return { ...formModel };
};

// 重置字段
const resetFields = () => {
  formRef.value.resetFields();
};

// 验证表单
const validate = async () => {
  return await formRef.value.validate();
};

// 暴露方法
defineExpose({
  setFieldsValue,
  getFieldsValue,
  resetFields,
  validate,
});

// 监听model变化
watch(
  () => props.model,
  (newVal) => {
    if (newVal) {
      Object.assign(formModel, newVal);
    }
  },
  { deep: true, immediate: true }
);

// 初始化
initFormModel();
</script>

<style scoped lang="less">
.form-actions {
  margin-bottom: 0;
  
  :deep(.ant-form-item-control-input) {
    min-height: auto;
  }
}
</style>
