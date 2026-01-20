<template>
  <a-config-provider :locale="zhCN" :theme="themeStore.antdTheme">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import { useThemeStore } from '@/stores/theme';

dayjs.locale('zh-cn');

const themeStore = useThemeStore();

onMounted(() => {
  themeStore.initTheme();
});
</script>

<style>
#app {
  width: 100%;
  height: 100%;
}
</style>
