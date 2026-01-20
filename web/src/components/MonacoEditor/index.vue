<template>
  <div class="monaco-editor-wrapper" :data-theme="theme">
    <div v-if="$slots.toolbar" class="monaco-editor-toolbar">
      <slot name="toolbar"></slot>
    </div>
    <vue-monaco-editor
      v-model:value="content"
      :language="language"
      :theme="theme"
      :height="height"
      :options="editorOptions"
      @mount="handleMount"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { VueMonacoEditor } from '@guolao/vue-monaco-editor';

interface Props {
  modelValue?: string;
  language?: string;
  height?: string | number;
  theme?: 'vs-dark' | 'vs' | 'hc-black';
  readonly?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  language: 'json',
  height: '300px',
  theme: 'vs-dark',
  readonly: false,
});

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}>();

const content = ref(props.modelValue);

const editorOptions = computed(() => ({
  minimap: { enabled: false },
  formatOnPaste: true,
  formatOnType: true,
  automaticLayout: true,
  scrollBeyondLastLine: false,
  readOnly: props.readonly,
  tabSize: 2,
  fontSize: 14,
}));

watch(() => props.modelValue, (newValue) => {
  if (newValue !== content.value) {
    content.value = newValue;
  }
});

watch(content, (newValue) => {
  emit('update:modelValue', newValue);
  emit('change', newValue);
});

const handleMount = (editor: any) => {
  // 自动格式化
  setTimeout(() => {
    editor.getAction('editor.action.formatDocument')?.run();
  }, 100);
};
</script>

<style scoped>
.monaco-editor-wrapper {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.monaco-editor-toolbar {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 10;
  padding: 2px 4px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(4px);
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.monaco-editor-wrapper[data-theme="vs-dark"] .monaco-editor-toolbar {
  background: rgba(30, 30, 30, 0.8);
}
</style>
