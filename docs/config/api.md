# go-zero API DSL（.api）核心速查

面向本仓库（Go + go-zero HTTP）的复习笔记。对照示例：`apps/user/api/user.api`、`apps/user/api/domain.api`。

前端类比：`.api` ≈ **OpenAPI / 路由表 + TypeScript interface**，由 `goctl api go` 生成 handler / types / routes；业务写在 `logic`，契约写在 `.api`。

和 `.proto` 的分工：`.proto` 管 gRPC 内部服务；`.api` 管对外 HTTP。本仓库 User 服务两者字段应对齐（如 Sex / Status 枚举数值）。

---

## 1. 它是什么

| 概念 | 说明 |
|------|------|
| IDL | HTTP 接口描述语言，只定义契约，不含业务实现 |
| 用途 | 路由、请求/响应结构、JWT、中间件、分组前缀 |
| 产物 | `goctl api go` 生成 `handler`、`types`、`routes.go` 等 |

流程：写 `.api` → `goctl api go` → 在 `logic` 填业务 → HTTP 入口自动绑定。

```bash
goctl api go -api user.api -dir .
```

---

## 2. 文件骨架

```api
syntax = "v1"                 // 必填，语法版本

info (                        // 可选，元数据（文档用）
	title:  "用户服务"
	author: "richard"
)

import (                      // 可选，拆分类型定义
	"domain.api"
)

@server (                     // 可选，作用于后面整个 service 块
	prefix: v1/user
	group:  user
)
service user {
	@doc "用户注册"
	@handler register
	post /register (RegisterReq) returns (RegisterResp)
}
```

| 元素 | 作用 | 前端类比 |
|------|------|----------|
| `syntax` | DSL 版本 | OpenAPI `openapi: 3.0` |
| `info` | 标题、作者等元数据 | package.json 描述信息 |
| `import` | 引用别的 `.api`（常放 type） | `import type` |
| `type` | 请求/响应结构体 | interface / DTO |
| `@server` | 前缀、分组、JWT、中间件 | 路由 `prefix` + middleware |
| `service` | HTTP 路由集合 | Router 模块 |
| `@handler` | 生成的 handler/logic 名 | 控制器方法名 |
| `@doc` | 接口说明（文档） | JSDoc / swagger summary |

本仓库拆法：

| 文件 | 内容 |
|------|------|
| `domain.api` | 只放 `type`（DTO） |
| `user.api` | `import` + `service`（路由） |

约定：被 import 的文件以类型为主；`service` 名在各文件间保持一致。

---

## 3. `type`：请求 / 响应 DTO

类型直接映射成 Go struct，字段用 **Go struct tag** 标明参数来源：

```api
type (
	LoginReq {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	GetUserInfoReq {
		Id string `form:"id"`           // GET query: ?id=
	}
	// 若路由是 /user/:id，则用 path
	// UserByIdReq {
	//   Id string `path:"id"`
	// }
)
```

### 3.1 Tag 从哪里取值

| Tag | 来源 | 典型场景 |
|-----|------|----------|
| `json` | Request / Response Body | POST / PUT / PATCH |
| `form` | URL Query（也可用于 form） | GET `?name=&phone=` |
| `path` | 路径参数 `/user/:id` | RESTful 资源 id |
| `header` | HTTP Header | `Authorization` 等 |

本仓库现状：写操作用 `json`；读操作用 `form`（如 `GetUserInfoReq`、`FindUserReq`）。

### 3.2 可选与默认值

