# Protocol Buffers（.proto）核心速查

面向本仓库（Go + gRPC）的复习笔记。对照示例：`apps/user/rpc/user.proto`。

前端类比：`.proto` ≈ **共享 TypeScript interface + API 路由表**，再由工具生成各端 SDK。

---

## 1. 它是什么

| 概念 | 说明 |
|------|------|
| IDL | 接口描述语言，只定义契约，不含业务实现 |
| 用途 | 微服务 / gRPC 的请求响应结构、服务方法列表 |
| 产物 | `protoc` / goctl 生成 `.pb.go`、客户端 stub、服务端接口 |

流程：写 `.proto` → 生成代码 → 服务端实现 RPC → 客户端调 stub。

---

## 2. 文件骨架

```protobuf
syntax = "proto3";          // 语法版本，建议第一行
package user;               // 命名空间，防跨文件重名

import "google/protobuf/timestamp.proto";  // 引用其他 proto

option go_package = "github.com/rechardlc/go-im-chat/apps/user/rpc/user";
// ↑ 生成 Go 代码时的 import path，需与 go.mod module 对齐
```

| 元素 | 作用 | 前端类比 |
|------|------|----------|
| `syntax` | proto2 / proto3 | tsconfig 的语言版本 |
| `package` | 命名空间 | namespace / 模块名 |
| `import` | 引用别的契约 | `import type` |
| `option` | 给代码生成器的元数据 | 编译器配置 |

---

## 3. 三类核心定义

一份业务 proto 通常就这三块：

```text
enum     → 有限取值集合
message  → 数据结构（DTO）
service  → RPC 方法列表
```

### 3.1 `enum` 枚举

```protobuf
enum Sex {
  SEX_UNSPECIFIED = 0;  // proto3 必须有 0，表示「未指定/默认」
  SEX_UNKNOWN = 1;
  SEX_MALE = 2;
  SEX_FEMALE = 3;
}
```

要点：

- 默认值永远是 `0`，所以 `0` 不要赋业务含义（如「男」）
- 命名习惯：`ENUM类型_取值`，避免不同 enum 生成代码冲突
- 比裸 `int32` / `string` 更不容易前后端约定漂移

### 3.2 `message` 消息（DTO）

```protobuf
message UserEntity {
  string id = 1;           // 类型 字段名 = 字段编号;
  string nickname = 3;
  UserStatus status = 5;   // 可用 enum / 其他 message
  Sex sex = 6;
}
```

| 概念 | 说明 |
|------|------|
| 字段编号 `= N` | 序列化身份；**已发布后不要改含义、不要复用删掉的编号** |
| `1–15` | 编码更省，留给高频字段 |
| 字段名 | 可读性用；改名通常不影响线兼容（编号不变即可） |

字段基本形态：`[修饰符] 类型 字段名 = 编号 [选项];`

例：`repeated UserEntity users = 1;`

| 部分 | 角色 | 本例 |
|------|------|------|
| 修饰符（label） | 基数：几个值 | `repeated` |
| 类型 | 标量 / enum / message | `UserEntity` |
| 字段名 | 可读名 | `users` |
| `= 编号` | 序列化身份 | `1` |

#### 字段修饰一览（proto3）

**① 基数 label（真正的修饰符）**

| 写法 | 含义 | 前端类比 |
|------|------|----------|
| （省略） | 单个值；未设时是类型零值 | `name: string` |
| `optional` | 单个值，能区分「未设置」与零值 | `email?: string`（带 presence） |
| `repeated` | 列表，0～N 个 | `users: UserEntity[]` |

```protobuf
string name = 1;                 // 默认单值
optional string email = 2;       // 可判断有没有传
repeated UserEntity users = 3;   // 列表
```

注意：proto2 的 `required` 在 **proto3 已移除**，不要写。

**② 结构写法（不是 label，但常一起记）**

| 写法 | 含义 | 前端类比 |
|------|------|----------|
| `map<K, V>` | 字典；key 限整数 / bool / string | `Record<K, V>` |
| `oneof { ... }` | 一组字段互斥，最多设一个 | 联合类型 / 二选一表单 |
| `UserEntity user` | 引用另一个 message（组合） | 嵌套对象 |

```protobuf
map<string, string> extra = 4;

oneof credential {
  string password = 5;
  string sms_code = 6;
}

UserEntity user = 7;   // 组合：响应里带一个 user 对象
```

**③ 字段选项（写在 `[]` 里）**

