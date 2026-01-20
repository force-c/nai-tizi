<template>
  <a-space :size="8">
    <template v-for="(action, index) in getActions" :key="index">
      <a-popconfirm
        v-if="action.popConfirm"
        :title="action.popConfirm.title"
        @confirm="action.popConfirm.confirm"
        @cancel="action.popConfirm.cancel"
      >
        <a-button
          :type="action.color === 'primary' ? 'primary' : 'link'"
          :danger="action.color === 'error'"
          :disabled="action.disabled"
          size="small"
        >
          {{ action.label }}
        </a-button>
      </a-popconfirm>
      <a-button
        v-else
        :type="action.color === 'primary' ? 'primary' : 'link'"
        :danger="action.color === 'error'"
        :disabled="action.disabled"
        size="small"
        @click="action.onClick"
      >
        {{ action.label }}
      </a-button>
    </template>
  </a-space>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Action {
  label: string;
  onClick?: () => void;
  color?: 'default' | 'primary' | 'success' | 'warning' | 'error';
  disabled?: boolean;
  loading?: boolean;
  ifShow?: boolean;
  popConfirm?: {
    title: string;
    confirm?: () => void;
    cancel?: () => void;
    okText?: string;
    cancelText?: string;
  };
}

interface Props {
  actions?: Action[];
  divider?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  actions: () => [],
  divider: true,
});

const getActions = computed(() => {
  return props.actions.filter((action) => action.ifShow !== false);
});
</script>
