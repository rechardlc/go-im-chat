user-rpc-dev:
	@make -C deploy/mk/user-rpc.mk build-test run-test

release-test: user-rpc-dev