```protobuf
string old_nick = 8 [deprecated = true, json_name = "nick"];
```

| 常见选项 | 作用 |
|----------|------|
| `deprecated = true` | 标记废弃，生成代码可能告警 |
| `json_name = "..."` | 指定 JSON 序列化时的字段名 |

速记：

```text
label:     optional / repeated /（默认单值）
结构:      map / oneof / 引用 message
元数据:    [option = ...]
```

本仓库示例：`FindUserResp` 用 `repeated UserEntity users = 1` 表示用户列表。

**组合 ≠ 嵌套定义：**

```protobuf
// 组合（本仓库用法）：顶层定义 UserEntity，再引用
message LoginResp {
  string token = 1;
  int64 expire = 2;
  UserEntity user = 3;   // 响应里带一个 user 对象
}

// 嵌套定义：类型写在 message 内部（较少用，复用差）
message LoginResp {
  message UserBrief { string id = 1; }
  UserBrief user = 1;
}
```

`LoginResp` 序列化后的直观形态：

```json
{
  "token": "xxx",
  "expire": 7200,
  "user": {
    "id": "1",
    "nickname": "张三",
    "phone": "138..."
  }
}
```

### 3.3 `service` / `rpc`（gRPC 接口）

```protobuf
service User {
  rpc Login(LoginReq) returns (LoginResp);           // 一元：一请求一响应
  rpc FindUser(FindUserReq) returns (FindUserResp);
}
```

四种 RPC 形态：

| 形态 | 写法 | 场景 |
|------|------|------|
| 一元 | `rpc A(Req) returns (Resp)` | 登录、查用户（最常用） |
| 服务端流 | `returns (stream Resp)` | 服务端持续推送 |
| 客户端流 | `rpc A(stream Req)` | 客户端批量上传 |
| 双向流 | `stream` 两边都有 | IM 长连接式交互 |

前端类比：`service` ≈ API 模块；每个 `rpc` ≈ 一个函数签名。

---

## 4. 常用标量类型

| proto | Go（常见） | 备注 |
|-------|------------|------|
| `string` | `string` | UTF-8 |
| `bool` | `bool` | |
| `bytes` | `[]byte` | 二进制 |
| `int32` / `int64` | `int32` / `int64` | 变长；负数用 `sint*` 更省 |
| `uint32` / `uint64` | `uint*` | 无符号 |
| `float` / `double` | `float32` / `float64` | 浮点 |

本仓库约定：`LoginResp.expire` / `RegisterResp.expire` 为 **token 有效期（秒）**，类型 `int64`。

---

## 5. 进阶语法（知道即可）

`oneof` / `map` / 字段选项见上文「字段修饰一览」。此处补充兼容性与官方类型：

```protobuf
reserved 10, 11 to 15;      // 预留编号，防止旧字段被误复用
reserved "password";        // 预留字段名

import "google/protobuf/timestamp.proto";
google.protobuf.Timestamp created_at = 7;  // 官方时间类型
```

---

## 6. 契约设计原则（结合本仓库 User 服务）

1. **公开 DTO 不含敏感字段**：`UserEntity` 永不带 `password`；密码只出现在 `LoginReq` / `RegisterReq`。
2. **Req / Resp 成对命名**：`LoginReq` ↔ `LoginResp`，语义一眼可读。
3. **复用实体**：登录、注册、查资料、更新资料都复用 `UserEntity`，避免多套「半截用户结构」。
4. **enum 表达状态/性别**：比魔法数字 / 随意字符串稳。
5. **编号稳定性**：上线后只追加字段；删除用 `reserved`，不要复用旧编号。

当前 User RPC 一览：

| RPC | 作用 |
|-----|------|
| `Ping` | 连通性探测 |
| `Login` | 手机号 + 密码登录 → token + user |
| `Register` | 注册并自动登录 → token + user |
| `GetUserInfo` | 按 id 查单个用户 |
| `FindUser` | 按昵称 / 手机 / ids 查找列表 |
| `UpdateUserInfo` | 更新昵称 / 头像 / 性别 |

---

## 7. 记忆口诀

```text
.proto = package 元数据 + message/enum 数据结构 + service RPC 接口
字段靠编号兼容，不靠字段名
修饰：optional / repeated；结构：map / oneof；选项：[...]
password 进 Req，不进 Entity
```

更完整的官方文档：搜索 `Protocol Buffers Language Guide proto3`。
