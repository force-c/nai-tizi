import { request } from '@/utils/request';
import type { Organization } from '@/types/system';

export const organizationApi = {
  // 获取组织树
  tree: () => request.get<Organization[]>('/api/v1/org/tree'),

  // 获取组织列表（分页）
  page: (params?: { orgName?: string; status?: number; pageNum?: number; pageSize?: number }) =>
    request.post<any>('/api/v1/org/page', params),

  // 获取组织详情
  detail: (orgId: number) => request.get<Organization>(`/api/v1/org/${orgId}`),

  // 创建组织
  create: (data: Partial<Organization>) => request.post('/api/v1/org', data),

  // 更新组织
  update: (id: number, data: Partial<Organization>) =>
      request.put(`/api/v1/org/${id}`, data),

  // 删除组织
  delete: (id: number) => request.delete(`/api/v1/org/${id}`),

  // 批量删除
  batchDelete: (ids: number[]) =>
    request.delete('/api/v1/org/batch', { data: { ids } }),
};
