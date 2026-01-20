<template>
  <a-popover
    v-model:visible="visible"
    trigger="click"
    placement="bottomLeft"
    overlay-class-name="icon-picker-popover"
  >
    <template #content>
      <div class="icon-picker-content">
        <a-input
          v-model:value="searchText"
          placeholder="搜索图标..."
          allow-clear
          class="mb-2"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </a-input>
        
        <div class="icon-list">
          <div
            v-for="iconName in filteredIcons"
            :key="iconName"
            class="icon-item"
            :class="{ active: modelValue === iconName }"
            @click="handleSelect(iconName)"
          >
            <component :is="getIcon(iconName)" class="icon" />
            <span class="icon-name">{{ iconName }}</span>
          </div>
        </div>
      </div>
    </template>
    
    <a-input
      :value="modelValue"
      placeholder="请选择图标"
      readonly
      :disabled="disabled"
      class="icon-picker-input"
    >
      <template #prefix>
        <component :is="getIcon(modelValue)" v-if="modelValue" />
        <AppstoreOutlined v-else />
      </template>
      <template #suffix>
        <CloseCircleOutlined
          v-if="modelValue && !disabled"
          @click.stop="handleClear"
          class="clear-icon"
        />
      </template>
    </a-input>
  </a-popover>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { SearchOutlined, AppstoreOutlined, CloseCircleOutlined } from '@ant-design/icons-vue';
import { getIcon, getAvailableIcons } from '@/config/icons';

interface Props {
  modelValue?: string;
  disabled?: boolean;
}

interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  disabled: false,
});

const emit = defineEmits<Emits>();

const visible = ref(false);
const searchText = ref('');

// 所有可用图标
const allIcons = getAvailableIcons();

// 过滤后的图标列表
const filteredIcons = computed(() => {
  if (!searchText.value) return allIcons;
  const search = searchText.value.toLowerCase();
  return allIcons.filter(icon => icon.toLowerCase().includes(search));
});

// 选择图标
const handleSelect = (iconName: string) => {
  emit('update:modelValue', iconName);
  emit('change', iconName);
  visible.value = false;
  searchText.value = '';
};

// 清除选择
const handleClear = () => {
  emit('update:modelValue', '');
  emit('change', '');
};
</script>

<style scoped lang="less">
.icon-picker-input {
  cursor: pointer;
  
  :deep(.ant-input) {
    cursor: pointer;
  }
  
  :deep(.ant-input[disabled]) {
    cursor: not-allowed;
  }
  
  .clear-icon {
    color: rgba(0, 0, 0, 0.25);
    cursor: pointer;
    transition: color 0.3s;
    
    &:hover {
      color: rgba(0, 0, 0, 0.45);
    }
  }
}
</style>

<style lang="less">
// 使用全局样式确保弹出层样式生效
.icon-picker-popover {
  .ant-popover-inner {
    padding: 12px;
  }
  
  .icon-picker-content {
    width: 500px;
    
    .mb-2 {
      margin-bottom: 12px;
    }
    
    .icon-list {
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 8px;
      max-height: 400px;
      overflow-y: auto;
      padding: 8px 0;
      
      &::-webkit-scrollbar {
        width: 6px;
      }
      
      &::-webkit-scrollbar-thumb {
        background-color: rgba(0, 0, 0, 0.2);
        border-radius: 3px;
        
        &:hover {
          background-color: rgba(0, 0, 0, 0.3);
        }
      }
      
      .icon-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 12px 8px;
        border: 1px solid #d9d9d9;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.3s;
        min-height: 70px;
        background-color: #fff;
        
        &:hover {
          border-color: #1890ff;
          background-color: #e6f7ff;
        }
        
        &.active {
          border-color: #1890ff;
          background-color: #e6f7ff;
          box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
        }
        
        .icon {
          font-size: 24px;
          margin-bottom: 6px;
          color: rgba(0, 0, 0, 0.85);
        }
        
        .icon-name {
          font-size: 11px;
          color: rgba(0, 0, 0, 0.65);
          text-align: center;
          word-break: break-all;
          line-height: 1.3;
          max-width: 100%;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }
      }
    }
  }
}
</style>