```api
RegisterReq {
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar,optional"`   // 可不传，零值
	Sex      int32  `json:"sex,optional"`
}
SearchReq {
	Page int `form:"page,default=1"`           // 缺省用默认值
	Size int `form:"size,default=20"`
}
```

| 写法 | 含义 | 前端类比 |
|------|------|----------|
| （无修饰） | 必填（校验失败则 400） | required prop |
| `,optional` | 可选；缺省为零值 | `avatar?: string` |
| `,default=N` | 缺省填默认值 | 参数默认值 |

### 3.3 嵌套与切片

```api
LoginResp {
	Token  string     `json:"token"`
	Expire int64      `json:"expire"`
	User   UserEntity `json:"user"`      // 嵌套对象
}
FindUserResp {
	Users []UserEntity `json:"users"`    // 切片 ≈ TS 数组
}
```

前端类比：`UserEntity` ≈ 共享 interface；`[]UserEntity` ≈ `UserEntity[]`。

---

## 4. `service`：路由定义

```api
@server (
	prefix: v1/user      // 最终路径 = prefix + route
	group:  user         // 生成代码的子目录/分组名
)
service user {
	@doc "用户登录"
	@handler login
	post /login (LoginReq) returns (LoginResp)

	@doc "获取用户信息"
	@handler getUserInfo
	get /info (GetUserInfoReq) returns (GetUserInfoResp)
}
```

### 4.1 单条路由形态

```text
@doc "说明"                          // 可选
@handler 名称                        // 必填 → 生成 XxxHandler / XxxLogic
HTTP方法 /路径 (请求类型) returns (响应类型)
```

支持方法：`get` / `post` / `put` / `patch` / `delete` / `head`。

无 body 时可省略请求类型，例如：`get /health returns (HealthResp)`。  
无响应体时可省略 returns，例如：`delete /user/:id (UserReq)`。

### 4.2 路径规则

1. 必须以 `/` 开头  
2. 节点用 `/` 分隔  
3. 路径参数写成 `/:id`，对应字段必须有 `path:"id"`  
4. 节点可用字母、数字、`_`、`-`（数字节点需较新 goctl）

```api
get /info              // → /v1/user/info（加 prefix 后）
get /user/:id          // 需 Id `path:"id"`
get /foo/bar-baz
```

### 4.3 `@server` 常用键

| 键 | 作用 | 前端类比 |
|----|------|----------|
| `prefix` | 该块所有路由的 URL 前缀 | `app.use('/v1/user', router)` |
| `group` | 生成 handler/logic 的分组目录 | 按模块拆文件夹 |
| `jwt` | 启用 JWT；值为 yaml 里的配置节名 | axios 鉴权拦截器 |
| `middleware` | 自定义中间件（逗号分隔） | 路由级 middleware |
| `timeout` | 该块请求超时（如 `3s`） | fetch timeout |

同一 `service` 名可写多块，按鉴权拆分：

```api
service user {
	@handler login
	post /login (LoginReq) returns (LoginResp)   // 公开
}

@server (
	jwt: Auth
	prefix: v1/user
)
service user {
	@handler getUserInfo
	get /info (GetUserInfoReq) returns (GetUserInfoResp)  // 需登录
}
```

---

## 5. 本仓库 User HTTP 一览

完整路径 = `prefix` + route → 均在 `v1/user` 下。

| Handler | 方法 | 路径 | 请求 tag | 作用 |
|---------|------|------|----------|------|
| `register` | POST | `/register` | json | 注册并返回 token + user |
| `login` | POST | `/login` | json | 登录 |
| `getUserInfo` | GET | `/info` | form | 按 id 查单个用户 |
| `findUser` | GET | `/find` | form（均可选） | 按昵称/手机/ids 搜索 |
| `updateUserInfo` | PUT | `/info` | json | 更新昵称/头像/性别 |

契约原则（与 proto 一致）：

1. **公开 DTO 不含密码**：`UserEntity` 永不带 `password`  
2. **Req / Resp 成对命名**  
3. **复用 `UserEntity`**，避免多套半截用户结构  
4. **枚举用 int32 与 proto 对齐**：Sex / UserStatus 数值约定写在 `domain.api` 注释里  

---

## 6. 生成物与改哪里

```text
.api（契约）
  ↓ goctl
handler/     → 薄：解析参数、调 logic（少改）
types/       → 由 type 生成（勿手改）
routes.go    → 路由注册（勿手改）
logic/       → 业务实现（主要改这里）
```

| 改动意图 | 改哪里 |
|----------|--------|
| 增删字段 / 改路由 | 改 `.api` → 重新 goctl |
| 业务逻辑 | 只改 `logic/*` |
| JWT / 超时 / 前缀 | 改 `@server` + `etc/*.yaml` |

前端类比：`.api` ≈ swagger；`logic` ≈ service；`handler` ≈ 薄 controller。

---

## 7. 记忆口诀

```text
.api = syntax + type(DTO) + service(路由)
参数靠 tag：json / form / path / header
可选写 ,optional；默认写 ,default=
公开路由与鉴权路由用多块 service + jwt 拆开
改契约跑 goctl；改业务写 logic
HTTP(.api) 与 gRPC(.proto) 字段语义对齐
```

更完整官方文档：https://go-zero.dev/reference/api-dsl/
