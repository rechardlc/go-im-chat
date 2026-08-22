# 生成rpc代码
goctl rpc protoc ./apps/social/rpc/social.proto --go_out=./apps/social/rpc/social --go-grpc_out=./apps/social/rpc/social --zrpc_out=./apps/social/rpc --style=go_zero

# go-zero 官方推荐使用 postgresql 生成 model 代码的方式
# postgresql 生成 model 代码（本地 docker；users 表在 chat schema，非 public）
goctl model pg datasource -url='postgres://chat:123456@127.0.0.1:5432/chat?sslmode=disable' -schema='chat' -table='friends,friend_requests,groups,group_members,group_requests' -style=go_zero -dir='./apps/social/models/' -c

# 生成 api 代码
goctl api go -api ./apps/social/api/social.api -dir ./apps/social/api -style go_zero
     