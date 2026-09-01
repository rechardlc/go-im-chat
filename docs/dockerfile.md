# Dockerfile 学习（对照本仓库）

前端类比：Dockerfile ≈ 一份可重复的「从源码到可运行产物」的 CI 脚本；镜像 ≈ `dist` + 运行时；容器 ≈ 用这份产物起的一个进程。

本仓库真正被 `docker-compose.prod.yml` 使用的是 **一份** `deploy/dockerfile/Dockerfile.app`。四个业务镜像（user/social × api/rpc）共用它，靠构建参数 `APP` / `TYPE` 区分编哪个包、拷哪份 yaml。

---

## 1. 谁在用、怎么传入

compose 里每个服务都是：

```yaml
build:
  context: .                              # 构建上下文 = 仓库根（才能 COPY go.mod、apps/）
  dockerfile: deploy/dockerfile/Dockerfile.app
  args:
    APP: user
    TYPE: rpc
```

`docker build` 时 `ARG APP=user TYPE=rpc`，对应源码路径 `./apps/user/rpc`。

| compose 服务 | APP | TYPE | 编出来的二进制 | 拷进镜像的配置 |
|---|---|---|---|---|
| user-rpc | user | rpc | `./apps/user/rpc` | `apps/user/rpc/etc/user.yaml` |
| social-rpc | social | rpc | `./apps/social/rpc` | `apps/social/rpc/etc/social.yaml` |
| user-api | user | api | `./apps/user/api` | `apps/user/api/etc/user.yaml` |
| social-api | social | api | `./apps/social/api` | `apps/social/api/etc/social.yaml` |

`context: .` 必须是仓库根。在 `deploy/` 下单独 build 会找不到 `go.mod` / `apps/`。

---

## 2. 多阶段构建：builder 编、alpine 跑

文件里有两个 `FROM`，这是 **multi-stage**：

| 阶段 | 基础镜像 | 干什么 | 会不会进最终镜像 |
|---|---|---|---|
| `AS builder` | `golang:*-alpine` | `go build` | 否（只有这一阶段用得到的层） |
| 第二段（无名，默认最终） | `alpine:3.21` | 只放二进制 + yaml + 探活用工具 | 是 |

前端类比：CI 里 `npm run build` 用带 node_modules 的环境，线上 Nginx 镜像只拷 `dist/`。第一阶段的 Go 工具链大约几百 MB，最终镜像只剩一个静态二进制。

```dockerfile
# 阶段 1：有编译器
FROM golang:${GO_VERSION}-alpine AS builder
...
RUN go build -o /out/app ./apps/${APP}/${TYPE}

# 阶段 2：干净运行时
FROM alpine:3.21
COPY --from=builder /out/app /app/bin
```

`COPY --from=builder` = 从上一个阶段拿文件，不是从宿主机。

---

## 3. 逐条对照 `Dockerfile.app`

### 语法与构建参数

```dockerfile
# syntax=docker/dockerfile:1
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder
ARG GOPROXY=https://goproxy.cn,direct
```

- `# syntax=`：启用 BuildKit 的 Dockerfile 前端，按 1.x 解析。
- **`ARG` 是构建期变量**，只在 `docker build` 期间存在；**`ENV` 才会进容器进程环境**。
- `ARG` 写在 `FROM` 之前，只能给 `FROM` 用；`FROM` 之后要再用 `APP`/`TYPE`，必须再声明一次（本文件 builder 和最终阶段都写了 `ARG APP` / `ARG TYPE`）。

| | ARG | ENV |
|---|---|---|
| 何时有效 | `docker build` | 容器跑起来之后 |
| 谁传入 | compose `build.args` / `--build-arg` | compose `environment` / `env_file` |
| 本仓库例子 | `APP`、`TYPE`、`GO_VERSION` | `ETCD_HOSTS`、`POD_IP`、`JWT_SECRET` |

