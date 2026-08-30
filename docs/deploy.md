# 本地两套环境

日常开发与「模拟线上」共用仓库根目录 `./data`（Postgres / Redis）。**同一时刻只能起一套**，否则会把库写坏。

## 日常开发

```bash
make infra          # 停模拟线上，只起 Postgres + Redis（5432 / 6379）
make local          # 本机 go run user-rpc + user-api（etc/dev，Endpoints 直连）
```

社交服务同样在对应目录 `go run .`，配置走 `etc/dev/*.yaml`。

## 模拟线上

与真机单机 Docker 同构：Nginx 唯一入口、etcd 注册发现、库和 RPC 不映射到宿主机。

```bash
cp deploy/.env.prod.example deploy/.env.prod   # 首次；密码须与 ./data 初次初始化一致
make prod                                      # 停日常库，构建并启动全栈
```

入口：`http://127.0.0.1/v1/user/...`、`http://127.0.0.1/v1/social/...`（默认 80，被占用则改 `deploy/.env.prod` 里的 `NGINX_PORT`）。网关探活：`GET /healthz`。

```bash
make prod-ps
make prod-logs
make prod-down
make infra          # 切回日常
```

必须在仓库根执行 `make`。不要在 `deploy/` 下直接 `docker compose -f docker-compose.prod.yml`，相对路径会指错 `./data`。

## 和真线上的差异

本机不做：K8s、HTTPS、etcd 三节点/TLS、云托管数据库。etcd 不落盘（RPC 启动时重新注册）。真服务器只需把 compose 里的 `./data` 换成数据盘路径（如 `/data/postgres`）。
