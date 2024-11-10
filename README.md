<h1 align="center">Gnboot Service</h1>

├── api                  // 各个微服务的proto/go文件
│   ├── reason           //错误码pb
│   ├── xxx              // xxx服务所需的go文件，在porto文件创建proto后，通过命令生成：make api
│   ├── proto            // proto文件
│   └── ...
├── cmd                  
│   └── gnboot             // 项目名
│       ├── main.go      // 程序主入口
│       ├── wire.go      // wire依赖注入，自动生成：make gen
│       └── wire_gen.go
├── configs              // 配置文件目录
│   ├── config.yml       // 主配置文件
│   ├── client.yml       // 配置grpc服务client，这期用不到
│   ├── gen.yml          // gen gorm或migrate会用到的配置文件，见文件内命令执行
│   └── ...              // 其他自定义配置文件以yml/yaml结尾均可
├── internal             // 内部逻辑代码
│   ├── biz              // 业务逻辑的组装层, 类似 DDD 的 domain 层, data 类似 DDD 的 repo, 而 repo 接口在这里定义, 使用依赖倒置的原则.
│   │   ├── biz.go       //用来构造对象给wire使用
│   │   ├── reason.go    // 定义错误描述
│   │   └── xxx.go       // 具体业务
│   ├── conf
│   │   ├── conf.pb.go
│   │   └── conf.proto   // 内部使用的config的结构定义, 使用proto格式生成
│   ├── data             // 业务数据访问, 包含 cache、db 等封装, 实现了 biz 的 repo 接口. 我们可能会把 data 与 dao 混淆在一起, data 偏重业务的含义, 它所要做的是将领域对象重新拿出来, 我们去掉了 DDD 的 infra层.
│   │   ├── model        // gorm gen生成model目录
│   │   ├── query        // gorm gen生成query目录
│   │   ├── cache.go     // cache层, 防缓存击穿/缓存穿透/缓存雪崩
│   │   ├── client.go    // 各个微服务client初始化
│   │   ├── data.go      // 数据初始化, 如DB/Redis，用来构造对象给wire使用
│   │   ├── xxx.go       // 具体业务repo
│   │   └── tracer.go    // 链路追踪tracer初始化
│   ├── db
│   │   ├── migrations   // sql迁移文件目录, 每一次数据库变更都放在这里, 参考https://github.com/rubenv/sql-migrate
│   │   │   ├── xxx.sql  // sql文件,文件名定义每次要比上一次大
│   │   │   └── ...
│   │   └── migrate.go   // embed sql文件，启动时会执行数据初始化和迁移（如果库里没执行过上诉sql）
│   ├── pkg              // 自定义扩展包
│   │   ├── idempotent   // 接口幂等性
│   │   ├── task         // 异步任务, 内部调用asynq
│   │   └── xxx          // 其他扩展
│   ├── server           // http和grpc实例的创建和配置
│   │   ├── middleware   // 自定义中间件
│   │   │   ├── locales  // i18n多语言map配置文件
│   │   │   └── xxx.go   // 一些中间件
│   │   ├── grpc.go      //启动server
│   │   ├── http.go      //启动server
│   │   └── server.go    //用来构造对象给wire使用
│   └── service          // 实现了 api 定义的服务层, 类似 DDD 的 application 层, 处理 DTO 到 biz 领域实体的转换(DTO -> DO), 同时协同各类 biz 交互, 但是不应处理复杂逻辑
│       ├── service.go   //用来构造对象给wire使用
│       └── xxx.go       // 业务接口入口
├── third_party          // api依赖的第三方proto, 编译proto文件需要用到
│   ├── cinch            // cinch公共依赖
│   ├── errors
│   ├── google
│   ├── openapi
│   │── validate
│   └── ...              //  其他自定义依赖
├─ .gitignore
├─ .gitmodules           // submodule配置文件
├─ .golangci.yml         // golangci-lint
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ LICENSE
├─ Makefile
└─ README.md




https://go-kratos.dev/docs/component/metrics/
https://blog.csdn.net/qq_43141726/article/details/140051769
https://go-cinch.github.io/docs/#/