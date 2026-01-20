<template>
  <div class="oper-log-page">
    <BasicTable
      :columns="columns"
      :api="loadData"
      :use-search-form="true"
      :form-config="searchFormConfig"
      :scroll="{ x: 2500 }"
      @register="registerTable"
    >
      <template #toolbar>
        <a-button v-permission="'log.delete'" danger @click="handleClean">
          <template #icon><DeleteOutlined /></template>
          清空日志
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'requestMethod'">
          <a-tag :color="getMethodColor(record.requestMethod)">
            {{ record.requestMethod }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'deviceType'">
          <a-tag :color="getDeviceColor(record.deviceType)">
            {{ record.deviceType }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <a-tag :color="record.status === '0' || record.status === 0 ? 'success' : 'error'">
            {{ record.status === '0' || record.status === 0 ? '成功' : '失败' }}
          </a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <a-button type="link" size="small" @click="handleView(record)">详情</a-button>
        </template>
      </template>
    </BasicTable>

    <!-- 详情弹窗 -->
    <a-modal
      v-model:open="detailVisible"
      title="操作日志详情"
      width="800px"
      :footer="null"
    >
      <a-descriptions :column="2" bordered>
        <a-descriptions-item label="系统模块">{{ currentRecord?.title }}</a-descriptions-item>
        <a-descriptions-item label="业务类型">{{ currentRecord?.businessType }}</a-descriptions-item>
        <a-descriptions-item label="请求方式">{{ currentRecord?.requestMethod }}</a-descriptions-item>
        <a-descriptions-item label="终端类型">{{ currentRecord?.deviceType }}</a-descriptions-item>
        <a-descriptions-item label="用户">{{ currentRecord?.operName }}</a-descriptions-item>
        <a-descriptions-item label="IP">{{ currentRecord?.operIp }}</a-descriptions-item>
        <a-descriptions-item label="地点">{{ currentRecord?.operLocation || '-' }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="currentRecord?.status === '0' || currentRecord?.status === 0 ? 'success' : 'error'">
            {{ currentRecord?.status === '0' || currentRecord?.status === 0 ? '成功' : '失败' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="时间" :span="2">{{ currentRecord?.operTime }}</a-descriptions-item>
        <a-descriptions-item label="耗时">{{ currentRecord?.costTime }}ms</a-descriptions-item>
        <a-descriptions-item label="请求URL" :span="2">{{ currentRecord?.operUrl }}</a-descriptions-item>
        <a-descriptions-item label="调用方法" :span="2">{{ currentRecord?.method }}</a-descriptions-item>
        <a-descriptions-item label="请求参数" :span="2">
          <pre style="max-height: 200px; overflow: auto;">{{ currentRecord?.operParam || '-' }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="返回结果" :span="2">
          <pre style="max-height: 200px; overflow: auto;">{{ currentRecord?.jsonResult || '-' }}</pre>
        </a-descriptions-item>
        <a-descriptions-item label="错误信息" :span="2">
          <span style="color: red;">{{ currentRecord?.errorMsg || '-' }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="User-Agent" :span="2">{{ currentRecord?.userAgent }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
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
  { title: '系统模块', dataIndex: 'title', width: 120 },
  { title: '业务类型', dataIndex: 'businessType', width: 100 },
  { title: '请求方式', dataIndex: 'requestMethod', width: 100 },
  { title: 'URL', dataIndex: 'operUrl', width: 250, ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '耗时(ms)', dataIndex: 'costTime', width: 100 },
  { title: '终端类型', dataIndex: 'deviceType', width: 100 },
  { title: 'IP', dataIndex: 'operIp', width: 130 },
  { title: '地点', dataIndex: 'operLocation', width: 150 },
  { title: '用户', dataIndex: 'operName', width: 180 },
  { title: '时间', dataIndex: 'operTime', width: 180 },
  { title: '操作', dataIndex: 'action', width: 100, fixed: 'right' },
];

const searchFormConfig: FormConfig = {
  schemas: [
    { field: 'title', label: '系统模块', component: 'Input', colProps: { span: 6 } },
    { field: 'operName', label: '用户', component: 'Input', colProps: { span: 6 } },
    { field: 'businessType', label: '业务类型', component: 'Input', colProps: { span: 6 } },
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

const detailVisible = ref(false);
const currentRecord = ref<any>(null);

const loadData = async (params: any) => {
  return await logApi.operLogPage(params);
};

const handleView = (record: any) => {
  currentRecord.value = record;
  detailVisible.value = true;
};

const handleClean = () => {
  Modal.confirm({
    title: '确认清空',
    content: '确定要清空所有操作日志吗？此操作不可恢复！',
    async onOk() {
      try {
        await logApi.cleanOperLog();
        message.success('清空成功');
        tableRef.value?.reload();
      } catch (error) {
        console.error('清空失败:', error);
      }
    },
  });
};

// 请求方式颜色映射
const getMethodColor = (method: string) => {
  const colorMap: Record<string, string> = {
    'GET': 'blue',
    'POST': 'green',
    'PUT': 'orange',
    'DELETE': 'red',
    'PATCH': 'purple',
  };
  return colorMap[method] || 'default';
};

// 终端类型颜色映射
const getDeviceColor = (deviceType: string) => {
  const colorMap: Record<string, string> = {
    'PC': 'cyan',
    'Mobile': 'magenta',
    'Tablet': 'geekblue',
    'Unknown': 'default',
  };
  return colorMap[deviceType] || 'default';
};
</script>

<style scoped lang="less">
.oper-log-page {
  padding: 0;
}
</style>
