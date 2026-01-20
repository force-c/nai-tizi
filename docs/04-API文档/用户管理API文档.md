# 用户管理 API 文档

## 概述

用户管理模块提供了完整的用户 CRUD 操作，包括新增、修改、删除、批量删除、批量导入和分页查询功能。

## 认证要求

所有用户管理接口都需要认证，请在请求头中携带 Token：

```
Authorization: Bearer {access_token}
```

## API 接口

### 1. 创建用户

**接口**: `POST /user`

**请求体**:
```json
{
  "orgId": 1,
  "userName": "zhangsan",
  "nickName": "张三",
  "password": "123456",
  "userType": "sys_user",
  "email": "zhangsan@example.com",
  "phonenumber": "13800138001",
  "sex": "0",
  "avatar": "",
  "status": "0",
  "remark": "测试用户"
}
```

**字段说明**:
- `orgId` (必填): 组织ID
- `userName` (必填): 用户名（登录账号）
- `nickName` (必填): 昵称（显示名称）
- `password` (必填): 密码（最少6位）
- `userType` (可选): 用户类型，默认 `sys_user`
- `email` (可选): 邮箱
- `phonenumber` (可选): 手机号
- `sex` (可选): 性别（0男 1女 2未知）
- `avatar` (可选): 头像URL
- `status` (可选): 状态（0正常 1停用），默认 `0`
- `remark` (可选): 备注

**响应**:
```json
{
  "code": 200,
  "data": {
    "userId": 123
  },
  "msg": "success"
}
```

### 2. 更新用户

**接口**: `PUT /user/:id`

**路径参数**:
- `id`: 用户ID

**请求体**:
```json
{
  "orgId": 1,
  "userName": "zhangsan",
  "nickName": "张三三",
  "userType": "sys_user",
  "email": "zhangsan@example.com",
  "phonenumber": "13800138001",
  "sex": "0",
  "avatar": "",
  "status": "0",
  "remark": "更新后的备注"
}
```

**注意**: 更新接口不会修改密码，如需修改密码请使用重置密码接口。

**响应**:
```json
{
  "code": 200,
  "data": "ok",
  "msg": "success"
}
```

### 3. 删除用户

**接口**: `DELETE /user/:id`

**路径参数**:
- `id`: 用户ID

**响应**:
```json
{
  "code": 200,
  "data": "ok",
  "msg": "success"
}
```

### 4. 批量删除用户

**接口**: `DELETE /user/batch`

**请求体**:
```json
{
  "userIds": [1, 2, 3, 4, 5]
}
```

**响应**:
```json
{
  "code": 200,
  "data": "ok",
  "msg": "success"
}
```

### 5. 根据ID查询用户

**接口**: `GET /user/:id`

**路径参数**:
- `id`: 用户ID

**响应**:
```json
{
  "code": 200,
  "data": {
    "userId": 1,
    "orgId": 1,
    "userName": "zhangsan",
    "nickName": "张三",
    "userType": "sys_user",
    "email": "zhangsan@example.com",
    "phonenumber": "13800138001",
    "sex": "0",
    "avatar": "",
    "status": "0",
    "loginIp": "127.0.0.1",
    "loginDate": 1703001234,
    "openId": "",
    "unionId": "",
    "remark": "测试用户",
    "createBy": 1,
    "createTime": "2024-01-01T10:00:00Z",
    "updateBy": 1,
    "updateTime": "2024-01-01T10:00:00Z"
  },
  "msg": "success"
}
```

### 6. 分页查询用户列表

**接口**: `GET /user/list`

**查询参数**:
- `pageNum` (可选): 页码，默认 1
- `pageSize` (可选): 每页数量，默认 10
- `username` (可选): 用户名（模糊查询）
- `phonenumber` (可选): 手机号（模糊查询）
- `status` (可选): 状态（0正常 1停用）

**示例**:
```
GET /user/list?pageNum=1&pageSize=10&username=zhang&status=0
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "userId": 1,
        "orgId": 1,
        "userName": "zhangsan",
        "nickName": "张三",
        "userType": "sys_user",
        "email": "zhangsan@example.com",
        "phonenumber": "13800138001",
        "sex": "0",
        "avatar": "",
        "status": "0",
        "loginIp": "127.0.0.1",
        "loginDate": 1703001234,
        "remark": "测试用户",
        "createTime": "2024-01-01T10:00:00Z",
        "updateTime": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "pageNum": 1,
    "pageSize": 10
  },
  "msg": "success"
}
```

### 7. 批量导入用户

**接口**: `POST /user/import`

