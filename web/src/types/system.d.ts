// 系统管理相关类型定义

// 用户信息（扩展）
export interface User {
  userId: number;
  orgId?: number;
  username: string;
  nickname: string;
  email?: string;
  phonenumber?: string;
  sex?: number;
  avatar?: string;
  password?: string;
  userType: number;
  status: number;
  remark?: string;
  createBy?: number;
  updateBy?: number;
  createdTime?: string;
  updatedTime?: string;
}

// 角色信息
export interface Role {
  roleId: number;
  roleName: string;
  roleKey: string;
  roleSort: number;
  dataScope: number;
  status: number;
  remark?: string;
  createTime?: string;
  updateTime?: string;
}

// 组织信息
export interface Organization {
  id?: number; // 后端返回的字段名
  orgId?: number; // 兼容字段，前端可能使用
  orgName: string;
  parentId: number;
  orderNum?: number;
  sort?: number; // 后端返回的字段名
  leader?: string;
  phone?: string;
  email?: string;
  status: number;
  children?: Organization[];
}

// 存储环境信息
export interface StorageEnv {
  id: number;
  name: string;
  code: string;
  storageType: string;
  config: Record<string, any>;
  isDefault: boolean;
  status: number;
  remark?: string;
  createTime?: string;
  updateTime?: string;
}

// 附件信息
export interface Attachment {
  attachmentId: number;
  fileName: string;
  fileSize: number;
  fileType: string;
  filePath: string;
  fileUrl: string;
  envId: number;
  businessType?: string;
  businessId?: number;
  createTime?: string;
}


// 字典信息
export interface Dict {
  dictId: number;
  dictType: string;
  dictLabel: string;
  dictValue: string;
  dictSort: number;
  cssClass?: string;
  listClass?: string;
  isDefault?: boolean;
  status: number;
  remark?: string;
  createTime?: string;
  updateTime?: string;
}

// 配置信息
export interface Config {
  configId: number;
  configName: string;
  configKey: string;
  configValue: string;
  configType: string;
  remark?: string;
  createTime?: string;
  updateTime?: string;
}

// 登录日志
export interface LoginLog {
  infoId: number;
  username: string;
  ipaddr: string;
  loginLocation?: string;
  browser?: string;
  os?: string;
  status: number;
  msg?: string;
  loginTime: string;
}

// 操作日志
export interface OperLog {
  operId: number;
  title: string;
  businessType: string;
  method: string;
  requestMethod: string;
  operatorType: string;
  operName: string;
  deptName?: string;
  operUrl: string;
  operIp: string;
  operLocation?: string;
  operParam?: string;
  jsonResult?: string;
  status: number;
  errorMsg?: string;
  operTime: string;
}
