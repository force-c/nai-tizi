# Role 模块 API 响应示例

## 修改说明

已将 role 模块的所有接口统一使用 `response.Success` 返回数据，去除了冗余的 message 字段。

### 分页接口 `/api/v1/role/page` 

**修改前：**
```json
{
  "code": 200,
  "msg": "分页查询角色列表成功",
  "rows": [...],
  "total": 100,
  "pageNum": 1,
  "pageSize": 10
}
```

**修改后：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "records": [
      {
        "id": 1,
        "roleKey": "admin",
        "roleName": "管理员",
        "sort": 1,
        "status": 1,
        "dataScope": 1,
        "isSystem": true,
        "remark": "系统管理员",
        "createTime": "2024-01-01T00:00:00Z",
        "updateTime": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "size": 10,
    "current": 1,
    "pages": 10
  }
}
```

### 其他接口示例

#### 创建角色 `POST /api/v1/role`

**修改前：**
```json
{
  "code": 200,
  "msg": "创建角色成功",
  "data": { ... }
}
```

**修改后：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "roleKey": "admin",
    "roleName": "管理员",
    ...
  }
}
```

#### 更新角色 `PUT /api/v1/role`

**修改前：**
```json
{
  "code": 200,
  "msg": "更新角色成功",
  "data": null
}
```

**修改后：**
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 删除角色 `DELETE /api/v1/role/{roleId}`

**修改前：**
```json
{
  "code": 200,
  "msg": "删除角色成功",
  "data": null
}
```

**修改后：**
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

#### 获取角色详情 `GET /api/v1/role/{roleId}`

**修改前：**
```json
{
  "code": 200,
  "msg": "获取角色详情成功",
  "data": { ... }
}
```

**修改后：**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "roleKey": "admin",
    "roleName": "管理员",
    ...
  }
}
```

## 主要改动

1. **Service 层**：`Page` 方法直接返回 `*pagination.Page[model.Role]` 而不是 `([]model.Role, int64, error)`
2. **Controller 层**：所有接口统一使用 `response.Success(ctx, data)` 替代 `response.SuccessWithMsg(ctx, msg, data)`
3. **分页接口**：直接使用 `Paginator.Find()` 的第一个返回值，包含完整的分页信息（records, total, size, current, pages）

## 优势

- 统一的响应格式，前端处理更简单
- 分页数据结构更清晰，包含更多有用信息（如总页数 pages）
- 代码更简洁，减少重复的 message 定义
- 符合 RESTful 规范，成功响应使用简单的 "success" 消息