**请求体**:
```json
{
  "users": [
    {
      "orgId": 1,
      "userName": "user001",
      "nickName": "用户001",
      "password": "123456",
      "userType": "sys_user",
      "email": "user001@example.com",
      "phonenumber": "13800138001",
      "sex": "0",
      "status": "0"
    },
    {
      "orgId": 1,
      "userName": "user002",
      "nickName": "用户002",
      "password": "123456",
      "userType": "sys_user",
      "email": "user002@example.com",
      "phonenumber": "13800138002",
      "sex": "1",
      "status": "0"
    }
  ]
}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "successCount": 2,
    "failCount": 0,
    "errors": []
  },
  "msg": "success"
}
```

**失败示例**:
```json
{
  "code": 200,
  "data": {
    "successCount": 1,
    "failCount": 1,
    "errors": [
      "第2行: 用户名已存在"
    ]
  },
  "msg": "success"
}
```

### 8. 重置用户密码

**接口**: `PUT /user/:id/password`

**路径参数**:
- `id`: 用户ID

**请求体**:
```json
{
  "newPassword": "newpass123"
}
```

**响应**:
```json
{
  "code": 200,
  "data": "ok",
  "msg": "success"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未登录或登录已过期 |
| 500 | 服务器内部错误 |

## 业务错误说明

| 错误信息 | 说明 |
|---------|------|
| 用户名已存在 | 创建用户时用户名重复 |
| 手机号已存在 | 创建用户时手机号重复 |
| 邮箱已存在 | 创建用户时邮箱重复 |
| 用户不存在 | 查询、更新或删除的用户不存在 |
| 用户名已被占用 | 更新用户时用户名被其他用户占用 |
| 手机号已被占用 | 更新用户时手机号被其他用户占用 |
| 邮箱已被占用 | 更新用户时邮箱被其他用户占用 |

## 使用示例

### cURL 示例

#### 创建用户
```bash
curl -X POST http://localhost:9009/user \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "orgId": 1,
    "userName": "zhangsan",
    "nickName": "张三",
    "password": "123456",
    "email": "zhangsan@example.com",
    "phonenumber": "13800138001"
  }'
```

#### 查询用户列表
```bash
curl -X GET "http://localhost:9009/user/list?pageNum=1&pageSize=10&username=zhang" \
  -H "Authorization: Bearer {token}"
```

#### 更新用户
```bash
curl -X PUT http://localhost:9009/user/123 \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "nickName": "张三三",
    "status": "0"
  }'
```

#### 删除用户
```bash
curl -X DELETE http://localhost:9009/user/123 \
  -H "Authorization: Bearer {token}"
```

#### 批量删除用户
```bash
curl -X DELETE http://localhost:9009/user/batch \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "userIds": [1, 2, 3]
  }'
```

### JavaScript 示例

```javascript
// 创建用户
async function createUser() {
  const response = await fetch('http://localhost:9009/user', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + token,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      orgId: 1,
      userName: 'zhangsan',
      nickName: '张三',
      password: '123456',
      email: 'zhangsan@example.com',
      phonenumber: '13800138001'
    })
  });
  const data = await response.json();
  console.log(data);
}

// 查询用户列表
async function getUserList() {
  const response = await fetch('http://localhost:9009/user/list?pageNum=1&pageSize=10', {
    headers: {
      'Authorization': 'Bearer ' + token
    }
  });
  const data = await response.json();
  console.log(data);
}

// 批量导入用户
async function importUsers() {
  const response = await fetch('http://localhost:9009/user/import', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + token,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      users: [
        {
          orgId: 1,
          userName: 'user001',
          nickName: '用户001',
          password: '123456',
          email: 'user001@example.com'
        },
        {
          orgId: 1,
          userName: 'user002',
          nickName: '用户002',
          password: '123456',
          email: 'user002@example.com'
        }
      ]
    })
  });
  const data = await response.json();
  console.log(data);
}
```

## 注意事项

1. **密码安全**: 
   - 密码最少6位
   - 密码会自动使用 bcrypt 加密存储
   - 更新用户接口不会修改密码

2. **唯一性约束**:
   - 用户名必须唯一
   - 手机号必须唯一（如果提供）
   - 邮箱必须唯一（如果提供）

3. **批量导入**:
   - 导入失败的记录不会影响其他记录
   - 返回结果包含成功数、失败数和错误详情
   - 建议分批导入，每批不超过100条

4. **分页查询**:
   - 默认按创建时间倒序排列
   - 支持用户名和手机号模糊查询
   - 支持按状态筛选

5. **审计字段**:
   - `createBy` 和 `updateBy` 会自动设置为当前登录用户ID
   - `createTime` 和 `updateTime` 由数据库自动管理

## 参考

- [认证系统使用指南](./认证系统使用指南.md)
- [开发规范指南](./开发规范指南.md)
