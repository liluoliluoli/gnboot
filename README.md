<h1 align="center">gnboot</h1>

<div align="center">
由gin + gorm + jwt + casbin组合实现的RBAC权限管理脚手架Golang版, 搭建完成即可快速、高效投入业务开发
</div>

## 特性

- `RESTful API` 设计规范
- `Gin` 一款高效的golang web框架
- `MySQL` 数据库存储
- `Jwt` 用户认证, 登入登出一键搞定
- `Casbin` 基于角色的访问控制模型(RBAC)
- `Gorm` 数据库ORM管理框架, 可自行扩展多种数据库类型(主分支已支持gorm 2.0)
- `Validator` 请求参数校验, 版本V9
- `Log` v1.2.2升级后日志支持两种常见的高性能日志 logrus / zap (移除日志写入本地文件, 强烈建议使用docker日志或其他日志收集工具)
- `Viper` 配置管理工具, 支持多种配置文件类型
- `Embed` go 1.16文件嵌入属性, 轻松将静态文件打包到编译后的二进制应用中
- `DCron` 分布式定时任务，同一task只在某台机器上执行一次(需要配置redis)
- `GoFunk` 常用工具包, 某些方法无需重复造轮子
- `FiniteStateMachine` 有限状态机, 常用于审批流程管理(没有使用工作流, 一是go的轮子太少, 二是有限状态机基本可以涵盖常用的审批流程)
- `Uploader` 大文件分块上传/多文件、文件夹上传Vue组件[vue-uploader](https://github.com/simple-uploader/vue-uploader/)
- `MessageCenter` 消息中心(websocket长连接保证实时性, 活跃用户上线时新增消息表, 不活跃用户不管, 有效降低数据量)
- `testing` 测试标准包, 快速进行单元测试
- `Grafana Loki` 轻量日志收集工具loki, 支持分布式日志收集(需要通过docker运行[gnboot-docker](https://github.com/piupuer/gnboot-docker))
- `Minio` 轻量对象存储服务(需要通过docker运行[gnboot-docker](https://github.com/piupuer/gnboot-docker))
- `Swagger` Swagger V2接口文档
- `Captcha` 密码输错次数过多需输入验证码
- `Sign` API接口签名(防重放攻击、防数据篡改)
- `Opentelemetry` 链路追踪, 快速分析接口耗时

## 中间件

- `Rate` 访问速率限制中间件 -- 限制访问流量
- `Exception` 全局异常处理中间件 -- 使用golang recover特性, 捕获所有异常, 保存到日志, 方便追溯
- `Transaction` 全局事务处理中间件 -- 每次请求无异常自动提交, 有异常自动回滚事务, 无需每个service单独调用(GET/OPTIONS跳过)
- `AccessLog` 请求日志中间件 -- 每次请求的路由、IP自动写入日志
- `Cors 跨域中间件` -- 所有请求均可跨域访问
- `Jwt` 权限认证中间件 -- 处理登录、登出、无状态token校验
- `Casbin` 权限访问中间件 -- 基于Cabin RBAC, 对不同角色访问不同API进行校验
- `Idempotence` 接口幂等性中间件 -- 保证接口不受网络波动影响而重复点击或提交(目前针对create接口加了处理，可根据实际情况更改)

## 默认菜单

- 首页
- 系统管理
    - 菜单管理
    - 角色管理
    - 用户管理
    - 接口管理
    - 数据字典
    - 操作日志
    - 消息推送
    - 机器管理
- 状态机
    - 状态机配置
    - 我的请假条
    - 待审批列表
- 上传组件
    - 上传示例1
    - 上传示例2(主要是针对ZIP压缩包上传及解压)
- 测试页面
    - 测试用例

## 快速开始

```
git clone https://github.com/piupuer/gnboot
cd gnboot
# 强烈建议使用golang官方包管理工具go mod, 无需将代码拷贝到$GOPATH/src目录下
# 确保go环境变量都配置好, 运行main文件
go run main.go
```

> 启动成功之后, 可在浏览器中输入: [http://127.0.0.1:10000/api/ping](http://127.0.0.1:10000/api/ping), 若不能访问请检查Go环境变量或数据库配置是否正确

## 项目结构概览

```
├── api
│   └── v1 # v1版本接口目录(类似于Java中的controller), 如果有新版本可以继续添加v2/v3
├── conf # 配置文件目录(包含测试/预发布/生产环境配置参数及casbin模型配置)
├── docker-conf # docker相关配置文件
├── initialize # 数据初始化目录
│   ├── db # 数据库初始化脚本目录, 遵循sql-migrate规范
│   └── xxx.go # 包含各种需要初始化的全局变量, 如mysql/redis
├── middleware # 中间件目录
├── models # 存储层模型定义目录
├── pkg # 公共模块目录
│   ├── cache_service # redis缓存服务目录
│   ├── global # 全局变量目录
│   ├── redis # redis查询工具目录
│   ├── request # 请求相关结构体目录
│   ├── response # 响应相关结构体目录
│   ├── service # 数据DAO服务目录
│   ├── utils # 工具包目录
│   └── wechat # 微信接口目录
├── router # 路由目录
├── tests # 本地单元测试配置目录
├── upload # 上传文件默认目录
├── Dockerfile # docker镜像构建文件(生产环境)
├── Dockerfile.stage # docker镜像构建文件(预发布环境)
├── go.mod # go依赖列表
├── go.sum # go依赖下载历史
├── main.go # 程序主入口
├── README.md # 说明文档
├── TIPS.md # 个人踩坑记录
├── TODO.md # 已完成/待完成列表
```



