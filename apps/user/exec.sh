# 生成rpc代码
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc/user --go-grpc_out=./apps/user/rpc/user --zrpc_out=./apps/user/rpc --style=go_zero

# go-zero 官方推荐使用 postgresql 生成 model 代码的方式
# postgresql 生成 model 代码（本地 docker；users 表在 chat schema，非 public）
goctl model pg datasource -url='postgres://chat:123456@127.0.0.1:5432/chat?sslmode=disable' -schema='chat' -table='users' -style=go_zero -dir='./apps/user/models/' -c