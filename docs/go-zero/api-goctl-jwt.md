# 学习记录：改 .api 之后谁会更新、JWT 如何注入 context

> 日期：2026-08-14  
> 范围：`apps/user/api/user.api`、`domain.api`、`internal/handler`、`internal/logic`、`pkg/ctxdata`  
> 对照速查：`docs/config/api.md`（DSL 语法）

触发问题：`GET /v1/user/info` 返回 `field "id" is not set`；删掉 `GetUserInfoReq` 再跑 `goctl api go`，handler 没变。

## 1. 两套身份：token 里的「我」vs 页面传来的 id

| | 调用者（subject） | 资源（object） |
|---|---|---|
| 从哪来 | JWT middleware 注入 context | URL path / query / body |
| 谁说了算 | 服务端签发的 token | 页面要看的那条数据 |
| 典型路由 | `GET /v1/user/info`（当前用户） | `GET /v1/user/find?ids=` |
| 前端类比 | axios 拦截器里的登录态 | 路由参数 `params.id` |

**查自己**：id 只从 token 取，请求里不要带 id。  
**查别人**：id 来自页面，再用 token 做鉴权。

不要做成「有 id 查别人、没 id 查自己」。前端可改请求里的 id，这就是 **IDOR**（Insecure Direct Object Reference）：接口直接用客户端传来的资源 id 查/改数据，却没验证当前用户有没有权碰它。

本仓库现状：`get /info returns (GetUserInfoResp)`，logic 用 `ctx.Value("user_id")`。

## 2. 改 .api 之后，哪些文件会更新

`goctl api go` **不是全量覆盖**。看文件头：

| 文件头 | 代表文件 | 再跑 goctl |
|---|---|---|
| `DO NOT EDIT` | `types/types.go`、`handler/routes.go` | **每次覆盖** |
| `Safe to edit` | `handler/*`、`logic/*`、`svc/`、`config/` | **已存在就跳过** |

这次删掉 `GetUserInfoReq`：

- `types.go` 里类型没了（codegen 已更新）
- `get_user_info_handler.go` 仍 `Parse` 旧 Req → **编译失败 / 行为陈旧**
- `get_user_info_logic.go` 签名也不会自动改

契约变了（增删 Req、改 handler 入参）时：

1. 改 `.api` / `domain.api`
2. 跑 `goctl api go`
3. **手改** 对应 handler / logic 的签名与 Parse

强制重生成只能删掉已有 handler/logic 再跑 goctl，会丢掉已写业务，一般不这么做。

前端类比：改了 OpenAPI 再 codegen，DTO 会更新，手写过的 service 不会被抹掉。

## 3. 日常改哪、契约变了改哪

| 意图 | 改哪里 |
|---|---|
| 业务（查 RPC、鉴权后的规则） | 只改 `logic/` |
| 增删字段 / 改路由 / 是否 JWT | 改 `.api` → goctl → 必要时手改 handler/logic |
| HTTP 协议层（Cookie、上传、特殊状态码） | 改 `handler/`，保持薄 |
| 注入 RPC / DB | `svc/service_context.go` |
| 统一错误 JSON | `main` 里 `httpx.SetErrorHandlerCtx`，不要每个 handler 复制 |

**handler 可以改，但只当 HTTP 适配器**（Parse → 调 logic → 写响应）。  
该改的例子：去掉 Req 后不再 Parse；登录 `Set-Cookie`；multipart 上传。  
不该改的例子：从 token 取 id、调 UserRpc —— 放 logic。

## 4. JWT 自动注入（go-zero 在做什么）

链路：

```text
.api 里 jwt: JwtAuth
  → routes.go：rest.WithJwt(AccessSecret)
  → middleware 校验 Authorization
  → 自定义 claims 写入 r.Context()
  → logic：ctx.Value("claim 的 json key")
```

配置对齐：

| 位置 | 名字 |
|---|---|
| `.api` `@server jwt:` | `JwtAuth` |
| `config.Config` | `JwtAuth.AccessSecret` |
| `etc/dev/user.yaml` | `JwtAuth.AccessSecret` |

签发在 `pkg/ctxdata.TokenManager`，claims 的 json 名必须和读取时的 key 一致：

```go
UserID string `json:"user_id"`
Phone  string `json:"phone"`
```

logic 读取：

```go
userID, ok := l.ctx.Value("user_id").(string)
```

### 4.1 哪些字段会进 context

go-zero 把 JWT 当 `MapClaims` 遍历，**标准字段故意不注入**：

`iss` / `exp` / `nbf` / `iat` / `aud` / `sub` / `jti`

只注入自定义字段（本仓库：`user_id`、`phone`）。

前端类比：axios 拦截器校验登录态后，把 `userId` 挂到「本次请求的上下文」，组件用 `useUser()`，不必每次自己 parse token。

### 4.2 请求怎么带 token

middleware 用 `Authorization` 头（`Bearer <token>` 或裸 token 均可）。  
没带 / 签错 / 过期 → **401**，进不了 handler。  
yaml 的 `AccessSecret` 必须和签发时的 secret 相同。

## 5. 为什么当时没有「错误日志」

`field "id" is not set` 发生在 handler 的 `httpx.Parse`：**logic 根本没跑**，所以 `fmt.Println` 不会出现。

另外：

- `httpx.ErrorCtx` 只写 HTTP 响应，自己不 `logx.Error`
- go-zero access log：**< 500 打 Info**，**>= 500 才打 Error**
- 400 是客户端缺参，不是服务端故障

缺必填 query（`form:"id"` 无 `optional`）会在进业务前被拦住。无请求体接口应写成 `get /info returns (GetUserInfoResp)`，id 从 JWT context 取。

## 6. 口诀

```text
查自己：id 从 token；查别人：id 从 URL，再用 token 鉴权
改 .api → goctl 只覆盖 types/routes；handler/logic 手改签名
jwt: JwtAuth 对齐 yaml；ctx.Value 用 json key（user_id）
标准 claims 不进 context；400 是 Info 不是 Error
handler 保持薄；业务在 logic
```

## 关键词

`goctl incremental generation` · `IDOR` · `jwt MapClaims` · `context.WithValue` · `httpx.Parse` · `subject vs resource`
