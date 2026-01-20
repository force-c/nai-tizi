<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    :width="width"
    :confirm-loading="confirmLoading"
    :mask-closable="maskClosable"
    :destroy-on-close="destroyOnClose"
    :body-style="bodyStyle"
    :wrap-class-name="wrapClassName"
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <slot></slot>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

interface Props {
  visible: boolean;
  title?: string;
  width?: number;
  showOkButton?: boolean;
  showCancelButton?: boolean;
  okText?: string;
  cancelText?: string;
  confirmLoading?: boolean;
  maskClosable?: boolean;
  destroyOnClose?: boolean;
  bodyStyle?: Record<string, any>;
  wrapClassName?: string;
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  title: '',
  width: 520,
  showOkButton: true,
  showCancelButton: true,
  okText: '确定',
  cancelText: '取消',
  confirmLoading: false,
  maskClosable: true,
  destroyOnClose: true,
  bodyStyle: undefined,
  wrapClassName: undefined,
});

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'ok'): void;
  (e: 'cancel'): void;
}>();

const visible = ref(props.visible);

const handleOk = () => {
  emit('ok');
};

const handleCancel = () => {
  visible.value = false;
  emit('update:visible', false);
  emit('cancel');
};

watch(
  () => props.visible,
  (val) => {
    visible.value = val;
  }
);

watch(visible, (val) => {
  emit('update:visible', val);
});

defineExpose({
  close: handleCancel,
});
</script>
