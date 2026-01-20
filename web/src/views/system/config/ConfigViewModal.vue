<template>
  <BasicModal
    v-model:visible="visible"
    title=""
    :width="800"
    :footer="null"
    wrap-class-name="config-view-modal"
    :body-style="{ maxHeight: '80vh', padding: '16px' }"
    @cancel="handleCancel"
  >
    <div class="config-view-content">
      <MonacoEditor
        :model-value="jsonContent"
        height="300px"
        :theme="editorTheme"
        :readonly="true"
        language="json"
      >
        <template #toolbar>
          <a-button size="small" type="text" @click="handleCopy" :title="'复制JSON数据'">
            <template #icon><CopyOutlined /></template>
          </a-button>
          <a-button size="small" type="text" @click="toggleTheme" :title="'切换主题'">
            {{ editorTheme === 'vs-dark' ? '🌙' : '☀️' }}
          </a-button>
        </template>
      </MonacoEditor>
    </div>
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { message } from 'ant-design-vue';
import { CopyOutlined } from '@ant-design/icons-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import MonacoEditor from '@/components/MonacoEditor/index.vue';

const props = defineProps<{ 
  visible: boolean; 
  data: any;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
}>();

const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const editorTheme = ref<'vs-dark' | 'vs'>('vs-dark');

const jsonContent = computed(() => {
  if (!props.data) return '';
  try {
    if (typeof props.data === 'string') {
      return JSON.stringify(JSON.parse(props.data), null, 2);
    }
    return JSON.stringify(props.data, null, 2);
  } catch (e) {
    return props.data;
  }
});

const toggleTheme = () => {
  editorTheme.value = editorTheme.value === 'vs-dark' ? 'vs' : 'vs-dark';
};

const handleCopy = async () => {
  try {
    const text = jsonContent.value;
    if (!text) {
      message.warning('没有可复制的内容');
      return;
    }
    
    // 使用 Clipboard API 复制文本
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text);
      message.success('复制成功');
    } else {
      // 降级方案：使用传统方法
      const textarea = document.createElement('textarea');
      textarea.value = text;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.select();
      try {
        document.execCommand('copy');
        message.success('复制成功');
      } catch (err) {
        message.error('复制失败');
      } finally {
        document.body.removeChild(textarea);
      }
    }
  } catch (error) {
    console.error('复制失败:', error);
    message.error('复制失败');
  }
};

const handleCancel = () => {
  visible.value = false;
};
</script>

<style scoped lang="less">
.config-view-content {
  min-height: 300px;
  
  :deep(.monaco-editor-wrapper) {
    min-height: 300px;
  }
}
</style>

<style lang="less">
.config-view-modal {
  .ant-modal {
    max-width: 90vw;
  }
  
  .ant-modal-content {
    max-height: 90vh;
    display: flex;
    flex-direction: column;
  }
  
  .ant-modal-header {
    display: none;
  }
  
  .ant-modal-body {
    flex: 1;
    overflow-y: auto;
    max-height: calc(90vh - 110px);
  }
}
</style>
