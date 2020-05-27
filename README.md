# go-api-demo

go version: go1.11+

类Unix系統可以通过make命令运行
```
make build #编译
make tidy #更新包
make run_api #直接运行，拉起http服务
./bin/test api -c ./config/config.toml #执行二进制文件，拉起http服务
```

Windows 使用go原生命令运行
```
go build -o bin/test -v #编译
go tidy -v #更新包
go run main.go api -c ./config/config.toml #拉起http服务
```
