### 简介

`gin-wire` 是一个基于 `cobra` 和 `gin` 框架的脚手架，使用 wire 依赖注入进行初始化组件，解决组件之间的耦合。

### 目录结构
```
├─app
│  ├─command  命令行
│  ├─compo  组件
│  ├─cron  定时任务
│  ├─data  数据处理层
│  ├─domain  领域模型
│  ├─handler  控制层
│  ├─middleware  中间件
│  ├─model  数据库模型
│  ├─pkg  功能类库
│  │  ├─error  错误码
│  │  ├─request  请求模型
│  │  └─reposese  响应处理
│  └─service  业务逻辑层
│ 
├─bin  二进制文件目录
├─cmd  编译入口
├─conf  配置文件
├─config  配置模型
├─router  路由
├─static  静态资源（允许外部访问）
├─storage  其他静态资源存储
│  ├─app
│  │  └─public  静态资源（允许外部访问）
│  └─logs  日志目录
│ 
└─utils  工具函数
```

### 运行

- go build

  ```sh
  $ go generate
  $ go build -o ./bin/ ./...
  $ ./bin/app
  ```

- go run

  ```sh
  $ go generate
  $ go run cmd/app/main.go cmd/app/wire_gen.go cmd/app/app.go
  ```

- make

  ```sh
  $ make generate
  $ make build
  $ ./bin/app
  ```

  