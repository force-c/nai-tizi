<template>
  <div class="login-log-page">
    <BasicTable
      :columns="columns"
      :api="loadData"
      :use-search-form="true"
      :form-config="searchFormConfig"
      @register="registerTable"
    >
      <template #toolbar>
        <a-button v-permission="'log.delete'" danger @click="handleClean">
          <template #icon><DeleteOutlined /></template>
          清空日志
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status'">
          <a-tag :color="record.status === 0 ? 'success' : 'error'">
            {{ record.status === 0 ? '成功' : '失败' }}
          </a-tag>
        </template>
      </template>
    </BasicTable>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { DeleteOutlined } from '@ant-design/icons-vue';
import BasicTable from '@/components/Table/BasicTable.vue';
import { logApi } from '@/api/log';
import type { FormConfig } from '@/types/form';

const columns = [
  { title: '用户名', dataIndex: 'userName', width: 120 },
  { title: 'IP地址', dataIndex: 'ipaddr', width: 150 },
  { title: '登录地点', dataIndex: 'loginLocation', width: 150 },
  { title: '浏览器', dataIndex: 'browser', width: 120 },
  { title: '操作系统', dataIndex: 'os', width: 120 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '提示消息', dataIndex: 'msg', width: 200 },
  { title: '登录时间', dataIndex: 'loginTime', width: 180 },
];

const searchFormConfig: FormConfig = {
  schemas: [
    { field: 'userName', label: '用户名', component: 'Input', colProps: { span: 6 } },
    { field: 'ipaddr', label: 'IP地址', component: 'Input', colProps: { span: 6 } },
    { 
      field: 'status', 
      label: '状态', 
      component: 'Select', 
      componentProps: { 
        options: [
          { label: '全部', value: null },
          { label: '成功', value: 0 }, 
          { label: '失败', value: 1 }
        ],
        allowClear: true,
        placeholder: '请选择状态'
      }, 
      colProps: { span: 6 } 
    },
  ],
};

const tableRef = ref();
const registerTable = (methods: any) => { tableRef.value = methods; };

const loadData = async (params: any) => {
  return await logApi.loginLogPage(params);
};

const handleClean = () => {
  Modal.confirm({
    title: '确认清空',
    content: '确定要清空所有登录日志吗？此操作不可恢复！',
    async onOk() {
      try {
        await logApi.cleanLoginLog();
        message.success('清空成功');
        tableRef.value?.reload();
      } catch (error) {
        console.error('清空失败:', error);
      }
    },
  });
};
</script>

<style scoped lang="less">
.login-log-page {
  padding: 0;
}
</style>
