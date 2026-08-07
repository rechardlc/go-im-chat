# 学习记录：Users Model / DAO 查询与缓存

> 日期：2026-08-07  
> 范围：`apps/user/models/users_model.go`、`users_model_gen.go`

## 1. 分层位置

| 文件 | 角色 | 类比前端 |
|---|---|---|
| `users_model_gen.go` | goctl 生成的 CRUD + 缓存读 | 脚手架生成的 API client |
| `users_model.go` | 自定义扩展方法（如 `FindUser`） | 手写业务 data fetching |

自定义查询写在 `customUsersModel`，不改 gen 文件。

## 2. 动态 SQL（`FindUser`）

手写可选条件 + 参数绑定，企业里很常见：

- `conditions`：WHERE 片段（如 `nickname = $1`）
- `args`：实际参数值
- `argIdx`：PostgreSQL 占位符序号（`$1` → `$2` …）

要点：

- 结构可以动态拼，**值必须参数化**（防 SQL 注入）
- 条件之间当前是 `OR`（搜索语义）
- 无条件时直接返回空切片，避免全表扫描

可深入：`parameterized query`、query builder（squirrel）

## 3. CachedConn 常用 API（go-zero `sqlc`）

命名规律：`操作` + `是否 NoCache` + `是否 Ctx`

| 场景 | 方法 |
|---|---|
| 主键单条 + 缓存 | `QueryRowCtx` |
| 唯一索引单条 + 缓存 | `QueryRowIndexCtx` |
| 单条不走缓存 | `QueryRowNoCacheCtx` |
| 多行列表 | `QueryRowsNoCacheCtx`（无带缓存版本） |
| 写库并删缓存 | `ExecCtx` |
| 写库不动缓存 | `ExecNoCacheCtx` |
| 事务 | `TransactCtx` |

## 4. `FindOne` 编码含义（cache-aside）

```go
QueryRowCtx(ctx, &resp, cacheKey, func(...) {
    // 仅缓存未命中时执行：查 DB
})
```

流程：先缓存 → 命中直接返回 → 未命中跑回调查库并回写缓存。

类比：React Query 的 `queryKey` + `queryFn`。

错误处理：`nil` 成功 / `ErrNotFound` 无数据 / 其它为真实错误。

## 5. 为什么列表不走缓存？

1. **框架**：go-zero 不为多行查询提供带缓存 API（一致性难保证）
2. **业务**：`name/phone/ids` 组合爆炸、命中率低、失效难（改一条用户不知删哪些搜索 key）

| | `FindOne` | `FindUser` |
|---|---|---|
| API | `QueryRowCtx` | `QueryRowsNoCacheCtx` |
| 缓存 | 有（按 id） | 无 |
| 适合 | 详情高频读 | 动态搜索 |

## 6. 选型口诀

- 按主键 / 唯一键读单条 → 走缓存
- 动态条件 / 列表搜索 → NoCache 直查 DB
- 写操作改了被缓存的行 → `ExecCtx` 带上对应 cache keys 删除

## 关键词

`cache-aside` · `parameterized query` · `cache invalidation` · `sqlc.CachedConn`
