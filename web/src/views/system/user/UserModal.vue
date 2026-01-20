<template>
  <BasicModal
    v-model:visible="visible"
    :title="isEdit ? '编辑用户' : '新增用户'"
    :width="600"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <BasicForm
      ref="formRef"
      :schemas="formSchemas"
      :model="formData"
      :label-width="100"
      :show-action-buttons="false"
    />
  </BasicModal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import BasicModal from '@/components/Modal/BasicModal.vue';
import BasicForm from '@/components/Form/BasicForm.vue';
import { userApi } from '@/api/user';
import { organizationApi } from '@/api/organization';
import type { FormSchema } from '@/types/form';

const props = defineProps<{
  visible: boolean;
  userId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'success'): void;
}>();

const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

const isEdit = computed(() => !!props.userId);
const loading = ref(false);
const formRef = ref();
const formData = ref<Record<string, any>>({});
const orgTreeData = ref<any[]>([]);

// 表单配置
const formSchemas = computed<FormSchema[]>(() => [
  {
    field: 'orgId',
    label: '所属组织',
    component: 'TreeSelect',
    rules: [{ required: true, message: '请选择所属组织' }],
    componentProps: {
      treeData: orgTreeData.value,
      placeholder: '请选择所属组织',
      treeDefaultExpandAll: true,
      fieldNames: { label: 'title', value: 'value', children: 'children' },
    },
  },
  {
    field: 'username',
    label: '用户名',
    component: 'Input',
    rules: [
      { required: true, message: '请输入用户名' },
      { min: 3, message: '用户名长度不能少于3位' },
    ],
    componentProps: {
      disabled: isEdit.value,
    },
  },
  {
    field: 'nickname',
    label: '昵称',
    component: 'Input',
    rules: [{ required: true, message: '请输入昵称' }],
  },
  {
    field: 'password',
    label: '密码',
    component: 'Input',
    rules: [
      { required: !isEdit.value, message: '请输入密码' },
      { min: 6, message: '密码长度不能少于6位' },
    ],
    componentProps: {
      type: 'password',
    },
    show: !isEdit.value,
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'Input',
    rules: [
      { type: 'email', message: '请输入正确的邮箱格式' },
    ],
  },
  {
    field: 'phonenumber',
    label: '手机号',
    component: 'Input',
    rules: [
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' },
    ],
  },
  {
    field: 'sex',
    label: '性别',
    component: 'Select',
    componentProps: {
      options: [
        { label: '男', value: 0 },
        { label: '女', value: 1 },
        { label: '未知', value: 2 },
      ],
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '正常', value: 0 },
        { label: '停用', value: 1 },
      ],
    },
  },
  {
    field: 'remark',
    label: '备注',
    component: 'Textarea',
    componentProps: {
      rows: 3,
    },
  },
]);

// 加载用户详情
const loadUserDetail = async () => {
  if (!props.userId) return;
  
  try {
    loading.value = true;
    const data = await userApi.detail(props.userId);
    formData.value = data;
    formRef.value?.setFieldsValue(data);
  } catch (error) {
    console.error('加载用户详情失败:', error);
  } finally {
    loading.value = false;
  }
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    loading.value = true;
    
    const values = formRef.value?.getFieldsValue();
    
    if (isEdit.value) {
      await userApi.update(props.userId!, values);
      message.success('操作成功');
    } else {
      await userApi.create(values);
      message.success('操作成功');
    }
    
    emit('success');
  } catch (error: any) {
    console.error('提交失败:', error);
    
    // 检查错误对象中的 code 和 message
    let errorCode: number | undefined;
    let errorMsg: string | undefined;
    
    // 优先从 error.code 和 error.message 获取（request 拦截器附加的）
    if (error?.code !== undefined) {
      errorCode = error.code;
      errorMsg = error.message || error.msg;
    } 
    // 其次从 error.response.data 获取（axios 响应错误）
    else if (error?.response?.data) {
      const responseData = error.response.data;
      errorCode = responseData.code;
      errorMsg = responseData.message || responseData.msg;
    }
    
    // 根据 code 显示错误消息
    if (errorCode !== undefined && errorCode !== 200) {
      // code 500 或其他错误码，显示后台返回的 msg
      message.error(errorMsg || '操作失败');
    } else {
      // 其他错误（网络错误等）
      message.error(errorMsg || error?.message || '操作失败，请稍后重试');
    }
  } finally {
    loading.value = false;
  }
};

// 取消
const handleCancel = () => {
  formRef.value?.resetFields();
  formData.value = {};
};

// 转换组织树数据，添加 key 属性
const transformOrgTree = (orgs: any[]): any[] => {
  return orgs.map(org => {
    // 后端返回的字段是 id，前端类型定义是 orgId，需要兼容处理
    const orgId = org.id !== undefined ? org.id : org.orgId;
    return {
      ...org,
      key: orgId,
      value: orgId,
      title: org.orgName,
      children: org.children ? transformOrgTree(org.children) : undefined
    };
  });
};

// 加载组织树数据
const loadOrgTree = async () => {
  try {
    const data = await organizationApi.tree();
    orgTreeData.value = transformOrgTree(data);
  } catch (error) {
    console.error('加载组织树失败:', error);
  }
};

// 监听弹窗显示
watch(
  () => props.visible,
  (val) => {
    if (val) {
      if (props.userId) {
        loadUserDetail();
      } else {
        formRef.value?.resetFields();
        formData.value = {};
      }
    }
  }
);

// 组件挂载时加载组织树
onMounted(() => {
  loadOrgTree();
});
</script>
