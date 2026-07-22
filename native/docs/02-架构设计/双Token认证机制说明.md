# 双 Token 认证机制

## 会话模型

登录成功后创建独立 Session，并签发：

- Access Token：短期 JWT，用于访问 API。
- Refresh Token：长期随机令牌，只用于换取下一组 Token。

JWT 包含用户、客户端、设备、签发时间、过期时间、issuer 和唯一 `jti`。服务端只接受 HS256 且 issuer 为 `quick-admin` 的 Token。

Redis 保存 Token 哈希和 Session 关联，不保存明文 Refresh Token：

```text
auth:access:<access-jti>       -> session-id
auth:refresh:<refresh-hash>   -> session-id
auth:session:<session-id>     -> session data
auth:user-sessions:<user-id>:<client-id> -> session ids
```

因此 Access Token 即使 JWT 签名和时间仍有效，只要 Redis 会话已撤销，也不能继续访问。

## 生命周期

1. 登录：创建 Session 和一组 Token。
2. 鉴权：同时校验 JWT 和 Redis 中的活动 Access Token。
3. 刷新：原子消费 Refresh Token，撤销旧 Session，创建新 Session 与新 Token。
4. 登出：删除当前 Access Token 对应的 Session 和两个 Token 索引。
5. 禁止并发登录：新登录前撤销该用户、该客户端的所有旧 Session。

Refresh Token 不能重复使用。客户端应原子替换本地保存的两枚 Token；刷新失败时清空登录状态并重新登录。

## 安全边界

- JWT 密钥仅通过 `QUICK_ADMIN_JWT_SECRET` 注入，至少 32 字符。
- 日志不记录 Token、密码、验证码和第三方密钥。
- 浏览器客户端使用公开逻辑 ID `web-admin`，不设计无法保密的客户端 Secret。
- 修改密码、停用用户等需要全量失效会话的场景调用用户会话撤销能力。
