<template>
  <div class="error-page">
    <a-result status="500" title="500" sub-title="抱歉，服务器出错了">
      <template #extra>
        <a-space>
          <a-button @click="reload">刷新页面</a-button>
          <a-button type="primary" @click="goHome">返回首页</a-button>
        </a-space>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const reload = () => {
  window.location.reload();
};

const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.error-page {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100vh;
}
</style>
