<template>
  <div class="file-page">
    <BasicTable
      :columns="columns"
      :api="loadData"
      :use-search-form="true"
      :form-config="searchFormConfig"
      :show-action-column="true"
      @register="registerTable"
    >
      <template #toolbar>
        <a-upload
          v-permission="'attachment.upload'"
          :custom-request="handleUpload"
          :show-upload-list="false"
          :multiple="true"
        >
          <a-button type="primary">
            <template #icon><UploadOutlined /></template>
            上传文件
          </a-button>
        </a-upload>
        <a-button v-permission="'attachment.delete'" danger @click="handleBatchDelete">
          <template #icon><DeleteOutlined /></template>
          批量删除
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'fileName'">
          <a @click="handlePreview(record)">{{ record.fileName }}</a>
        </template>
        <template v-else-if="column.dataIndex === 'fileSize'">
          {{ formatFileSize(record.fileSize) }}
        </template>
        <template v-else-if="column.dataIndex === 'fileType'">
          <a-tag>{{ record.fileType }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <TableAction
            :actions="[
              {
                label: '预览',
                onClick: () => handlePreview(record),
                ifShow: isPreviewable(record.fileType),
              },
              {
                label: '下载',
                onClick: () => handleDownload(record),
                ifShow: hasPermission('attachment.download'),
              },
              {
                label: '删除',
                color: 'error',
                popConfirm: {
                  title: '确定删除该文件吗？',
                  onConfirm: () => handleDelete(record.attachmentId),
                },
                ifShow: hasPermission('attachment.delete'),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>

    <!-- 预览弹窗 -->
    <a-modal
      v-model:open="previewVisible"
      title="文件预览"
      :width="800"
      :footer="null"
    >
      <div v-if="previewFile" class="preview-container">
        <img
          v-if="previewFile.fileType.startsWith('image/')"
          :src="previewUrl"
          style="max-width: 100%"
        />
        <iframe
          v-else-if="previewFile.fileType === 'application/pdf'"
          :src="previewUrl"
          style="width: 100%; height: 600px; border: none"
        />
        <div v-else>
          <a-result
            status="info"
            title="该文件类型不支持预览"
            sub-title="请下载后查看"
          >
            <template #extra>
              <a-button type="primary" @click="handleDownload(previewFile)">
                下载文件
              </a-button>
            </template>
          </a-result>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { message } from 'ant-design-vue';
import { UploadOutlined, DeleteOutlined } from '@ant-design/icons-vue';
import BasicTable from '@/components/Table/BasicTable.vue';
import TableAction from '@/components/Table/TableAction.vue';
import { attachmentApi } from '@/api/attachment';
import { usePermission } from '@/composables/usePermission';
import type { FormConfig } from '@/types/components';

const { hasPermission } = usePermission();

// 表格列配置
const columns = [
  {
    title: '文件名',
    dataIndex: 'fileName',
    width: 200,
  },
  {
    title: '文件大小',
    dataIndex: 'fileSize',
    width: 120,
  },
  {
    title: '文件类型',
    dataIndex: 'fileType',
    width: 150,
  },
  {
    title: '业务类型',
    dataIndex: 'businessType',
    width: 120,
  },
  {
    title: '上传人',
    dataIndex: 'uploadBy',
    width: 120,
  },
  {
    title: '上传时间',
    dataIndex: 'uploadTime',
    width: 180,
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 200,
    fixed: 'right',
  },
];

// 搜索表单配置
const searchFormConfig: FormConfig = {
  schemas: [
    {
      field: 'fileName',
      label: '文件名',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'fileType',
      label: '文件类型',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'businessType',
      label: '业务类型',
      component: 'Input',
      colProps: { span: 6 },
    },
  ],
};

// 表格实例
const tableRef = ref();
const registerTable = (methods: any) => {
  tableRef.value = methods;
};

// 加载数据
const loadData = async (params: any) => {
  const res = await attachmentApi.page(params);
  return res;
};

// 上传文件
const handleUpload = async ({ file }: any) => {
  try {
    message.loading({ content: '上传中...', key: 'upload' });
    await attachmentApi.uploadFile(file);
    message.success({ content: '上传成功', key: 'upload' });
    tableRef.value?.reload();
  } catch (error) {
    message.error({ content: '上传失败', key: 'upload' });
    console.error('上传失败:', error);
  }
};

// 预览相关
const previewVisible = ref(false);
const previewFile = ref<any>(null);
const previewUrl = ref('');

// 判断是否可预览
const isPreviewable = (fileType: string) => {
  return fileType.startsWith('image/') || fileType === 'application/pdf';
};

// 预览文件
const handlePreview = async (record: any) => {
  try {
    const res = await attachmentApi.getUrl(record.attachmentId);
    previewUrl.value = res.url;
    previewFile.value = record;
    previewVisible.value = true;
  } catch (error) {
    console.error('获取预览URL失败:', error);
  }
};

// 下载文件
const handleDownload = async (record: any) => {
  try {
    await attachmentApi.download(record.attachmentId);
  } catch (error) {
    console.error('下载失败:', error);
  }
};

// 删除文件
const handleDelete = async (attachmentId: number) => {
  try {
    await attachmentApi.delete(attachmentId);
    message.success('删除成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

// 批量删除
const handleBatchDelete = async () => {
  const selectedRows = tableRef.value?.getSelectRows();
  if (!selectedRows || selectedRows.length === 0) {
    message.warning('请选择要删除的文件');
    return;
  }
  try {
    const attachmentIds = selectedRows.map((row: any) => row.attachmentId);
    await attachmentApi.batchDelete(attachmentIds);
    message.success('批量删除成功');
    tableRef.value?.reload();
  } catch (error) {
    console.error('批量删除失败:', error);
  }
};

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
};
</script>

<style scoped lang="less">
.file-page {
  padding: 0;
}

.preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
