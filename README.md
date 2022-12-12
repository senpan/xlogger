# xlogger

日志作为整个代码行为的记录，是程序执行逻辑和异常最直接的反馈 ，`xlogger`日志组件，支持标准输出和高性能磁盘写入,可基于配置文件完成快速启动，使用起来简单方便，5个日志级别满足项目中各种需求。

本工具包主要是对于`zap`和`logrus`的封装

Fork from : https://github.com/tal-tech/loggerX


#### 下载安装
```bash
go get github.com/senpan/xlogger
```
#### 初始化logger：

- 在`main`函数中初始化
```golang
import (
    logger "github.com/senpan/xlogger"
)

version := "日志版本号" // 自定义
logger.InitXLogger(version)
```
- 配置文件内容
```yaml
Logger:
  # 日志组件:logrus/zap
  lib: zap
  # 日志模式:stdout/file
  mode: stdout
  # level:DEBUG/INFO/WARNING/ERROR
  level: DEBUG
  # 日志文件
  filename: ./logs/pangu.log
  # 日志文件大小，单位:MB
  maxSize: 128
  # 最大过期日志保留个数
  maxBackups: 2
  # 保留过期文件最大时间，单位:天
  maxAge: 3
  #是否压缩日志，默认是不压缩
  compress: true
```

#### 打印日志方法

* 支持不同级别打印日志的方法：

```golang
import (
    logger "github.com/senpan/xlogger"
)

logger.D(tag string, args interface{}, v ...interface{})
logger.I(tag string, args interface{}, v ...interface{})
logger.W(tag string, args interface{}, v ...interface{})
logger.E(tag string, args interface{}, v ...interface{})
logger.F(tag string, args interface{}, v ...interface{}) // 触发panic
logger.Dx(ctx context.Context, tag string, args interface{}, v ...interface{})
logger.Ix(ctx context.Context, tag string, args interface{}, v ...interface{})
logger.Wx(ctx context.Context, tag string, args interface{}, v ...interface{})
logger.Ex(ctx context.Context, tag string, args interface{}, v ...interface{})
logger.Fx(ctx context.Context, tag string, args interface{}, v ...interface{}) // 触发panic
```

* 使用用例:

```golang
import (
    logger "github.com/senpan/xlogger"
)

logger.I("info.tag", "data save to mysql, uid:%d ,name:%s", 2015, "Cindy")
logger.Ix(ctx, "info.tag.ctx","data save to mysql, uid:%d ,name:%s", 2015, "Cindy")
logger.E("error.tag", "get redis error：%v, uid:%d ,name:%s", err, 2015, "Cindy")
logger.Ex(ctx, "error.tag.ctx","get redis error:%v, uid:%d ,name:%s", err, 2015, "Cindy")
```

* 支持携带ctx的打印方法:

```golang
// 每次调用携带全局变量，支持log特殊需求，如链路追踪等
// 一次请求开始时写入
ctx = context.WithValue(ctx, "__svc_start__", time.Now()) // 每条日志会计算出相对接口开始时间耗时
```

#### 注意事项：

* logger库是并发不安全的，所以全局只能有一个实例。在写单元测时，有可能会多次初始化，此时一定要在包测试完之后进行`Close()`操作。