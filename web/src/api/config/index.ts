import { request } from '@/utils/request';
import type { PageParams, PageResponse } from '@/types/api';

export interface Config {
  id: number;
  name: string;
  code: string;
  data?: any;
  remark?: string;
  createBy?: number;
  updateBy?: number;
  createdTime?: string;
  updatedTime?: string;
}

// 创建配置请求参数
export interface CreateConfigParams {
  name: string;
  code: string;
  data?: any;
  remark?: string;
  createBy: number;
  updateBy: number;
}

// 更新配置请求参数
export interface UpdateConfigParams {
  id: number;
  name: string;
  code: string;
  data?: any;
  remark?: string;
  updateBy: number;
}

export const configApi = {
  // 获取配置列表（分页）
  page: (params: PageParams & { name?: string; code?: string }) =>
    request.post<PageResponse<Config>>('/api/v1/config/page', params),

  // 根据配置编码获取配置列表
  getByCode: (code: string) =>
    request.get<Config[]>(`/api/v1/config/code`, { params: { code } }),

  // 根据配置编码获取配置数据
  getData: (code: string) =>
    request.get<any>(`/api/v1/config/data`, { params: { code } }),

  // 获取配置详情
  detail: (id: number) =>
    request.get<Config>(`/api/v1/config/${id}`),

  // 创建配置
  create: (data: CreateConfigParams) =>
    request.post('/api/v1/config', data),

  // 更新配置
  update: (data: UpdateConfigParams) =>
    request.put(`/api/v1/config/${data.id}`, data),

  // 删除配置
  delete: (id: number) =>
    request.delete(`/api/v1/config/${id}`),

  // 批量删除配置
  batchDelete: (ids: number[]) =>
    request.delete('/api/v1/config/batch', { data: { ids } }),
};
