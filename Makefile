# user-rpc-dev:
# 	@make -C deploy/mk/user-rpc.mk build-test run-test

# release-test: user-rpc-dev

# make会根据makefile所在文件的目录文件去执行命令，这里使用phony表示不用去找文件，直接执行命令
# phony 表示这个命令是一个伪命令，不会真正执行，只是为了方便管理
# 定义的命令如果后面加冒号，则表示这个命令是一个目标，会执行这个命令
# .phony 可以有多个，表示多个伪命令
.PHONY: user-api-dev user-rpc-dev local infra prod prod-down prod-logs prod-ps env-prod

user-api-dev:
	cd apps/user/api && go run .
user-rpc-dev:
	cd apps/user/rpc && go run .

# 两个都是常驻进程，必须并行；写成依赖会卡在第一个（第二个不会启动）
# -j2 让内层 Make 同时跑 2 个任务；无 -j 则串行
# make local 起外层；$(MAKE) 再起内层并把 -j2 传给它
local:
	$(MAKE) -j2 user-rpc-dev user-api-dev

# 模拟线上必须从仓库根解析 ./data 与 env_file（不要在 deploy/ 下直接 compose）
# COMPOSE_PROD := 语法是定义一个变量，后面是命令，这里定义了一个docker compose的命令，后面是参数
# --project-directory . 表示项目目录是当前目录
# -f deploy/docker-compose.prod.yml 表示使用docker-compose.prod.yml文件
# --env-file deploy/.env.prod 表示使用deploy/.env.prod文件
COMPOSE_PROD := docker compose --project-directory . -f deploy/docker-compose.prod.yml --env-file deploy/.env.prod

env-prod:
	@test -f deploy/.env.prod || cp deploy/.env.prod.example deploy/.env.prod

# 日常：只起 Postgres/Redis，本机 go run；会先停模拟线上
infra: env-prod
	-$(COMPOSE_PROD) down
	docker compose up -d

# 模拟线上：etcd + 四个业务容器 + Nginx；会先停日常库（同一 ./data 不能双开）
prod: env-prod
	-docker compose down
	$(COMPOSE_PROD) up -d --build
# 停止模拟线上
prod-down:
	-$(COMPOSE_PROD) down

# 查看模拟线上日志
prod-logs:
	$(COMPOSE_PROD) logs -f

# 查看模拟线上进程
prod-ps:
	$(COMPOSE_PROD) ps