yaml 里的 `${ETCD_HOSTS}` **不会在 build 时被替换**。镜像只是把 yaml 原样拷进去；启动时 go-zero `conf.UseEnv()` 读进程环境。

### 依赖层缓存

```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

先拷 `go.mod`/`go.sum` 再 `mod download`，再拷整个源码。前端类比：`COPY package.json` → `npm ci` → 再 `COPY src`。`go.mod` 没变时，依赖层可以复用，改业务代码不必重新下包。

### 编译

```dockerfile
ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on
RUN go build -ldflags="-s -w" -o /out/app ./apps/${APP}/${TYPE}
```

- `CGO_ENABLED=0`：纯 Go 静态链接，alpine 里不依赖 glibc。
- `-ldflags="-s -w"`：去掉符号表/DWARF，体积更小（换调试信息）。
- 四个服务都输出成同一个名字 `/out/app`，最终容器里统一是 `/app/bin`。

### 运行时镜像

```dockerfile
RUN apk add --no-cache ca-certificates tzdata curl netcat-openbsd
COPY --from=builder /out/app /app/bin
COPY apps/${APP}/${TYPE}/etc/${APP}.yaml /app/conf/config.yaml
EXPOSE 8080 8888
ENTRYPOINT ["/app/bin", "-f", "/app/conf/config.yaml"]
```

- `ca-certificates`：出容器访问 HTTPS。
- `tzdata` + `TZ=Asia/Shanghai`：日志时间。
- `curl` / `nc`：给 compose **healthcheck** 用（api 用 curl 打 8888，rpc 用 `nc -z 8080`），不是业务依赖。
- `EXPOSE` **只是文档/元数据**，不会真的把端口映射到宿主机。映射看 compose 的 `ports:`（本仓库业务容器没对外映射，只给 Nginx 走内部网络）。
- `ENTRYPOINT` 用 JSON 数组：不经 shell，信号能直接到进程（停容器时能优雅退出）。等价于本机 `./bin -f /app/conf/config.yaml`。

配置文件统一改名为 `config.yaml`，避免 api/rpc 各自 `etc/user.yaml` 路径不一致。

---

## 4. 构建产物 vs 运行时注入（和 compose 的分工）

```
docker build（镜像里有什么）
  └─ /app/bin              静态二进制
  └─ /app/conf/config.yaml 仍带 ${ETCD_HOSTS} 等占位符

docker run / compose up（进程环境）
  └─ environment <<: *app-env   → ETCD_HOSTS=etcd:2379 …
  └─ env_file: .env.prod        → POSTGRES_PASSWORD、JWT_SECRET …
  └─ POD_IP: user-rpc           → go-zero 往 etcd 登记的对外地址
```

前端类比：webpack 打进包的是 `import.meta.env`（构建期）；这里 yaml 占位符是 **启动期** 读 `process.env`。密钥不进镜像层，进容器环境。

---

## 5. 和 `Dockerfile.user_rpc_dev` 的差别

`deploy/dockerfile/Dockerfile.user_rpc_dev` 是旧/开发向写法：**宿主机先编好二进制**，镜像只负责 `COPY ./bin/...`。没有 Go 工具链，也没有多阶段。

| | `Dockerfile.app`（现用） | `Dockerfile.user_rpc_dev` |
|---|---|---|
| 在哪编译 | 镜像内 `go build` | 宿主机编完再 COPY |
| 配置 | `etc/${APP}.yaml`（带 `${}`，走环境变量） | `etc/dev/*.yaml`（本地直连） |
| 复用 | 一份文件 × 四个服务 | 写死 user-rpc |

prod 路径只认 `Dockerfile.app`。

---

## 6. 本地自己 build 一次（理解用）

必须在**仓库根**：

```bash
docker build -f deploy/dockerfile/Dockerfile.app \
  --build-arg APP=user --build-arg TYPE=rpc \
  -t easy-im-user-rpc:prod .
```

日常用 `make prod`，compose 会带上同样的 `args`。
