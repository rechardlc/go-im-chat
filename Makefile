# user-rpc-dev:
# 	@make -C deploy/mk/user-rpc.mk build-test run-test

# release-test: user-rpc-dev

.PHONY: user-api-dev user-rpc-dev release-test

user-api-dev:
	cd apps/user/api && go run .
user-rpc-dev:
	cd apps/user/rpc && go run .

# 两个都是常驻进程，必须并行；写成依赖会卡在第一个
local:
	$(MAKE) -j2 user-rpc-dev user-api-dev
