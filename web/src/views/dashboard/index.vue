<template>
  <div class="dashboard">
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="用户总数"
            :value="112893"
            :value-style="{ color: '#3f8600' }"
          >
            <template #prefix>
              <user-outlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="今日访问"
            :value="9280"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <eye-outlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="文件总数"
            :value="1128"
            :value-style="{ color: '#cf1322' }"
          >
            <template #prefix>
              <file-outlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="存储使用"
            :value="93"
            suffix="/ 100 GB"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <cloud-outlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="24">
        <a-card title="欢迎使用 NAI-TIZI 管理后台">
          <p>这是一个基于 Vue 3 + TypeScript + Ant Design Vue 构建的现代化管理后台系统。</p>
          <p>当前登录用户：{{ authStore.nickname || authStore.username }}</p>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { UserOutlined, EyeOutlined, FileOutlined, CloudOutlined } from '@ant-design/icons-vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
</script>

<style scoped lang="less">
.dashboard {
  background: #fff;
  padding: 24px;
  border-radius: 4px;
}
</style>
