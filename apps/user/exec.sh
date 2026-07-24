# 生成rpc代码
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc/user --go-grpc_out=./apps/user/rpc/user --zrpc_out=./apps/user/rpc --style=go_zero